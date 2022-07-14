package models

import (
    //"fmt"
    "os"
    "testing"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_ItemInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_ItemInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/items.json")
    if err != nil {
        t.Errorf("Test_ItemInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    iNodes := ItemNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&iNodes, &testFile)
    if err != nil {
        t.Errorf("Test_ItemInterfaces 01 %s", err)
    }
    err = UpsertHandler(&iNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ItemInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&iNodes, pgDb)
    if err != nil {
        t.Errorf("Test_ItemInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&iNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ItemInterfaces 04 %s", err)
    }
    pgDb.Close()
}

func Test_ItemLoadAndValidate(t *testing.T) {
    testFile, err := os.ReadFile("./test_data/items.json")
    if err != nil {
        t.Errorf("Test_ItemLoadAndValidate %s", err)
    }
    itemNodes := ItemNodes{}
    itemNodes.Load(&testFile)
    itemNodes.Validate()
    for i, node := range itemNodes.Nodes {
        testId := "2940d429-9d97-4d9e-a80f-67d56e42d226"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
