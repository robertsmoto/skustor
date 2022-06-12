package models

import (
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	//"github.com/pborman/uuid"
	"github.com/robertsmoto/skustor/tools"
)

func TestBrand_AbsFunctions(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/brands.json")
	if err != nil {
		t.Errorf("TestBrands %s", err)
		log.Print("error", err)
	}
	brands := BrandNodes{}
	// open the db
	devPostgres := tools.PostgresDev{}
	devDb, err := tools.Open(&devPostgres)

	userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	date := "2022-06-01"
	userDir := "111111111111"

	// test the JsonLoadValidateUpsert interface
	err = JsonLoadValidateUpsert(&brands, testFile, devDb, userId)
	if err != nil {
		t.Error("Failed at JsonLoadValidateUpsert", err)
	}
	// test the RelatedTableUpsert
	err = RelatedTableUpsert(&brands, devDb, userId)
	if err != nil {
		t.Error("Failed at RelatedTableUpsert", err)
	}
	// test the ImageSizeUpsert interface
	err = ImageSizeUpsert(&brands, devDb, date, userId, userDir)
	if err != nil {
		t.Error("Failed at ImageSizeUpsert", err)
	}
}

func TestBrand_LoadValidate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/brands.json")
	if err != nil {
		t.Errorf("TestBrands %s", err)
		log.Print("error", err)
	}
	brands := BrandNodes{}
	// test the loader and validator
	err = brands.Load(&testFile)
	err = brands.Validate()
	if err != nil {
		t.Error("Brand Loader error", err)
	}
	node0 := brands.Nodes[0]
	if node0.Id == "eeb75266-7f4a-4d8e-9a8a-2c0ada73e7b1" != true {
		t.Error("Error loading Brands")
	}
}
