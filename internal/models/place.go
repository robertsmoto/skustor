package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Place struct {
    BaseData
    AllIdNodes
	//Name        string   `json:"name" validate:"omitempty,lte=100"`
	//Description string   `json:"description" validate:"omitempty,lte=200"`
	//Phone       string   `json:"phone" validate:"omitempty,lte=20"`
	//Email       string   `json:"email" validate:"omitempty,lte=100,email"`
	//Website     string   `json:"website" validate:"omitempty,lte=100,url"`
	//Domain      string   `json:"domain" validate:"omitempty,lte=100"`
	//AddressIds  []string `json:"addressIds" validate:"dive"`
	//Address
	//Addresses
}

type PlaceNodes struct {
	Nodes []*Place `json:"placeNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *PlaceNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "placeNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("PlaceNodes.Load() %s", err)
	}
	return nil
}

func (s *PlaceNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("PlaceNodes.Validate() %s", err)
	}
	return nil
}

func (s *PlaceNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO place (
                id, account_id, type, document
            )
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE place.id = $1;`
		_, err = db.Exec(qstr, node.Id, accountId, node.Type, node.Document)
		if err != nil {
			return fmt.Errorf("PlaceNodes.Upsert %s", err)
		}
	}
	return nil
}

func (s *PlaceNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
        var qstr string
        if node.ParentId != "" {
            qstr = `
                UPDATE place
                SET parent_id = $2
                WHERE place.id = $1;`

            _, err = db.Exec(qstr, node.Id, node.ParentId)

		} else {
            // need this query to set fk uuid to null
            qstr = `
                UPDATE place
                SET parent_id = null
                WHERE place.id = $1;`
            _, err = db.Exec(qstr, node.Id)
        }

		if err != nil {
			return fmt.Errorf("PlaceNodes.ForeignKeyUpdate %s", err)
		}
	}
	return nil
}

func (s *PlaceNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        structArray := []Upserter{}
        ascendentColumn := "place_id"
        if node.CollectionIdNodes.Nodes != nil {
            node.collectionJson = s.Gjson.Array()[i].Get("collectionIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
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
        //if node.PlaceIdNodes.Nodes != nil {
            //node.placeJson = s.Gjson.Array()[i].Get("placeIdNodes")
            //node.PlaceIdNodes.ascendentColumn = ascendentColumn
            //node.PlaceIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.PlaceIdNodes)
        //}
        if node.PersonIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("personIdNodes")
            node.PersonIdNodes.ascendentColumn = ascendentColumn
            node.PersonIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PersonIdNodes)
        }
        for _, sa := range structArray {
            err = UpsertHandler(sa, accountId, db)
            if err != nil {
                return fmt.Errorf("PlaceNodes.RelatedTableUpsert %s", err)
            }

        }
    }
	return nil
}

func (s *PlaceNodes) Delete(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		qstr := `
            DELETE FROM place
            WHERE place.id = $1;
            `
		_, err = db.Exec(qstr, node.Id)
		if err != nil {
			return fmt.Errorf("PlaceNodes.Delete %s", err)
		}
	}
	return nil
}
