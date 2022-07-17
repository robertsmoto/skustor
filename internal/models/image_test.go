package models

import (
    //"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/robertsmoto/skustor/internal/configs"
	"github.com/robertsmoto/skustor/internal/postgres"
	"github.com/tidwall/gjson"
)

func SetEnvConfig() (err error) {
	configs.Load(&configs.Config{})
	return nil
}

func Test_CheckProcessed(t *testing.T) {
    // check manually in db to see if record was created

    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
    
    i := Image{}
    i.Url = "http://www.example.com/test01.jpg"
    i.Id = Md5Hasher([]string{accountId, i.Url})
    i.Type = "testImage"
    i.Process = 0
    
    imageNodes := ImageNodes{}
    // need to mock the doc
    request := `{"url": "http://www.example.com/test01.jpg"}`
    imageNodes.Gjson = gjson.Parse(request)

    imageNodes.Nodes = append(imageNodes.Nodes, &i)
    devPostgres := postgres.PostgresDb{}
    pgDb, err := postgres.Open(&devPostgres)
    defer pgDb.Close()
    err = imageNodes.Upsert(accountId, pgDb)
    if err != nil {
        t.Error("Db error.", err)
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

func Test_ExistingUrlCheck(t *testing.T) {
    var exists int8
    var err error
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
    // open database
    db, err := postgres.Open(&postgres.PostgresDb{})
    if err != nil {
        t.Error("Test_ExistingUrlCheck 01", err)
    }
    // check does not exist
    i1 := Image{}
    i1.Url = "http://www.example.com"
    i1.Id = Md5Hasher([]string{accountId, i1.Url})
    exists, err = i1.RecordValidate(i1.Id, db)
    if exists != 0 {
        t.Errorf("Test_ExistingUrlCheck exists %d != 0", exists)
    }
    // check exists
    i2 := Image{}
    i2.Url = "http://www.example.com/test01.jpg"
    i2.Id = Md5Hasher([]string{accountId, i2.Url})
    exists, err = i2.RecordValidate(i2.Id, db)
    if exists != 1 {
        t.Errorf("Test_ExistingUrlCheck exists %d != 1", exists)
    }
    // close db
    db.Close()
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

func Test_WebImage(t *testing.T) {

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
    i.Type = "testImage"

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    iNodes := ImageNodes{}

    // need to mock the request
    request := `{"mock": "for testing"}`
    iNodes.Gjson = gjson.Parse(request)

    iNodes.Nodes = append(iNodes.Nodes, &i)
    var err error

    // open database
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // preProcess imageNodes
    err = iNodes.PreProcess(accountId, pgDb)
    if err != nil {
        t.Error("Test_WebImage 01 ", err)
    }

    // upsert imageNodes
    err = iNodes.Upsert(accountId, pgDb)
    if err != nil {
        t.Error("Test_WebImage 01 ", err)
    }

    // work on Image.Process == 1
    for _, node := range iNodes.Nodes {
        if node.Process == 0 {
            continue
        }

        err = node.Download()
        if err != nil {
            t.Error("Download error: ", err)
        }

        err = node.Resize()
        if err != nil {
            t.Error("Resize image ", err)
        }

        // upsert new sizes to db
        err = node.ResizedImageUpsert(accountId, pgDb)
        if err != nil {
            t.Errorf("Upsert %s", err)
        }

        // upload to cdn
        err = node.UploadToSpaces()
        if err != nil {
            t.Errorf("Upload to spaces %s", err)
        }

        // delete temp file
        err = node.RemoveTempFile()
        if err != nil {
            t.Errorf("Upload to spaces %s", err)
        }
    }
    pgDb.Close()
}

func Test_ImageInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    configs.Load(&configs.Config{})
    if err != nil {
        t.Errorf("Test_ImageInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/images.json")
    if err != nil {
        t.Errorf("Test_ImagesInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    cNodes := ImageNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&cNodes, &testFile)
    if err != nil {
        t.Errorf("Test_ImageInterfaces 01 %s", err)
    }

    // special preProcessing
    err = PreProcessHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ImageInterfaces 02 %s", err)
    }

    err = UpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ImageInterfaces 03 %s", err)
    }

    err = ForeignKeyUpdateHandler(&cNodes, pgDb)
    if err != nil {
        t.Errorf("Test_ImageInterfaces 04 %s", err)
    }

    err = RelatedTableUpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ImageInterfaces 05 %s", err)
    }

    // works on image.ResizedImages if image.Process == 1
    for _, node := range cNodes.Nodes {

        node.Date = "2022-06-01"
        node.AccountDir = "111111111111"
        if node.Process == 0 {
            continue
        }

        err = ResizedImgHandler(node, accountId, pgDb)
        if err != nil {
            t.Errorf("Test_ImageInterfaces 05 %s", err)
        }
    }
}

