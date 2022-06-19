package models

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	//"github.com/pborman/uuid"
	//"github.com/robertsmoto/skustor/configs"
	"github.com/robertsmoto/skustor/tools"
)

func TestUser_Upsert(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/users.json")
	if err != nil {
		t.Error("TestUser_Upsert() ", err)
	}
	users := UserNodes{}
	// open the db
	devPostgres := tools.PostgresDev{}
	devDb, err := tools.Open(&devPostgres)

	// test the JsonLoadValidateUpsert interface
	// this loads the []structs
	LoaderHandler(&users, testFile)

	// now loop through eash struct indvidually
	//userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	//date := "2022-06-01"
	//userDir := "111111111111"

	for _, user := range users.Nodes {
		UpsertHandler(&user, devDb)
	}

	// test the Image interface
}

func TestUser_Validate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/users.json")
	if err != nil {
		t.Errorf("TestUsers %s", err)
	}
	users := UserNodes{}
	// test the loader and validator
	err = users.Load(&testFile)
	err = users.Validate()
	if err != nil {
		t.Error("User Loader error", err)
	}
	node0 := users.Nodes[0]
	if node0.Id == "f8b0f997-1dcc-4e56-915c-9f62f52345ee" != true {
		t.Error("Error loading User")
	}
	if err != nil {
		t.Error("TestUser_Validate Error: ", err)
	}
}
