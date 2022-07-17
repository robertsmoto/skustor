package models

import (
	"os"
	"testing"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_CollectionInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/collections.json")
    if err != nil {
        t.Errorf("Test_CollectionInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    //collection := Collection{}
    cNodes := CollectionNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&cNodes, &testFile)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces 01 %s", err)
    }
    err = UpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&cNodes, pgDb)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces 04 %s", err)
    }
    pgDb.Close()
}

func Test_CollectionLoadAndValidate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/collections.json")
	if err != nil {
		t.Errorf("TestCollection_LoadAndValidate %s", err)
	}
	collections := CollectionNodes{}

	// test the loader and validator
	//collection.Load(&testFile)
	//collection.Validate()
	collections.Load(&testFile)
	collections.Validate()

	//if collection.Id != "0f93a63a-13db-40e8-aa65-eecd37a86e8e" {
	//t.Error("Collection.Id 01 ", err, collection.Id)
	//}
	for i, collection := range collections.Nodes {
		if i == 0 && collection.Id != "f9758d28-f580-4da5-bcbf-097d101f8270" {
			t.Error("Collection.Id 02 ", err)
		}
	}
}
