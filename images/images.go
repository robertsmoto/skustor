package images

import (
	//"errors"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	_ "image/jpeg"
    _ "image/png"

	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
)

type Downloader interface {
	Download() (err error)
}

type Sizer interface {
	Size() (err error)
}

type Uploader interface {
    UploadImagesToSpaces() (err error)
}

type DownloaderSizerUploader interface {
	Downloader
	Sizer
    Uploader
}

func Resize(i DownloaderSizerUploader) (err error) {
	// should return []strings
	err = i.Download()
	err = i.Size()
    err = i.UploadImagesToSpaces()  // aws.go
	return err
}

type WebImage struct {
	Url            string // constructed
	TempFileDir    string // constructed :: location to download temp files
	UploadPrefix   string // constructed :: eg. "media"
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
	NewSizes       map[string][]string //localFilePath[uploadFilePath, height, width]
}

func (i *WebImage) Download() (err error) {
	//Downloads image from url to tempFileDir
	// validates url and filetype
	isValid := ValidateUrl(i.Url)
	if isValid == false {
		log.Fatal("Image url or file-type is not valid.")
		os.Exit(1)
	}

	i.baseFileName, i.fullFileName = GetFileName(i.Url)

	i.filePath = CreateTempFilePath(i.TempFileDir, i.fullFileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = DownloadFile(i.Url, i.filePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return err
}

func (i *WebImage) Size() (err error) {
	// Resizes, renames and encodes original image.

	fmt.Println(i.filePath)
	imgFile, err := os.Open(i.filePath)
	if err != nil {
		fmt.Println("error here")
		log.Fatal(err)
		os.Exit(1)
	}
	defer imgFile.Close()

	imgConfig, _, err := image.DecodeConfig(imgFile)
	imgWidth := imgConfig.Width
	imgHeight := imgConfig.Height

	imgFile, err = os.Open(i.filePath)
	if err != nil {
		fmt.Println("error here")
		log.Fatal(err)
		os.Exit(1)
	}
	defer imgFile.Close()

	decodedImage, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("error here")
		log.Fatal(err)
	}

	sizeMap := map[int]float64{
		1: 1.0,
		2: 0.5,
		3: 0.25,
	}

	// resize images in various sizes
	for _, v := range sizeMap {

		// calculate new image sizes
		newWidth, newHeight := CalcNewSize(imgWidth, imgHeight, v)
		fmt.Println(newWidth, newHeight)

		// create new file name and dirs
		newFileName := CreateNewFileName(i.baseFileName, newWidth, newHeight)
		fmt.Println(newFileName)

		tempFilePath := CreateTempFilePath(i.TempFileDir, newFileName)

		// assign values to the newSizes map
		if i.NewSizes == nil {
			i.NewSizes = make(map[string][]string)
		}
		i.NewSizes[tempFilePath] = append(
			i.NewSizes[tempFilePath],
			CreateUploadPath(i.UploadPrefix, i.UserDir, newFileName, i.Date),
		)
		i.NewSizes[tempFilePath] = append(
			i.NewSizes[tempFilePath],
			fmt.Sprint(newWidth),
		)
		i.NewSizes[tempFilePath] = append(
			i.NewSizes[tempFilePath],
			fmt.Sprint(newHeight),
		)
		// create local dir
		// if successful, the created file can be used for I/O
		var f *os.File
		f, err = os.Create(tempFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// resize the image in memory
		// resizedImage is of image.Image type
		resizedImage := ResizeBaseImage(decodedImage, imgWidth, imgHeight, v)

		// encode image
		// Encode takes two arguments, io.writer and image.Image
		err = webpbin.Encode(f, resizedImage)
		if err != nil {
			f.Close()
			log.Fatal(err)
		}
        f.Close()
	}
	return
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

func ResizeBaseImage(decodedImg image.Image, w, h int, sizeFactor float64) image.Image {
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

func CalcNewSize(width, height int, sizeFactor float64) (int, int) {
	newWidth := int(math.Round(float64(width) * sizeFactor))
	newHeight := int(math.Round(float64(height) * sizeFactor))
	return newWidth, newHeight
}
