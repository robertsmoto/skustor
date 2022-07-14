package models

import (
    //"fmt"
    "os"
    "testing"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_ContentInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_ContentInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/content.json")
    if err != nil {
        t.Errorf("Test_ContentInterfaces %s", err)
    }

    // open the db connections
    pgDb, err := postgres.Open(&postgres.PostgresDb{})

    // instantiate the structs
    //collection := Collection{}
    cNodes := ContentNodes{}

    // Little Johnnie account
    accountId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    err = LoadValidateHandler(&cNodes, &testFile)
    if err != nil {
        t.Errorf("Test_ContentInterfaces 01 %s", err)
    }
    err = UpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ContentInterfaces 02 %s", err)
    }
    err = ForeignKeyUpdateHandler(&cNodes, pgDb)
    if err != nil {
        t.Errorf("Test_ContentInterfaces 03 %s", err)
    }
    err = RelatedTableUpsertHandler(&cNodes, accountId, pgDb)
    if err != nil {
        t.Errorf("Test_ContentInterfaces 04 %s", err)
    }

    pgDb.Close()

}

func Test_ContentLoadAndValidate(t *testing.T) {
    testFile, err := os.ReadFile("./test_data/content.json")
    if err != nil {
        t.Errorf("Test_ContentLoadAndValidate %s", err)
    }
    contentNodes := ContentNodes{}
    contentNodes.Load(&testFile)
    contentNodes.Validate()
    for i, node := range contentNodes.Nodes {
        testId := "2f877877-7669-42b6-abed-6ebc20ba4c5b"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
