package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	//"github.com/robertsmoto/skustor/internal/configs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
)

type downloader interface {
	Download() (err error)
}

type resizer interface {
	Resize() (err error)
}

type processor interface {
	Process()
}

type uploader interface {
	UploadToSpaces() (err error)
}

type downloaderResizerProcessorUpserterUploader interface {
	downloader
	resizer
	processor
	upserter
	uploader
}

func ImgHandler(i downloaderResizerProcessorUpserterUploader, userId string, db *sql.DB) {
	var err error
	err = i.Download()
	err = i.Resize()
	i.Process()
	i.Upsert(userId, db)
	err = i.UploadToSpaces() // to AWS cdn
	if err != nil {
		log.Print("Error image.Resize() ", err)
	}
}

type ImgSize struct {
	ratio float64
	size  string
}

type ResizedImage struct {
	tempFilePath string
	key          string
	url          string
	width        string
	height       string
	size         string // eg "LG" "MD" or "SM"
}

type Image struct {
	Id             string `json:"id" validate:"omitempty,uuid4"`
	Url            string `json:"url" validate:"omitempty,url"`
	urlHash        string `json:"-"`  
	Title          string `json:"title" validate:"omitempty,lte=200"`
	Alt            string `json:"alt" validate:"omitempty,lte=100"`
	Caption        string `json:"caption" validate:"omitempty,lte=200"`
	Position       uint8  `json:"position" validate:"omitempty,number"`
	Process        uint8  `json:"process" validate:"required,number,oneof=0 1"`
	Height         string `json:"height" validate:"omitempty,lte=20"`
	Width          string `json:"width" validate:"omitempty,lte=20"`
	Size           string `json:"size" validate:"omitempty,lte=20,oneof=LG MD SM"`
    SvUserId       string `json:"-"`  // constructed
	TempFileDir    string `json:"-"`  // constructed :: location to download temp files
	UploadPrefix   string `json:"-"`  // constructed :: eg. "media"
	VanityUrl      string `json:"-"`  // constructed
	UserDir        string `json:"-"`  // constructed :: eg. "98c56d78fe3a"
	Date           string `json:"-"`  // constructed :: eg. "2020-05-04"
	DoBucket       string `json:"-"`  // constructed
	DoCacheControl string `json:"-"`  // constructed eg. "max-age=60"
	DoContentType  string `json:"-"`  // constructed eg. "image/webp"
	DoEndpointUrl  string `json:"-"`  // constructed
	DoAccessKey    string `json:"-"`  // constructed
	DoSecret       string `json:"-"`  // constructed
	DoRegionName   string `json:"-"`  // constructed
	BaseFileName   string `json:"-"`  
	fullFileName   string `json:"-"`  
	filePath       string `json:"-"`  
	ImgSizes       []ImgSize `json:"-"`  
	ResizedImages  []ResizedImage `json:"-"`  
}

type ImageNodes struct {
	Nodes []*Image `json:"imageNodes" validate:"dive"`
}

func (s *ImageNodes) Load(fileBuffer []byte) {
	var err error
	json.Unmarshal(fileBuffer, &s)
	if err != nil {
		log.Print("Image.Load() ", err)
	}
}

func (s *ImageNodes) Validate() {
	var err error
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("Image.Validate() ", err)
	}
}

func (s *ImageNodes) RecordOriginalImage(db *sql.DB) (err error) {
    // Checks if image record exists in db.
    for _, node := range s.Nodes {
        node.urlHash = Md5Hasher([]string{node.Url})
        // store all original image data in db for future checks,
        // so images are only processed one time.
        documentMap := map[string]interface{}{
            "url": node.Url, 
            "title": node.Title,
            "alt": node.Alt,
            "caption": node.Caption,
            "position": node.Position,
            "height": node.Height,
            "width": node.Width,
            "size": node.Size,
        }
        documentJson, _ := json.Marshal(documentMap)
        qstr := `
            INSERT INTO image (id, document)
            values ($1, $2)
            ON CONFLICT (id) DO UPDATE
            SET document = $2
            WHERE image.id = $1;`
        _, err = db.Exec(qstr, node.urlHash, string(documentJson))
    }
    return nil
}

func (s *ImageNodes) Download() (err error) {
    for _, node := range s.Nodes {
        if node.Process == 0 {
            continue
        }
        isValid := ValidateUrl(node.Url)
        if isValid == false {
            return fmt.Errorf("Image url is not valid %s", node.Url)
        }
        baseFileName, fullFileName := GetFileName(node.Url)
        node.BaseFileName = baseFileName
        node.filePath = filepath.Join(
            os.Getenv("TMPDIR"),
            "images/downloads",
            fullFileName,
        )
        if err != nil {
            return fmt.Errorf("Internal error Image.Download() 01")
        }
        err = DownloadFile(node.Url, node.filePath)
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
    }
	return err
}

func (s *ImageNodes) Resize() (err error) {
    for _, node := range s.Nodes {
        // Resizes, renames and encodes original image.
        if node.Process == 0 {
            continue
        }
        if node.ImgSizes == nil {
            return err
        }
        // Check that download file exists
        _, err = os.Stat(node.filePath)
        if err != nil {
            return fmt.Errorf("Download file does not exist: %s", err)
        }
        for _, size := range node.ImgSizes {
            imgFile, err := os.Open(node.filePath)
            if err != nil {
                return fmt.Errorf("Resize open error 01: %s", err)
            }
            defer imgFile.Close()
            imgConfig, _, err := image.DecodeConfig(imgFile)
            imgWidth := imgConfig.Width
            imgHeight := imgConfig.Height
            imgFile, err = os.Open(node.filePath)
            if err != nil {
                return fmt.Errorf("Resize open error 02: %s", err)
            }
            defer imgFile.Close()
            decodedImage, _, err := image.Decode(imgFile)
            if err != nil {
                return fmt.Errorf("Resize decode: %s", err)
            }
            // calculate new image sizes
            newWidth, newHeight := CalcNewSize(imgWidth, imgHeight, size.ratio)
            // create new file name and dirs
            newFileName := CreateNewFileName(node.BaseFileName, newWidth, newHeight)
            tempFilePath := filepath.Join(
                os.Getenv("TMPDIR"),
                "images/resized",
                newFileName,
            )
            uploadPath := CreateUploadPath(node.UploadPrefix,
                node.UserDir, newFileName, node.Date)
            // create local dir
            // if successful, the created file can be used for I/O
            var f *os.File
            f, err = os.Create(tempFilePath)
            if err != nil {
                return fmt.Errorf("Create tempFilePath %s", err)
            }
            // resize the image in memory
            // resizedImage is of image.Image type
            resizedImage := ResizeBaseImage(decodedImage, newWidth, newHeight)
            // encode image
            // Encode takes two arguments, io.writer and image.Image
            err = webpbin.Encode(f, resizedImage)
            if err != nil {
                return fmt.Errorf("Image encoding error %s", err)
            }
            f.Close()
            // add resize information to struct
            rsi := ResizedImage{}
            rsi.tempFilePath = tempFilePath
            rsi.key = uploadPath // aws key
            rsi.url = filepath.Join(node.VanityUrl, uploadPath)
            rsi.width = fmt.Sprint(newWidth)
            rsi.height = fmt.Sprint(newHeight)
            rsi.size = size.size
            node.ResizedImages = append(node.ResizedImages, rsi)
        }
    }
	return err
}

func (s *ImageNodes) Upsert(db *sql.DB) (err error) {
	// now add the resizeImg.NewSizes to the image table
    for _, node := range s.Nodes {
        if node.Process == 0 {
            continue
        }
        for _, i := range node.ResizedImages {
            // creates unique id on url, size
            idHash := Md5Hasher([]string{i.url, i.size})
            documentMap := map[string]interface{}{
                "url": i.url, 
                "title": node.Title,
                "alt": node.Alt,
                "caption": node.Caption,
                "position": node.Position,
                "height": i.height,
                "width": i.width,
                "size": i.size,
            }
            documentJson, _ := json.Marshal(documentMap)
            qstr := `
                INSERT INTO image (id, parent_id, document)
                values ($1, $2, $3)
                ON CONFLICT (id) DO UPDATE
                SET parent_id = $2,
                    document = $3
                WHERE image.id = $1;`
            _, err = db.Exec(qstr, idHash, node.urlHash, string(documentJson))
            if err != nil {
                return fmt.Errorf("Image.Upsert %s", err)
            }
        }
    }
	return nil
}

func (s *ImageNodes) UploadToSpaces() (err error) {
	//img {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]
    for _, node := range s.Nodes {
        if node.Process == 0 {
            continue
        }
        if node.ResizedImages == nil {
            return err
        }
        customResolver := aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{URL: node.DoEndpointUrl}, nil
            },
        )

        cfg, err := config.LoadDefaultConfig(context.TODO(),
            config.WithEndpointResolverWithOptions(customResolver),
            config.WithCredentialsProvider(
                credentials.NewStaticCredentialsProvider(
                    node.DoAccessKey, node.DoSecret, "")),
        )

        if err != nil {
            panic("AWS configuration error, " + err.Error())
        }

        client := s3.NewFromConfig(cfg, func(o *s3.Options) {
            o.Region = node.DoRegionName
        })

        for _, rsi := range node.ResizedImages {
            file, err := os.Open(rsi.tempFilePath)
            if err != nil {
                log.Print("Unable to open file ", rsi.tempFilePath)
                return err
            }

            defer file.Close()

            input := &s3.PutObjectInput{
                Bucket:       &node.DoBucket,
                Key:          &rsi.key, // <-- uploadPath
                Body:         file,
                CacheControl: &node.DoCacheControl,
                ContentType:  &node.DoContentType,
                ACL:          "public-read",
            }

            _, err = PutFile(context.TODO(), client, input)
            if err != nil {
                fmt.Print(err)
                return err
            }
        }
    }
	return err
}

func (s *ImageNodes) RemoveTempFile() (err error) {
    for _, node := range s.Nodes {
        if node.Process == 0 {
            continue
        }
        e := os.Remove(node.filePath)
        if e != nil {
            return fmt.Errorf("Image.RemoveTempFile %s", err)
        }
    }
    return nil
}

func ValidateUrl(url string) bool {
	//ValidateUrl checks that the url and filetype is valid.
	if strings.HasPrefix(url, "http") == false {
		return false
	}
	allowedFormats := []string{"jpg", "JPG", "jpeg", "JPEG", "png", "PNG"}
	for _, f := range allowedFormats {
		if strings.HasSuffix(url, f) {
			return true
		}
	}
	return false
}

func GetFileName(url string) (string, string) {
	//GetFileName returns the fullFileName and baseFileName (with and w/o .extension)
	fullFileName := filepath.Base(url)
	baseFileName := strings.Split(fullFileName, ".")[0]
	return baseFileName, fullFileName
}

func DownloadFile(url, filePath string) (err error) {
	//DownloadFile downloads a file from a url to filePath.
	//It writes as it downloads rather than loading the entire file into memory.

	// checks the filePath dir if it doesn't exist creates it
	dir := filepath.Dir(filePath)
	_, err = os.Stat(dir)
	if err != nil {
		err := os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ResizeBaseImage(decodedImg image.Image, w, h int) image.Image {
	//ResizeBaseImage resizes a given image width, hieght by the sizeFactor.
	resizedImage := resize.Resize(uint(w), uint(h), decodedImg, resize.NearestNeighbor)
	return resizedImage
}

func CreateNewFileName(fileName string, w, h int) string {
	//CreateNewFileName creates the new filename with size and .webp extension.
	fileName += "_" + strconv.Itoa(w) + "x" + strconv.Itoa(h) + ".webp"
	return fileName
}

//func CreateTempFilePath(tempFileDir, fileName string) (tempPath string) {
//tempPath = filepath.Join(tempFileDir, fileName)
//return tempPath
//}

func CreateUploadPath(
	uploadPrefix, userDir, fileName, date string) string {
	//CreateUploadPath creates a new path based on uploadPrefix eg. "media"
	//userDir eg. "d91216fbb4d4" and fileName.
	//Path fomat /<uploadPrefix>/<userPath>/year/month/data/filename
	//Date should be in string format "2022-06-02"

	//year, month, day := time.Now().Date()
	dateSlice := strings.Split(date, "-")
	y := dateSlice[0]
	m := dateSlice[1]
	d := dateSlice[2]
	key := filepath.Join(uploadPrefix, userDir, y, m, d, fileName)
	return key
}

func CalcNewSize(width, height int, ratio float64) (int, int) {
	newWidth := int(math.Round(float64(width) * ratio))
	newHeight := int(math.Round(float64(height) * ratio))
	return newWidth, newHeight
}

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PutFile(
	c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (
	*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}
