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
		t.Errorf("Test_ContentInterfaces %s", err)
	}

	// read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/items.json")
	if err != nil {
		t.Errorf("Test_ContentInterfaces %s", err)
	}

	// open the db connections
	pgDb, err := postgres.Open(&postgres.PostgresDb{})

	// instantiate the structs
	itemNodes := ItemNodes{}

	// Little Johnnie user
	userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

	procStructs := []LoaderProcesserUpserter{&itemNodes}
	for _, s := range procStructs {
		err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
		if err != nil {
			t.Errorf("Test_ItemInterfaces %s", err)
		}
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
