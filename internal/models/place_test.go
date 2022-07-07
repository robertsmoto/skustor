package models

import (
    //"fmt"
    "os"
    "testing"


    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_PlaceInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    err = configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_PlaceInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/places.json")
    if err != nil {
        t.Errorf("Test_PlaceInterfaces %s", err)
    }

    // open the db connections
	pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    placeNodes := PlaceNodes{}

    // Little Johnnie user
    userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    procStructs := []LoaderProcesserUpserter{&placeNodes}
    for _, s := range procStructs {
        err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
        if err != nil {
        t.Errorf("Test_PlaceInterfaces %s", err)
        }
    }
    pgDb.Close()
}

func Test_PlaceLoadAndValidate(t *testing.T) {
    testFile, err := os.ReadFile("./test_data/places.json")
    if err != nil {
        t.Errorf("Test_PlaceLoadAndValidate %s", err)
    }
    placeNodes := PlaceNodes{}
    placeNodes.Load(&testFile)
    placeNodes.Validate()
    for i, node := range placeNodes.Nodes {
        testId := "58d94686-4dd6-4c10-b255-fc40ebfd56e1"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
