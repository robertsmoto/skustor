package models

import (
    //"fmt"
	"os"
	"testing"


    "github.com/robertsmoto/skustor/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_PersonInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    err = configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_PerrsonInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/people.json")
    if err != nil {
        t.Errorf("Test_PersonInterfaces %s", err)
    }

    // open the db connections
    postgres := postgres.PostgresDb{}
    pgDb, err := postgres.Open(&postgres)

    // instantiate the structs
    personNodes := PersonNodes{}

    // Little Johnnie user
    userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    procStructs := []LoaderProcesserUpserter{&personNodes}
    for _, s := range procStructs {
        err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
        if err != nil {
        t.Errorf("Test_PersonInterfaces %s", err)
        }
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
