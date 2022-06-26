package models

import (
    "os"
    "testing"
)

func Test_PlaceLoadAndValidate(t *testing.T) {
    // read file (will eventually come from the request)
    testFile, err := os.ReadFile("./test_data/places.json")
    if err != nil {
        t.Errorf("Test_PlaceLoadAndValidate %s", err)
    }

    // instantiate the address structs
    address := Address{}
    addresses := Addresses{}
    // test the loader and validator
    address.Load(&testFile)
    address.Validate()
    addresses.Load(&testFile)
    addresses.Validate()
    if address.Id != "b5d72565-a305-44c3-b979-24bba466b637" {
        t.Error("Address.Id 01 ", err)
    }
    for i, address := range addresses.Nodes {
        if i == 0 && address.Id != "a3f982c9-6537-440c-8191-39445044f2f9" {
            t.Error("Address.Id 02 ", err)
        }
    }

    // instantiate the structs
    place := Place{}
    places := Places{}
    // test the loader and validator
    place.Load(&testFile)
    place.Validate()
    places.Load(&testFile)
    places.Validate()
    if place.Id != "becf155e-77ee-419d-bfae-d63c8cd687b1" {
        t.Error("Place.Id 01 ", err)
    }
    for i, place := range places.Nodes {
        if i == 0 && place.Id != "58d94686-4dd6-4c10-b255-fc40ebfd56e1" {
            t.Error("Place.Id 02 ", err)
        }
    }
}
