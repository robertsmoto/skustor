package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Content struct {
	BaseData
    AllIdNodes

	//type: article, page, docs
	//PublishedTime string   `json:"publishedTime" validate:"omitempty,datetime=15:04 MST"`
	//Published     string   `json:"published" validate:"omitempty,datetime=2006-01-02"`
	//Modified      string   `json:"modified" validate:"omitempty,datetime=2006-01-02"`
	//Keywords      string   `json:"keywords"`
	//Title         string   `json:"title"`
	//Excerpt       string   `json:"excerpt"`
	//Body          string   `json:"body" validate:"omitempty"`
	//Footer        string   `json:"footer"`
}

type ContentNodes struct {
	Nodes []*Content `json:"contentNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *ContentNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "contentNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("ContentNodes.Load() %s", err)
	}
	return nil
}

func (s *ContentNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("ContentNodes.Validate() %s", err)
	}
	return nil
}

func (s *ContentNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO content (
                id, account_id, type, document
            )
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE content.id = $1;`
		_, err = db.Exec(qstr, node.Id, accountId, node.Type, node.Document)
		if err != nil {
			return fmt.Errorf("ContentNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *ContentNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
        var qstr string
        if node.ParentId != "" {
            qstr = `
                UPDATE content
                SET parent_id = $2
                WHERE content.id = $1;`

            _, err = db.Exec(qstr, node.Id, node.ParentId)

		} else {
            // need this query to set fk uuid to null
            qstr = `
                UPDATE content
                SET parent_id = null
                WHERE content.id = $1;`
            _, err = db.Exec(qstr, node.Id)
        }

		if err != nil {
			return fmt.Errorf("ContentNodes.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *ContentNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        ascendentColumn := "content_id"
        structArray := []Upserter{}
        if node.CollectionIdNodes.Nodes != nil {
            node.collectionJson = s.Gjson.Array()[i].Get("collectionIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        //if node.ContentIdNodes.Nodes != nil {
            //node.contentJson = s.Gjson.Array()[i].Get("contentIdNodes")
            //node.ContentIdNodes.ascendentColumn = ascendentColumn
            //node.ContentIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.ContentIdNodes)
        //}
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
                return fmt.Errorf("ContentNodes.RelatedTableUpsert %s", err)
            }
        }
    }
	return nil
}

func (s *ContentNodes) Delete(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		qstr := `
            DELETE FROM content
            WHERE content.id = $1;
            `
		_, err = db.Exec(qstr, node.Id)
		if err != nil {
			return fmt.Errorf("ContentNodes.Delete %s", err)
		}
	}
	return nil
}
