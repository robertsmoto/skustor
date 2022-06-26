package models

import (
	"os"
	"testing"
)

func Test_PersonLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/people.json")
    if err != nil {
        t.Errorf("TestPerson_LoadAndValidate %s", err)
    }
    // instantiate the structs
    person := Person{}
    people := People{}
    // test the loader and validator
    person.Load(&testFile)
    person.Validate()
    people.Load(&testFile)
    people.Validate()
    if person.Id != "0340a432-f7a9-428e-a70b-6a7a4f8bcdbc" {
        t.Error("Person.Id 01 ", err)
    }
    for i, person := range people.Nodes {
        if i == 0 && person.Id != "18788cb9-1abd-4822-9efa-f28d4443e042" {
            t.Error("Person.Id 02 ", err)
        }
    }
}
