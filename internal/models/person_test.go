package models

import (
    //"fmt"
    "os"
    "testing"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_PersonInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_PerrsonInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/people.json")
    if err != nil {
        t.Errorf("Test_PersonInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    pNodes := PersonNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&pNodes, &testFile)
    if err != nil {
        t.Errorf("Test_PersonInterfaces 01 %s", err)
    }
    err = UpsertHandler(&pNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_PersonInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&pNodes, pgDb)
    if err != nil {
        t.Errorf("Test_PersonInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&pNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_PersonInterfaces 04 %s", err)
    }
    pgDb.Close()
}

func Test_PersonLoadAndValidate(t *testing.T) {
    testFile, err := os.ReadFile("./test_data/people.json")
    if err != nil {
        t.Errorf("Test_PersonLoadAndValidate %s", err)
    }
    personNodes := PersonNodes{}
    personNodes.Load(&testFile)
    personNodes.Validate()
    for i, node := range personNodes.Nodes {
        testId := "18788cb9-1abd-4822-9efa-f28d4443e042"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
