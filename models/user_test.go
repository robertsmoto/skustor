package models

import (
    //"fmt"
    "os"
    "testing"


    "github.com/robertsmoto/skustor/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
)

func Test_UserInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    err = configs.Load(&conf)
    if err != nil {
        t.Errorf("Test_UserInterfaces %s", err)
    }

    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/users.json")
    if err != nil {
        t.Errorf("Test_UserInterfaces %s", err)
    }

    // open the db connections
    postgres := postgres.PostgresDb{}
    pgDb, err := postgres.Open(&postgres)

    // instantiate the structs
    userNodes := UserNodes{}

    // Little Johnnie user
    userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    procStructs := []LoaderProcesserUpserter{&userNodes}
    for _, s := range procStructs {
        err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
        if err != nil {
        t.Errorf("Test_UserInterfaces %s", err)
        }
    }
    pgDb.Close()
}
func Test_UserLoadAndValidate(t *testing.T) {
	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/users.json")
	if err != nil {
		t.Errorf("Test_UserLoadAndValidate %s", err)
	}

	// instantiate the user structs
	userNodes := UserNodes{}
	// test the loader and validator
	userNodes.Load(&testFile)
	userNodes.Validate()
    for i, node := range userNodes.Nodes {
        testId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
        if i == 0 && node.Id != testId {
            t.Errorf("node.Id 02 %s != %s ", node.Id, testId)
        }
    }
}
