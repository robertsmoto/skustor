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

	//"github.com/robertsmoto/skustor/configs"
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

func ImgHandler(i downloaderResizerProcessorUpserterUploader, db *sql.DB) {
	var err error
	err = i.Download()
	err = i.Resize()
	i.Process()
	err = i.Upsert(db)
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
	Title          string `json:"title" validate:"omitempty,lte=200"`
	Alt            string `json:"alt" validate:"omitempty,lte=100"`
	Caption        string `json:"caption" validate:"omitempty,lte=200"`
	Position       uint8  `json:"position" validate:"omitempty,number"`
	Bypass         uint8  `json:"bypassProcessing" validate:"omitempty,number,oneof=0 1"`
	Height         string `json:"height" validate:"omitempty,lte=20"`
	Width          string `json:"width" validate:"omitempty,lte=20"`
	Size           string `json:"size" validate:"omitempty,lte=20,oneof=LG MD SM"`
	SvUserId       string // constructed
	ClusterId      string // constructed
	ItemId         string // constructed
	TempFileDir    string // constructed :: location to download temp files
	UploadPrefix   string // constructed :: eg. "media"
	VanityUrl      string // constructed
	UserDir        string // constructed :: eg. "98c56d78fe3a"
	Date           string // constructed :: eg. "2020-05-04"
	DoBucket       string // constructed
	DoCacheControl string // constructed eg. "max-age=60"
	DoContentType  string // constructed eg. "image/webp"
	DoEndpointUrl  string // constructed
	DoAccessKey    string // constructed
	DoSecret       string // constructed
	DoRegionName   string // constructed
	baseFileName   string
	fullFileName   string
	filePath       string
	ImgSizes       []ImgSize
	ResizedImages  []ResizedImage
}

type ImageNodes struct {
	Nodes []Image `json:"imageNodes" validate:"dive"`
}

func (s *ImageNodes) Load(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

func (s *ImageNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	return err
}

func (s *Image) Download() (err error) {
	//Downloads image from url to tempFileDir
	// validates url and filetype
	isValid := ValidateUrl(s.Url)
	if isValid == false {
		log.Fatal("Image url is not valid. ", s.Url)
		os.Exit(1)
	}

	s.baseFileName, s.fullFileName = GetFileName(s.Url)

	s.filePath = CreateTempFilePath(s.TempFileDir, s.fullFileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = DownloadFile(s.Url, s.filePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return err
}

func (s *Image) Resize() (err error) {
	// Resizes, renames and encodes original image.
	if s.Bypass > 0 {
		return err
	}
	if s.ImgSizes == nil {
		return err
	}

	for _, size := range s.ImgSizes {

		imgFile, err := os.Open(s.filePath)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		defer imgFile.Close()

		imgConfig, _, err := image.DecodeConfig(imgFile)
		imgWidth := imgConfig.Width
		imgHeight := imgConfig.Height

		imgFile, err = os.Open(s.filePath)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		defer imgFile.Close()

		decodedImage, _, err := image.Decode(imgFile)
		if err != nil {
			log.Print(err)
		}

		// calculate new image sizes
		newWidth, newHeight := CalcNewSize(imgWidth, imgHeight, size.ratio)

		// create new file name and dirs
		newFileName := CreateNewFileName(s.baseFileName, newWidth, newHeight)

		tempFilePath := CreateTempFilePath(s.TempFileDir, newFileName)
		uploadPath := CreateUploadPath(s.UploadPrefix,
			s.UserDir, newFileName, s.Date)

		// create local dir
		// if successful, the created file can be used for I/O
		var f *os.File
		f, err = os.Create(tempFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// resize the image in memory
		// resizedImage is of image.Image type
		resizedImage := ResizeBaseImage(decodedImage, newWidth, newHeight)

		// encode image
		// Encode takes two arguments, io.writer and image.Image
		err = webpbin.Encode(f, resizedImage)
		if err != nil {
			f.Close()
			log.Fatal(err)
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

	return err
}

func (s *Image) Process() {
	var data []string
	data = append(data, s.SvUserId, s.ClusterId, s.ItemId, s.Size, string(s.Position))
	s.Id = Md5Hasher(data)
}

func (s *Image) Upsert(db *sql.DB) (err error) {

	// now add the resizeImg.NewSizes to the image table
	qstr := `
    INSERT INTO image (
        id,
        sv_user_id,
        cluster_id, 
        item_id,
        url,
        size,
        position,
        height,
        width,
        title,
        alt,
        caption
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    ON CONFLICT (id)
    DO UPDATE
    SET url = $5,
        size = $6,
        position = $7,
        height = $8,
        width = $9,
        title = $10,
        alt = $11,
        caption = $12
    WHERE
        image.id = $1;`

	// execute qstr
	//lgImg {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]

	if s.ResizedImages != nil {
		// loop through the new sizes
		for _, rsi := range s.ResizedImages {
			_, err = db.Exec(
				qstr,
				s.Id, FormatUUID(s.SvUserId), FormatUUID(s.ClusterId),
				FormatUUID(s.ItemId), rsi.url, rsi.size, s.Position, rsi.height,
				rsi.width, s.Title, s.Alt, s.Caption,
			)
		}
	}
	if s.Bypass > 0 {
		// record the bypass image
		_, err = db.Exec(
			qstr,
			s.Id, FormatUUID(s.SvUserId), FormatUUID(s.ClusterId),
			FormatUUID(s.ItemId), s.Url, s.Size, s.Position, s.Height, s.Width,
			s.Title, s.Alt, s.Caption,
		)
	}
	if err != nil {
		log.Print("\nCommiting image to db: ", err)
	}
	return err
}

func (s *Image) UploadToSpaces() (err error) {
	//img {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]

	if s.ResizedImages == nil {
		return err
	}
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
			log.Print("Unable to open file ", rsi.tempFilePath)
			return err
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
			fmt.Print(err)
			return err
		}
	}
	return err
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

func CreateTempFilePath(tempFileDir, fileName string) (tempPath string) {
	tempPath = filepath.Join(tempFileDir, fileName)
	return tempPath
}

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
