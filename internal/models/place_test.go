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
    configs.Load(&conf)
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
    pNodes := PlaceNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&pNodes, &testFile)
    if err != nil {
        t.Errorf("Test_PlaceInterfaces 01 %s", err)
    }
    err = UpsertHandler(&pNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_PlaceInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&pNodes, pgDb)
    if err != nil {
        t.Errorf("Test_PlaceInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&pNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_PlaceInterfaces 04 %s", err)
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
