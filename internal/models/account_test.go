package models

import (
    //"fmt"
    "os"
    "testing"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_AccountInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_AccountInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/accounts.json")
    if err != nil {
        t.Errorf("Test_AccountInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    aNodes := AccountNodes{}

    // Little Johnnie user
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&aNodes, &testFile)
    if err != nil {
        t.Errorf("Test_AccountInterfaces 01 %s", err)
    }
    err = UpsertHandler(&aNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_AccountInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&aNodes, pgDb)
    if err != nil {
        t.Errorf("Test_AccountInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&aNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_AccountInterfaces 04 %s", err)
    }
    pgDb.Close()
}
func Test_UserLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/accounts.json")
    if err != nil {
        t.Errorf("Test_AccountLoadAndValidate %s", err)
    }

    // instantiate the user structs
    accountNodes := AccountNodes{}
    // test the loader and validator
    accountNodes.Load(&testFile)
    accountNodes.Validate()
    for i, node := range accountNodes.Nodes {
        testId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
