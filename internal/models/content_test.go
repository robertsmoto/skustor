package models

import (
	//"fmt"
    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/internal/postgres"
	"os"
	"testing"
)

func Test_ContentInterfaces(t *testing.T) {
    var err error

    // loading env variables (will eventually be loaded by main)
    conf := configs.Config{}
    err = configs.Load(&conf)
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
    contentNodes := ContentNodes{}

    // Little Johnnie user
    userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

    procStructs := []LoaderProcesserUpserter{&contentNodes}
    for _, s := range procStructs {
        err = JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
        if err != nil {
        t.Errorf("Test_CollectionInterfaces %s", err)
        }
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
