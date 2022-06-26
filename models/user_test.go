package models

import (
    "os"
    "testing"
)

func Test_UserLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/users.json")
    if err != nil {
        t.Errorf("Test_UserLoadAndValidate %s", err)
    }

    // instantiate the user structs
    user := User{}
    useres := Users{}
    // test the loader and validator
    user.Load(&testFile)
    user.Validate()
    useres.Load(&testFile)
    useres.Validate()
    if user.Id != "5cc4a649-8138-4223-a62c-263cbb8a12b3" {
        t.Error("User.Id 01 ", err)
    }
    for i, user := range useres.Nodes {
        if i == 0 && user.Id != "f8b0f997-1dcc-4e56-915c-9f62f52345ee" {
            t.Error("User.Id 02 ", err)
        }
    }
}
