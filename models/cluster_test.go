package models

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	//"github.com/pborman/uuid"
	"github.com/robertsmoto/skustor/configs"
	"github.com/robertsmoto/skustor/tools"
)

func TestGroup_Upsert(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/groups.json")
	if err != nil {
		t.Error("TestGroup_Upsert() ", err)
	}
	groups := GroupNodes{}
	// open the db
	devPostgres := tools.PostgresDev{}
	devDb, err := tools.Open(&devPostgres)

	// test the JsonLoadValidateUpsert interface
	// this loads the []structs
	LoaderHandler(&groups, testFile)

	// now loop through eash struct indvidually
	userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	date := "2022-06-01"
	userDir := "111111111111"

	for _, group := range groups.Nodes {
		group.UserId = userId
		UpsertHandler(&group, devDb)
		if err != nil {
			t.Error("TestGroup_Upsert.UpsertHandler() ", err)
		}

		// check for image nodes
		if group.ImageNodes.Nodes != nil {

			conf := configs.Config{}
			configs.Load(&conf)

			lgSize := ImgSize{1.0, "LG"}
			mdSize := ImgSize{0.5, "MD"}
			smSize := ImgSize{0.25, "SM"}

			for _, imgNode := range group.ImageNodes.Nodes {
				// construct imgNode static data
				imgNode.ImgSizes = append(
					imgNode.ImgSizes,
					lgSize,
					mdSize,
					smSize,
				)

				imgNode.UserId = group.UserId // make sure to add ids
				imgNode.GroupId = group.Id    // make sure to add ids
				imgNode.ItemId = ""
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
			}
		}
	}

	// test the Image interface
}

func TestGroup_Validate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/groups.json")
	if err != nil {
		t.Errorf("TestBrands %s", err)
	}
	groups := GroupNodes{}
	// test the loader and validator
	err = groups.Load(&testFile)
	err = groups.Validate()
	if err != nil {
		t.Error("Brand Loader error", err)
	}
	node0 := groups.Nodes[0]
	if node0.Id == "eeb75266-7f4a-4d8e-9a8a-2c0ada73e7b1" != true {
		t.Error("Error loading Brands")
	}
	if err != nil {
		t.Error("TestGroupValidate Error: ", err)
	}
}
