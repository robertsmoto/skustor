package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
    "errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
    "time"

	//"github.com/robertsmoto/skustor/internal/configs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
	"github.com/tidwall/gjson"
)


type resizeConstructor interface {
    Construct(accountId string)
}

type downloader interface {
	Download() (err error)
}

type resizer interface {
	Resize() (err error)
}

type resizedImageUpserter interface {
    ResizedImageUpsert(accountId string, db *sql.DB) (err error)
}

type spacesUploader interface {
	UploadToSpaces() (err error)
}

type tempFileRemover interface {
    RemoveTempFile() (err error)
}

type resizedImgHandler interface {
    resizeConstructor
	downloader
	resizer
    resizedImageUpserter
	spacesUploader
    tempFileRemover
}

func ResizedImgHandler(i resizedImgHandler, accountId string, db *sql.DB) (err error){
    i.Construct(accountId)
	err = i.Download()
	if err != nil {
		return fmt.Errorf("ResizedImgHandler 01 %s", err)
	}
	err = i.Resize()
	if err != nil {
		return fmt.Errorf("ResizedImgHandler 02 %s", err)
	}
	err = i.ResizedImageUpsert(accountId, db)
	if err != nil {
		return fmt.Errorf("ResizedImgHandler 03 %s", err)
	}
	err = i.UploadToSpaces() // to AWS cdn
	if err != nil {
		return fmt.Errorf("ResizedImgHandler 04 %s", err)
	}
	err = i.RemoveTempFile()
	if err != nil {
		return fmt.Errorf("ResizedImgHandler 05 %s", err)
	}
    return nil
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

type ImageProcessingInfo struct {
	Url            string `json:"url" validate:"required,url"`
	Process        uint8  `json:"process" validate:"omitempty,number,oneof=0 1"`
	TempFileDir    string `json:"-"`  // constructed :: location to download temp files
	UploadPrefix   string `json:"-"`  // constructed :: eg. "media"
	VanityUrl      string `json:"-"`  // constructed
	AccountDir     string `json:"-"`  // constructed :: eg. "98c56d78fe3a"
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

type Image struct {
    // Image struct differs from base data, it creates the id
    // from Md5Hasher of account_id, url. Note: the image.Url is required
	Id       string `json:"id" validate:"omitempty"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
    Type string `json:"type" validate:"omitempty,lte=20"`
	AccountId string
	Document string
    AllIdNodes
    ImageProcessingInfo
}

type ImageNodes struct {
	Nodes []*Image `json:"imageNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *ImageNodes) Load(fileBuffer *[]byte) (err error){
	value := gjson.Get(string(*fileBuffer), "imageNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Image.Load %s", err)
	}
    return nil
}

func (s *ImageNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Image.Validate %s", err)
	}
    return nil
}

func (s *ImageNodes) PreProcess(accountId string, db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        // creates a unique id based on accountId, url
        node.Id = Md5Hasher([]string{accountId, node.Url})
        // constraints Image.Process == 0 if record already exists
        // so existing images won't be processed again
        exists, err := node.RecordValidate(node.Id, db)
        if exists == 1 {
            node.Process = 0
        }
        if err != nil {
            return fmt.Errorf("Image.PreProcess %s", err)
        }
    }
    return nil
}

func (s *ImageNodes) Upsert(accountId string, db *sql.DB) (err error) {
    // Record Original Image
    for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
        qstr := `
            INSERT INTO image (id, account_id, type, document)
            values ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE image.id = $1;`
        _, err = db.Exec(qstr, node.Id, accountId, node.Type, node.Document)
        if err != nil {
            return fmt.Errorf("ImageNodes.Upsert %s", err)
        }
    }
	return nil
}

func (s *ImageNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		if node.ParentId == "" {
			continue
		}
		qstr := `
            UPDATE image
            SET parent_id = $2
            WHERE image.id = $1;`
		_, err = db.Exec(qstr, node.Id, node.ParentId)
		if err != nil {
			return fmt.Errorf("ImageNodes.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *ImageNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        ascendentColumn := "image_id"
        structArray := []Upserter{}
        if node.CollectionIdNodes.Nodes != nil {
            node.collectionJson = s.Gjson.Array()[i].Get("collectionIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        if node.ContentIdNodes.Nodes != nil {
            node.contentJson = s.Gjson.Array()[i].Get("contentIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        //if node.ImageIdNodes.Nodes != nil {
            //node.imageJson = s.Gjson.Array()[i].Get("imageIdNodes")
            //node.ImageIdNodes.ascendentColumn = ascendentColumn
            //node.ImageIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.ImageIdNodes)
        //}
        if node.ItemIdNodes.Nodes != nil {
            node.itemJson = s.Gjson.Array()[i].Get("itemIdNodes")
            node.ItemIdNodes.ascendentColumn = ascendentColumn
            node.ItemIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ItemIdNodes)
        }
        if node.PlaceIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("placeIdNodes")
            node.PlaceIdNodes.ascendentColumn = ascendentColumn
            node.PlaceIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PlaceIdNodes)
        }
        if node.PersonIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("personIdNodes")
            node.PersonIdNodes.ascendentColumn = ascendentColumn
            node.PersonIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PersonIdNodes)
        }
        for _, sa := range structArray {
            err = UpsertHandler(sa, accountId, db)
            if err != nil {
                return fmt.Errorf("ImageNodes.RelatedTableUpsert %s", err)
            }
        }
    }
	return nil
}

func (s *Image) RecordValidate(record_id string, db *sql.DB) (exists int8, err error) {
    // checks if Image.Url exists in Image.Doc
    qstr := `
        SELECT COUNT(id)
        FROM image
        WHERE id = $1;
        `
    err = db.QueryRow(qstr, record_id).Scan(&exists)
    if err != nil {
        return 0, fmt.Errorf("Image.RecordValidate 01 %s",err)
    } 
    return exists, nil
}

func (s *Image) Construct(accountId string) {
    // adds additional information needed to process resized images

    // create sizes
    lgSize := ImgSize{1.0, "LG"}
    mdSize := ImgSize{0.5, "MD"}
    smSize := ImgSize{0.25, "SM"}
    s.ImgSizes = append(
        s.ImgSizes,
        lgSize,
        mdSize,
        smSize,
    )

    // create date string format "2022-06-02"
    t := time.Now()
    if s.Date == "" {
        s.Date = t.Format("2006-01-02")
    }

    // create it
    if s.AccountDir == "" {
        // take it from the end of the accountId
        s.AccountDir = accountId[len(accountId)-12:]
    }

    // construct additional information needed
    s.TempFileDir = os.Getenv("TMPDIR")
    s.UploadPrefix = os.Getenv("ULOADP")
    s.DoCacheControl = "max-age=2592000" // one month
    s.DoContentType = "image/webp"
    s.DoBucket = os.Getenv("DOBCKT")
    s.DoEndpointUrl = os.Getenv("DOENDU")
    s.DoAccessKey = os.Getenv("DOAKEY")
    s.DoSecret = os.Getenv("DOSECR")
    s.DoRegionName = os.Getenv("DOREGN")
    s.VanityUrl = os.Getenv("DOVANU")
}

func (s *Image) Download() (err error) {
    isValid := ValidateUrl(s.Url)
    if isValid == false {
        return fmt.Errorf("Image url is not valid %s", s.Url)
    }
    baseFileName, fullFileName := GetFileName(s.Url)
    s.BaseFileName = baseFileName
    s.filePath = filepath.Join(
        os.Getenv("TMPDIR"),
        "images/downloads",
        fullFileName,
    )
    if err != nil {
        return fmt.Errorf("Image.Download 01 %s", err)
    }
    err = DownloadFile(s.Url, s.filePath)
    if err != nil {
        return fmt.Errorf("Image.Download 02 %s", err)
    }
	return nil
}


func (s *Image) Resize() (err error) {
    // Check that download file exists
    _, err = os.Stat(s.filePath)
    if err != nil {
        return fmt.Errorf("Download file does not exist: %s", err)
    }

    for _, size := range s.ImgSizes {
        imgFile, err := os.Open(s.filePath)
        if err != nil {
            return fmt.Errorf("Resize open error 01: %s", err)
        }
        defer imgFile.Close()
        imgConfig, _, err := image.DecodeConfig(imgFile)
        imgWidth := imgConfig.Width
        imgHeight := imgConfig.Height
        imgFile, err = os.Open(s.filePath)
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
        newFileName := CreateNewFileName(s.BaseFileName, newWidth, newHeight)
        tempFilePath := filepath.Join(
            os.Getenv("TMPDIR"),
            "images/resized",
            newFileName,
        )

        uploadPath := CreateUploadPath(s.UploadPrefix,
            s.AccountDir, newFileName, s.Date)
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
        rsi.url = filepath.Join(s.VanityUrl, uploadPath)
        rsi.width = fmt.Sprint(newWidth)
        rsi.height = fmt.Sprint(newHeight)
        rsi.size = size.size
        s.ResizedImages = append(s.ResizedImages, rsi)

    }
	return nil
}

func (s *Image) ResizedImageUpsert(accountId string, db *sql.DB) (err error) {
    if s.ResizedImages == nil {
        return errors.New("ResizedImageUpsert needs Image.ResizedImages")
    }
    for _, i := range s.ResizedImages {


        // creates unique id on url, size
        idHash := Md5Hasher([]string{accountId, i.url})
        documentMap := map[string]interface{}{
            "url": i.url, 
            "height": i.height,
            "width": i.width,
            "size": i.size,
        }
        documentJson, _ := json.Marshal(documentMap)

        qstr := `
            INSERT INTO image (id, account_id, parent_id, type, document)
            values ($1, $2, $3, $4, $5)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                parent_id = $3,
                type = $4,
                document = $5
            WHERE image.id = $1;`
        _, err = db.Exec(
            qstr, idHash, accountId, s.Id, s.Type, string(documentJson))
        if err != nil {
            return fmt.Errorf("Image.Upsert %s", err)
        }
    }
    return nil
}


func (s *Image) UploadToSpaces() (err error) {
	//img {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]
    customResolver := aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{URL: s.DoEndpointUrl}, nil
        },
    )

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithEndpointResolverWithOptions(customResolver),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                s.DoAccessKey, s.DoSecret, "")),
    )

    if err != nil {
        panic("AWS configuration error, " + err.Error())
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.Region = s.DoRegionName
    })

    for _, rsi := range s.ResizedImages {
        file, err := os.Open(rsi.tempFilePath)
        if err != nil {
            return fmt.Errorf("Image.UploadToSpaces 01 %s %s", err, rsi.tempFilePath)
        }

        defer file.Close()

        input := &s3.PutObjectInput{
            Bucket:       &s.DoBucket,
            Key:          &rsi.key, // <-- uploadPath
            Body:         file,
            CacheControl: &s.DoCacheControl,
            ContentType:  &s.DoContentType,
            ACL:          "public-read",
        }

        _, err = PutFile(context.TODO(), client, input)
        if err != nil {
            return fmt.Errorf("Image.UploadToSpaces 02 %s", err)
        }
    }
	return nil
}

func (s *Image) RemoveTempFile() (err error) {
    e := os.Remove(s.filePath)
    if e != nil {
        return fmt.Errorf("Image.RemoveTempFile %s", err)
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
			return fmt.Errorf("DownlaodFile 01 %s", err)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
        return fmt.Errorf("DownlaodFile 02 %s", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
        return fmt.Errorf("DownlaodFile 03 %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
        return fmt.Errorf("DownlaodFile 04 %s", err)
	}

	return nil
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
	uploadPrefix, accountDir, fileName, date string) string {
	//CreateUploadPath creates a new path based on uploadPrefix eg. "media"
	//accountDir eg. "d91216fbb4d4" and fileName.
	//Path fomat /<uploadPrefix>/<accountPath>/year/month/data/filename
	//Date should be in string format "2022-06-02"

	//year, month, day := time.Now().Date()
	dateSlice := strings.Split(date, "-")
	y := dateSlice[0]
	m := dateSlice[1]
	d := dateSlice[2]
	key := filepath.Join(uploadPrefix, accountDir, y, m, d, fileName)
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
