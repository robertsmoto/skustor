package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Person struct {
    BaseData
    AllIdNodes
	PlaceId  string `json:"placeId" validate:"omitempty,uuid4"`
	//Salutation string `json:"salutation" validate:"omitempty,lte=20"`
	//Firstname  string `json:"firstname" validate:"omitempty,lte=100"`
	//Lastname   string `json:"lastname" validate:"omitempty,lte=100"`
	//Nickname   string `json:"nickname" validate:"omitempty,lte=100"`
	//Phone      string `json:"phone" validate:"omitempty,lte=50"`
	//Mobile     string `json:"mobile" validate:"omitempty,lte=50"`
	//Email      string `json:"email" validate:"omitempty,email"`
	//Address
	//Addresses
	//AddressIds []string
}

type PersonNodes struct {
	Nodes []*Person `json:"personNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *PersonNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "personNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("PersonNodes.Load() %s", err)
	}
	return nil
}

func (s *PersonNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("PersonNodes.Validate() %s", err)
	}
	return nil
}

func (s *PersonNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO person (
                id, account_id, type, document
            )
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE person.id = $1;`
		_, err = db.Exec(qstr, node.Id, accountId, node.Type, node.Document)
		if err != nil {
			return fmt.Errorf("PersonNodes.Upsert %s", err)
		}
	}
	return nil
}

func (s *PersonNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		if node.ParentId == "" {
			continue
		}
		qstr := `
            UPDATE person
            SET parent_id = $2
            WHERE person.id = $1;`

		_, err = db.Exec(
			qstr, node.Id, node.ParentId,
		)
		if err != nil {
			return fmt.Errorf("PersonNodes.ForeignKeyUpdate %s", err)
		}
	}
	return nil
}

func (s *PersonNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        structArray := []Upserter{}
        ascendentColumn := "person_id"
        if node.CollectionIdNodes.Nodes != nil {
            node.CollectionIdNodes.ascendentColumn = ascendentColumn
            node.CollectionIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.CollectionIdNodes)
        }
        if node.ContentIdNodes.Nodes != nil {
            node.contentJson = s.Gjson.Array()[i].Get("contentIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        if node.ImageIdNodes.Nodes != nil {
            node.imageJson = s.Gjson.Array()[i].Get("imageIdNodes")
            node.ImageIdNodes.ascendentColumn = ascendentColumn
            node.ImageIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ImageIdNodes)
        }
        if node.ItemIdNodes.Nodes != nil {
            node.itemJson = s.Gjson.Array()[i].Get("itemIdNodes")
            node.ItemIdNodes.ascendentColumn = ascendentColumn
            node.ItemIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ItemIdNodes)
        }
        if node.PlaceIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("placeIdNodes")
            node.PlaceIdNodes.ascendentColumn = ascendentColumn
            node.PlaceIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PlaceIdNodes)
        }
        //if node.PersonIdNodes.Nodes != nil {
            //node.placeJson = s.Gjson.Array()[i].Get("personIdNodes")
            //node.PersonIdNodes.ascendentColumn = ascendentColumn
            //node.PersonIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.PersonIdNodes)
        //}
        for _, sa := range structArray {
            err = UpsertHandler(sa, accountId, db)
            if err != nil {
                return fmt.Errorf("PersonNodes.RelatedTableUpsert %s", err)
            }

        }
    }
	return nil
}

func (s *PersonNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("PersonNodes.Delete() Not implemented.")
	return nil
}
