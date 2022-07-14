package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Collection struct {
    BaseData
	AllIdNodes
}

type CollectionNodes struct {
	Nodes []*Collection `json:"collectionNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *CollectionNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "collectionNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("CollectionNodes.Load() %s", err)
	}
	return nil
}

func (s *CollectionNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("CollectionNodes.Validate() %s", err)
	}
	return nil
}

func (s *CollectionNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO collection (
                id, account_id, type, document
            )
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE collection.id = $1;`

		_, err = db.Exec(
			qstr, node.Id, accountId, node.Type, node.Document,
		)
		if err != nil {
			return fmt.Errorf("CollectionNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *CollectionNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		if node.ParentId == "" {
			continue
		}
		qstr := `
            UPDATE collection
            SET parent_id = $2
            WHERE collection.id = $1;`

		_, err = db.Exec(
			qstr, node.Id, node.ParentId,
		)
		if err != nil {
			return fmt.Errorf("CollectionNodes.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *CollectionNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        ascendentColumn := "collection_id"
        structArray := []Upserter{}
        //if node.CollectionIdNodes.Nodes != nil {
            //node.CollectionIdNodes.ascendentColumn = ascendentColumn
            //node.CollectionIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.CollectionIdNodes)
        //}
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
        if node.PersonIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("personIdNodes")
            node.PersonIdNodes.ascendentColumn = ascendentColumn
            node.PersonIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PersonIdNodes)
        }
        for _, sa := range structArray {
            err = UpsertHandler(sa, accountId, db)
            if err != nil {
                return fmt.Errorf("CollectionNodes.RelatedTableUpsert %s", err)
            }
        }
    }
	return nil
}

func (s *CollectionNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("CollectionNodes.Delete() Not implemented.")
	return nil
}
