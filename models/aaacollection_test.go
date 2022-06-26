package models

import (
	"os"
	"testing"
    "fmt"
    "github.com/robertsmoto/skustor/tools"
    "github.com/robertsmoto/skustor/configs"
)

//func Test_AaaCollectionInterfaces(t *testing.T) {
    //var err error

    //// loading env variables (will eventually be loaded by main)
    //conf := configs.Config{}
    //err = configs.Load(&conf)
    //if err != nil {
        //t.Errorf("Test_CollectionInterfaces %s", err)
    //}

    //// read file (will eventually come from the request)
    //testFile, err := os.ReadFile("./test_data/collections.json")
    //if err != nil {
        //t.Errorf("Test_CollectionInterfaces %s", err)
    //}

    //// open the db connections
    //postgres := tools.PostgresDb{}
    //pgDb, err := tools.Open(&postgres)

    //// instantiate the structs
    //collection := Collection{}
    //collections := Collections{}

    //// Little Johnnie user
    //userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    //err = JsonLoaderUpserterHandler(&collection, userId, &testFile, pgDb)
    //if err != nil {
        //t.Errorf("Test_CollectionInterfaces %s", err)
    //}

    //err = JsonLoaderUpserterHandler(&collections, userId, &testFile, pgDb)
    //if err != nil {
        //t.Errorf("Test_CollectionInterfaces %s", err)
    //}

    //pgDb.Close()

//}

func Test_AaaCollectionLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/aaaCollections.json")
    if err != nil {
        t.Errorf("TestCollection_LoadAndValidate %s", err)
    }
    // instantiate the structs
    collection := AaaCollection{}
    collections := AaaCollections{}
    // test the loader and validator
    collection.Load(&testFile)
    collection.Validate()
    collections.Load(&testFile)
    collections.Validate()


    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    err = configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_CollectionInterfaces %s", err)
    }


    // open the db connections
    postgres := tools.PostgresDb{}
    pgDb, err := tools.Open(&postgres)


    collections.Upsert(pgDb)
    fmt.Println("##collection --> ", collection)
    //fmt.Println("##collections --> ", collections)
    if collection.Id != "0f93a63a-13db-40e8-aa65-eecd37a86e8e" {
        t.Error("Collection.Id 01 ", err)
    }
    //for i, collection := range collections.Nodes {
        //if i == 0 && collection.Id != "eeb75266-7f4a-4d8e-9a8a-2c0ada73e7b1" {
            //t.Error("Collection.Id 02 ", err)
        //}
    //}
}
