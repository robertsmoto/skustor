package models

import (
	"os"
	"testing"

)

func Test_ItemLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/items.json")
    if err != nil {
        t.Errorf("Test_ItemLoadAndValidate %s", err)
    }
    // instantiate the structs
    item := Item{}
    items := Items{}
    // test the loader and validator
    item.Load(&testFile)
    item.Validate()
    items.Load(&testFile)
    items.Validate()
    if item.Id != "723a1838-a3b2-4fcd-9ec9-24f8b2c2dde8" {
        t.Error("Item.Id 01 ", err)
    }
    for i, item := range items.Nodes {
        if i == 0 && item.Id != "2940d429-9d97-4d9e-a80f-67d56e42d226" {
            t.Error("Item.Id 02 ", err)
        }
    }
}
