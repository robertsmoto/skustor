package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/robertsmoto/skustor/internal/configs"
	"github.com/robertsmoto/skustor/internal/postgres"
)

func SetEnvConfig() (err error) {
	conf := configs.Config{}
	configs.Load(&conf)
	return err
}

func Test_CheckProcessed(t *testing.T) {
    // check manually in db to see if record was created
    imageNodes := ImageNodes{}
    i := Image{
        Url: "https://www.example.com/test_processed.jpg",
        Process: 0,
    }
    imageNodes.Nodes = append(imageNodes.Nodes, &i)
	devPostgres := postgres.PostgresDb{}
	pgDb, err := postgres.Open(&devPostgres)
    defer pgDb.Close()
    imageNodes.RecordOriginalImage(pgDb)
    if err != nil {
		t.Error("Db error.")
    }
}

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
	_ = SetEnvConfig()

	url := "http://go.dev/blog/go-brand/logos.jpg"
	filePath := filepath.Join(
		os.Getenv("TMPDIR"), "images/downloads", "logos.jpg")
	err := DownloadFile(url, filePath)
	if err != nil {
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

	//tempFileDir := "/tempdir"
	fileName := "prettyPicture_800x400.webp"
	newFilePath := filepath.Join(os.Getenv("TMPDIR"), "images/downloaded", fileName)
	if newFilePath != "/home/robertsmoto/dev/temp/sodavault/images/downloaded/prettyPicture_800x400.webp" {
		t.Error("New file path was not created correctly.", newFilePath)
	}
}

func TestCreateUploadPath(t *testing.T) {

	accountDir := "111111111111"
	date := "2022-05-15"

	fileName := "prettyPicture_800x400.webp"
	uploadPath := CreateUploadPath(os.Getenv("ULOADP"), accountDir, fileName, date)
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

	// create the env variables
	_ = SetEnvConfig()

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
    i.Process = 1

	i.TempFileDir = os.Getenv("TMPDIR")
	i.UploadPrefix = os.Getenv("ULOADP")
	i.AccountDir = "111111111111"
	i.Date = "2022-06-01"
	i.DoCacheControl = "max-age=2592000" // one month
	i.DoContentType = "image/webp"
	i.DoBucket = os.Getenv("DOBCKT")
	i.DoEndpointUrl = os.Getenv("DOENDU")
	i.DoAccessKey = os.Getenv("DOAKEY")
	i.DoSecret = os.Getenv("DOSECR")
	i.DoRegionName = os.Getenv("DOREGN")
	i.VanityUrl = os.Getenv("DOVANU")


    iNodes := ImageNodes{}
    iNodes.Nodes = append(iNodes.Nodes, &i)
    var err error

    // check database
    devPostgres := postgres.PostgresDb{}
    pgDb, err := postgres.Open(&devPostgres)
    defer pgDb.Close()

    err = iNodes.RecordOriginalImage(pgDb)
    if err != nil {
        t.Error("CheckProcessed error: ", err)
    }

    err = iNodes.Download()
    if err != nil {
        t.Error("Download error: ", err)
    }
    //_, err = os.Stat(iNodes.filePath)
    //if err != nil {
        //t.Errorf("Download file does not exist: %s", err)
    //}
    err = iNodes.Resize()
    if err != nil {
        t.Error("Resize image ", err)
    }
    //lgImg {0=tempFilePath, 1=key 2=url, 3=width, 4=height, 5=size eg "LG"]
    for _, rsi := range i.ResizedImages {
        _, err = os.Stat(rsi.tempFilePath)

        if err != nil {
            t.Errorf("Error making new image size: %s", err)
        }
    }
    // upsert new sizes to db
    err = iNodes.Upsert(pgDb)
    if err != nil {
        t.Errorf("Upsert %s", err)
    }
    // upload to cdn
    err = iNodes.UploadToSpaces()
    if err != nil {
        t.Errorf("Upload to spaces %s", err)
    }
    // delete temp file
    err = iNodes.RemoveTempFile()
    if err != nil {
        t.Errorf("Upload to spaces %s", err)
    }
}


