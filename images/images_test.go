package images

import (
	"fmt"
	"os"
	"testing"

    "github.com/robertsmoto/skustor/internal/conf"
)

func TestValidateUrl(t *testing.T) {
	fmt.Println("TestValidateUrl ....")
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
	fmt.Println("TestGetFileName ....")

    url := "http://go.dev/blog/go-brand/logos.jpg"
	baseFileName, fullFileName := GetFileName(url)
	if fullFileName != "logos.jpg" || baseFileName != "logos" {
		t.Error("Error parsing string in func GetFileName")
	}
}

func TestDownloadFile(t *testing.T) {
	// Test successful download.
	fmt.Println("TestDownloadFile ....")
	url := "http://go.dev/blog/go-brand/logos.jpg"
	filePath := "./test_data/logos.jpg"
	err := DownloadFile(url, filePath)
	if err != nil {
		fmt.Printf("Error downloading file: %s", err)
		t.Error("TestDownlodFile Error")
	}

	// Check that file exists.
	_, err = os.Stat(filePath)
	if err != nil {
		t.Errorf("Error downloading image: %s", err)
	} else {
		fmt.Println("Successfully downloaded image ...")
	}
	// Remove file.
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error removing image from test directory: %s", err)
	} else {
		fmt.Println("Successfully removed image from test directory ...")
	}
}

func TestCreateNewFilename(t *testing.T) {
	fmt.Println("TestCreateNewFilename ....")
	baseFileName := "prettyPicture"
	w := 800
	h := 400

	newName := CreateNewFileName(baseFileName, w, h)
	if newName != "prettyPicture_800x400.webp" {
		t.Error("New image filename was not created correctly.")
	}
}

func TestCreateTempFilePath(t *testing.T) {
	fmt.Println("TestCreateTempFilePath ....")

	tempFileDir := "/tempdir"
	fileName := "prettyPicture_800x400.webp"
	newFilePath := CreateTempFilePath(tempFileDir, fileName)
	if newFilePath != "/tempdir/prettyPicture_800x400.webp" {
		t.Error("New file path was not created correctly.")
	}
}

func TestCreateUploadPath(t *testing.T) {
	fmt.Println("TestCreateUploadPath ....")

	uploadPrefix := "media"
	userDir := "ed59b8fa8cb3"
	date := "2022-05-15"

	fileName := "prettyPicture_800x400.webp"
	uploadPath := CreateUploadPath(uploadPrefix, userDir, fileName, date)
	if uploadPath != "media/ed59b8fa8cb3/2022/05/15/prettyPicture_800x400.webp" {
		t.Errorf("New upload path was not created correctly. %s", uploadPath)
	}
}

func TestCalcNewSize(t *testing.T) {
	fmt.Println("TestCalcNewSize ....")

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
	fmt.Println("TestWebImageImplementation ....")


	fmt.Println("LoadingSvConf ....")
    SvConf := conf.Config{}
    err := SvConf.LoadJson("../internal/conf/test_data/test_config.json")
	if err != nil {
		t.Error(err)
	}
    // output SvConf
	fmt.Println(SvConf)

	fmt.Println("Instantiating WebImage struct ....")
    // instantiate the WebImage struct and assign variables
    i := WebImage{}
    i.url = "http://go.dev/blog/go-brand/logos.jpg"
    i.tempFileDir = "./test_data"
    i.uploadPrefix = "media"
    i.userDir = "test-user-id"
    i.date = "2022-06-01"
	i.doCacheControl = "max-age=2592000"  // one month
	i.doContentType = "image/webp"
    // these values will be assigned from conf values
    i.doBucket = SvConf.DoSpaces.BucketName
    i.doEndpointUrl = SvConf.DoSpaces.EndpointUrl
    i.doAccessKey = SvConf.DoSpaces.AccessKey
    i.doSecret = SvConf.DoSpaces.Secret
    i.doRegionName = SvConf.DoSpaces.RegionName
    
	fmt.Println("WebImage.Download() ....")
    err = i.Download()
    if err != nil {
        fmt.Println(err)
    }
    _, err = os.Stat(i.filePath)
	if err != nil {
		t.Errorf("Error downloading image: %s", err)
	} else {
		fmt.Println("Successfully downloaded image ...", i.filePath)
	}

	fmt.Println("WebImage.MakeNewSizes() ....")
    err = i.MakeNewSizes()
    if err != nil {
        fmt.Println(err)
    }

    // newSizes is map localDir, UploadDir
    for s := range(i.newSizes){
        _, err = os.Stat(s)
        if err != nil {
            t.Errorf("Error making new image size: %s", err)
        } else {
            fmt.Println("Successfully created new image size ...", s)
        }
    }

    // upload to cdn
	fmt.Println("WebImage.UploadImagesToSpaces() ....")
    err = i.UploadImagesToSpaces()
    if err != nil {
        fmt.Println(err)
    }
}
