package images

import (
	"log"
	"os"
	"testing"

	"github.com/robertsmoto/skustor/configs"
)

func TestValidateUrl(t *testing.T) {
	log.Print("TestValidateUrl ....")
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
	log.Print("TestGetFileName ....")

	url := "http://go.dev/blog/go-brand/logos.jpg"
	baseFileName, fullFileName := GetFileName(url)
	if fullFileName != "logos.jpg" || baseFileName != "logos" {
		t.Error("Error parsing string in func GetFileName")
	}
}

func TestDownloadFile(t *testing.T) {
	// Test successful download.
	log.Print("TestDownloadFile ....")
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
	} else {
		log.Print("Successfully downloaded image ...")
	}
	// Remove file.
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error removing image from test directory: %s", err)
	} else {
		log.Print("Successfully removed image from test directory ...")
	}
}

func TestCreateNewFilename(t *testing.T) {
	log.Print("TestCreateNewFilename ....")
	baseFileName := "prettyPicture"
	w := 800
	h := 400

	newName := CreateNewFileName(baseFileName, w, h)
	if newName != "prettyPicture_800x400.webp" {
		t.Error("New image filename was not created correctly.")
	}
}

func TestCreateTempFilePath(t *testing.T) {
	log.Print("TestCreateTempFilePath ....")

	tempFileDir := "/tempdir"
	fileName := "prettyPicture_800x400.webp"
	newFilePath := CreateTempFilePath(tempFileDir, fileName)
	if newFilePath != "/tempdir/prettyPicture_800x400.webp" {
		t.Error("New file path was not created correctly.")
	}
}

func TestCreateUploadPath(t *testing.T) {
	log.Print("TestCreateUploadPath ....")

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
	log.Print("TestCalcNewSize ....")

	width := 800
	height := 400
	sizeFactor := .25
	newWidth, newHeight := CalcNewSize(width, height, sizeFactor)
	if newWidth != 200 || newHeight != 100 {
		t.Errorf(
			"New width and height were not created correctly. %d, %d",
			newWidth, newHeight)
	}
}

func TestWebImage(t *testing.T) {
	// This puts all the webimage processing functions together
	log.Print("TestWebImageImplementation ....")

	log.Print("LoadingSvConf ....")
	conf := configs.Config{}
	configs.Load(&conf)

	log.Print("Instantiating WebImage struct ....")
	// instantiate the WebImage struct and assign variables
	i := WebImage{}
	i.Url = "http://go.dev/blog/go-brand/logos.jpg"
	i.TempFileDir = "./test_data"
	i.UploadPrefix = "media"

	i.UserDir = "111111111111"
	i.Date = "2022-06-01"
	i.DoCacheControl = "max-age=2592000" // one month
	i.DoContentType = "image/webp"
	// these values will be assigned from conf values
	i.DoBucket = conf.DoSpaces.BucketName
	i.DoEndpointUrl = conf.DoSpaces.EndpointUrl
	i.DoAccessKey = conf.DoSpaces.AccessKey
	i.DoSecret = conf.DoSpaces.Secret
	i.DoRegionName = conf.DoSpaces.RegionName

	log.Print("WebImage.Download() ....")
	err := i.Download()
	if err != nil {
		log.Print(err)
	}
	_, err = os.Stat(i.filePath)
	if err != nil {
		t.Errorf("Error downloading image: %s", err)
	} else {
		log.Print("Successfully downloaded image ...", i.filePath)
	}

	log.Print("WebImage.MakeNewSizes() ....")
	err = i.Size()
	if err != nil {
		log.Print(err)
	}

	// newSizes is map localDir: [UploadDir, newHeight, newWidth]
	log.Print("##New Sizes --> ", i.NewSizes)
	for s := range i.NewSizes {
		_, err = os.Stat(s)
		if err != nil {
			t.Errorf("Error making new image size: %s", err)
		} else {
			log.Print("Successfully created new image size ...", s)
		}
	}

	// upload to cdn
	log.Print("WebImage.UploadImagesToSpaces() ....")
	err = i.UploadImagesToSpaces()
	if err != nil {
		log.Print(err)
	}
}
