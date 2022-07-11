package models

import (
    "fmt"
    "log"
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
	collectionNodes := CollectionNodes{}

	// Little Johnnie user
	userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

	structs := []LoaderProcesserUpserter{&collectionNodes}
	for _, s := range structs {
		err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
		if err != nil {
			t.Errorf("Test_CollectionInterfaces %s", err)
		}
	}

    qstr := `
        SELECT document FROM collection
        WHERE id = 'eeb75266-7f4a-4d8e-9a8a-2c0ada73e7b1';
        `
    var x string
    result := pgDb.QueryRow(qstr).Scan(&x)

    fmt.Println("result ", x)
    fmt.Printf("result %T ", result)

    if err != nil {
        log.Print("Error creating or updating database.", err)
    } else {
        log.Print("Successfully updated postgresDb.")
    }

	//err = JsonLoaderUpserterHandler(&collection, userId, &testFile, pgDb)
	//if err != nil {
	//t.Errorf("Test_CollectionInterfaces %s", err)
	//}

	//err = JsonLoaderUpserterHandler(&collections, userId, &testFile, pgDb)
	//if err != nil {
	//t.Errorf("Test_CollectionInterfaces %s", err)
	//}

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

	//fmt.Println("### collection", collection)
	//fmt.Println("##id ", collection.Id)
	//if collection.Id != "0f93a63a-13db-40e8-aa65-eecd37a86e8e" {
	//t.Error("Collection.Id 01 ", err, collection.Id)
	//}
	for i, collection := range collections.Nodes {
		if i == 0 && collection.Id != "eeb75266-7f4a-4d8e-9a8a-2c0ada73e7b1" {
			t.Error("Collection.Id 02 ", err)
		}
	}
}
