package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Collection struct {
	Id       string   `json:"id" validate:"required,uuid4"`
	ParentId string   `json:"parentId" validate:"omitempty,uuid4"`
	ImageIds []string `json:"imageIds" validate:"dive,omitempty"`
	ItemIds  []string `json:"itemIds" validate:"dive,omitempty,uuid4"`
	SvUserId string
	Document string
}

type CollectionNodes struct {
	Nodes []Collection `json:"collectionNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *CollectionNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "collectionNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Collections.Load() %s", err)
	}
	return nil
}

func (s *CollectionNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Collections.Validate() %s", err)
	}
	return nil
}

func (s *CollectionNodes) Upsert(userId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO collection (
                id, sv_user_id, document
            )
            VALUES ($1, $2, $3)
            ON CONFLICT (id) DO UPDATE
            SET sv_user_id = $2,
                document = $3
            WHERE collection.id = $1;`

		_, err = db.Exec(
			qstr, FormatUUID(node.Id), FormatUUID(userId), node.Document,
		)
		if err != nil {
			return fmt.Errorf("Collections.Upsert() %s", err)
		}
	}
	return nil
}

func (s *CollectionNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		qstr := `
            UPDATE collection
            SET parent_id = $2
            WHERE collection.id = $1;`

		_, err = db.Exec(
			qstr, FormatUUID(node.Id), FormatUUID(node.ParentId),
		)
		if err != nil {
			return fmt.Errorf("Collections.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *CollectionNodes) RelatedTableUpsert(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		fmt.Println("node", node)
		//if s.ItemIds != nil {
		//for _, id := range s.ItemIds {
		//err = JoinCollectionItemUpsert(
		//db,
		//s.SvUserId,
		//s.Id,
		//id,
		//s.Position,
		//)
		//}
		//if err != nil {
		//return fmt.Errorf("Collection.RelatedTableUpsert() 01 %s", err)
		//}
		//}
		if err != nil {
			return fmt.Errorf("Collections.RelatedTableUpsert() %s", err)
		}
	}
	return nil
}

func (s *CollectionNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("CollectionNodes.Delete() Not implemented.")
	//for _, node := range s.Nodes {
	//if err != nil {
	//return fmt.Errorf("Collections.Delete() %s", err)
	//}
	//}
	return nil
}
