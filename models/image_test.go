package models

import (
	"github.com/robertsmoto/skustor/configs"
	"log"
	"os"
	"testing"
)

func TestValidateUrl(t *testing.T) {
	//log.Print("TestValidateUrl ....")
	var url string
	var isValid bool

	// Valid url should validate as true
	url = "http://go.dev/blog/go-brand/logos.jpg"
	isValid = ValidateUrl(url)
	if isValid == false {
		t.Error("Url is valid yet vaidates as false")
	}
	// filetype must be of jpg or png
	url = "http://go.dev/blog/go-brand/logos.webp"
	isValid = ValidateUrl(url)
	if isValid == true {
		t.Error("Filetype is of wrong type yet vaidates as true")
	}
	// url must begin with "http"
	url = "//go.dev/blog/go-brand/logos.jpg"
	isValid = ValidateUrl(url)
	if isValid == true {
		t.Error("Url does not begin with http yet vaidates as true")
	}
}

func TestGetFileName(t *testing.T) {
	url := "http://go.dev/blog/go-brand/logos.jpg"
	baseFileName, fullFileName := GetFileName(url)
	if fullFileName != "logos.jpg" || baseFileName != "logos" {
		t.Error("Error parsing string in func GetFileName")
	}
}

func TestDownloadFile(t *testing.T) {
	url := "http://go.dev/blog/go-brand/logos.jpg"
	filePath := "./test_data/logos.jpg"
	err := DownloadFile(url, filePath)
	if err != nil {
		log.Printf("Error downloading file: %s", err)
		t.Error("TestDownlodFile Error")
	}

	// Check that file exists.
	_, err = os.Stat(filePath)
	if err != nil {
		t.Errorf("Error downloading image: %s", err)
	}
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error removing image from test directory: %s", err)
	}
}

func TestCreateNewFilename(t *testing.T) {
	baseFileName := "prettyPicture"
	w := 800
	h := 400

	newName := CreateNewFileName(baseFileName, w, h)
	if newName != "prettyPicture_800x400.webp" {
		t.Error("New image filename was not created correctly.")
	}
}

func TestCreateTempFilePath(t *testing.T) {

	tempFileDir := "/tempdir"
	fileName := "prettyPicture_800x400.webp"
	newFilePath := CreateTempFilePath(tempFileDir, fileName)
	if newFilePath != "/tempdir/prettyPicture_800x400.webp" {
		t.Error("New file path was not created correctly.")
	}
}

func TestCreateUploadPath(t *testing.T) {

	uploadPrefix := "media"
	userDir := "111111111111"
	date := "2022-05-15"

	fileName := "prettyPicture_800x400.webp"
	uploadPath := CreateUploadPath(uploadPrefix, userDir, fileName, date)
	if uploadPath != "media/111111111111/2022/05/15/prettyPicture_800x400.webp" {
		t.Errorf("New upload path was not created correctly. %s", uploadPath)
	}
}

func TestCalcNewSize(t *testing.T) {

	newWidth, newHeight := CalcNewSize(800, 400, 0.25)
	if newWidth != 200 || newHeight != 100 {
		t.Errorf(
			"New width and height were not created correctly. %d, %d",
			newWidth, newHeight)
	}
}

func TestWebImage(t *testing.T) {

	conf := configs.Config{}
	configs.Load(&conf)

	// instantiate the WebImage struct and assign variables
	i := Image{}

	lgSize := ImgSize{1.0, "LG"}
	mdSize := ImgSize{0.5, "MD"}
	smSize := ImgSize{0.25, "SM"}
	i.ImgSizes = append(
		i.ImgSizes,
		lgSize,
		mdSize,
		smSize,
	)

	i.Url = "https://cdn-stage.sodavault.com/media/111111111111/svLogo.png"
	//i.Url = "http://go.dev/blog/go-brand/logos.jpg"
	//i.Url = "https://cdn-stage.sodavault.com/media/111111111111/testTree.jpg"
	i.TempFileDir = "./test_data"
	i.UploadPrefix = "media"
	i.UserDir = "111111111111"
	i.Date = "2022-06-01"
	i.DoCacheControl = "max-age=2592000" // one month
	i.DoContentType = "image/webp"
	i.DoBucket = conf.DoSpaces.BucketName
	i.DoEndpointUrl = conf.DoSpaces.EndpointUrl
	i.DoAccessKey = conf.DoSpaces.AccessKey
	i.DoSecret = conf.DoSpaces.Secret
	i.DoRegionName = conf.DoSpaces.RegionName
	i.VanityUrl = conf.DoSpaces.VanityUrl

	err := i.Download()
	if err != nil {
		log.Print(err)
	}
	_, err = os.Stat(i.filePath)
	if err != nil {
		t.Errorf("Error downloading image: %s", err)
	}
	err = i.Resize()
	if err != nil {
		log.Print(err)
	}

	//lgImg {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]
	for _, rsi := range i.ResizedImages {
		_, err = os.Stat(rsi.tempFilePath)

		if err != nil {
			t.Errorf("Error making new image size: %s", err)
		}
	}

	// upload to cdn
	err = i.UploadToSpaces()
	if err != nil {
		log.Print(err)
	}

	// use newSizes information to record entries in the db
}
