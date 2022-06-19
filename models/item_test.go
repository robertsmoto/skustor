package models

import (
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/robertsmoto/skustor/configs"
	"github.com/robertsmoto/skustor/tools"
)

func TestItem_Upsert(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/items.json")
	if err != nil {
		t.Errorf("TestItem_Upsert %s", err)
		log.Print("error", err)
	}
	items := ItemNodes{}
	// open the db
	devPostgres := tools.PostgresDev{}
	devDb, err := tools.Open(&devPostgres)
	if err != nil {
		t.Error("TstItem_Upsert Open() ", err)
	}

	// test the JsonLoadValidateUpsert interface
	// this loads the []structs
	LoaderHandler(&items, testFile)

	// now itterate over each struct
	userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	date := "2022-06-01"
	userDir := "111111111111"

	for _, item := range items.Nodes {
		log.Print("## item_test item --> ", item)
		item.UserId = userId
		UpsertHandler(&item, devDb)

		// check for image nodes
		if item.ImageNodes.Nodes != nil {

			conf := configs.Config{}
			configs.Load(&conf)

			lgSize := ImgSize{1.0, "LG"}
			mdSize := ImgSize{0.5, "MD"}
			smSize := ImgSize{0.25, "SM"}

			for _, imgNode := range item.ImageNodes.Nodes {
				// construct imgNode static data
				imgNode.ImgSizes = append(
					imgNode.ImgSizes,
					lgSize,
					mdSize,
					smSize,
				)

				imgNode.UserId = item.UserId // make sure to add ids
				imgNode.ItemId = item.Id
				imgNode.TempFileDir = conf.TempFileDir
				imgNode.UploadPrefix = conf.UploadPrefix
				imgNode.VanityUrl = conf.DoSpaces.VanityUrl
				imgNode.UserDir = userDir
				imgNode.Date = date
				imgNode.DoBucket = conf.DoSpaces.BucketName
				imgNode.DoCacheControl = "max-age=2628002" // one month
				imgNode.DoContentType = "image/webp"
				imgNode.DoEndpointUrl = conf.DoSpaces.EndpointUrl
				imgNode.DoAccessKey = conf.DoSpaces.AccessKey
				imgNode.DoSecret = conf.DoSpaces.Secret
				imgNode.DoRegionName = conf.DoSpaces.RegionName

				ImgHandler(&imgNode, devDb)
				if err != nil {
					t.Error("Failed at ImageSizeUpsert", err)
				}
			}
		}
	}
}

func TestItem_Validate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/items.json")
	if err != nil {
		t.Errorf("TestItem Reading json file. %s", err)
		log.Print("error", err)
	}
	items := ItemNodes{}
	// test the loader and validator
	err = items.Load(&testFile)
	err = items.Validate()
	if err != nil {
		t.Error("Item Load error", err)
	}
	node0 := items.Nodes[0]
	if node0.Id == "a4ce7648-e499-40d1-8e6b-0c2ef0a7856e" != true {
		t.Error("Error loading Items")
	}
	if err != nil {
		t.Error("TestGroupValidate Error: ", err)
	}
}
