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
    //Id       string `json:"id" validate:"omitempty,uuid4"`
    //SvUserId string
    //Document string

	//Id          string   `json:"id" validate:"omitempty,uuid4"`
	//SvUserId    string   `json:"svUserId" validate:"omitempty,uuid4"`
	//Type        string   `json:"type" valdidate:"omitempty,lte=100,oneof=company store warehouse website"`
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

func (s *PlaceNodes) Upsert(userId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO place (
                id, sv_user_id, document
            )
            VALUES ($1, $2, $3)
            ON CONFLICT (id) DO UPDATE
            SET sv_user_id = $2,
                document = $3
            WHERE place.id = $1;`
		_, err = db.Exec(qstr, node.Id, userId, node.Document)
		if err != nil {
			return fmt.Errorf("PlaceNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *PlaceNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("PlaceNodes.ForeignKeyUpdate Not implemented.")
	//for _, node := range s.Nodes {
	//qstr := `
	//UPDATE collection
	//SET parent_id = $2
	//WHERE collection.id = $1;`

	//_, err = db.Exec(
	//qstr, FormatUUID(node.Id), FormatUUID(node.ParentId),
	//)
	//if err != nil {
	//return fmt.Errorf("Collections.ForeignKeyUpdate() %s", err)
	//}
	//}
	return nil
}

func (s *PlaceNodes) RelatedTableUpsert(userId string, db *sql.DB) (err error) {
	fmt.Println("PlaceNodes.ForeignKeyUpdate Not implemented.")
	//for _, node := range s.Nodes {
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
	//if err != nil {
	//return fmt.Errorf("Collections.RelatedTableUpsert() %s", err)
	//}
	//}
	return nil
}

func (s *PlaceNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("PlaceNodes.Delete() Not implemented.")
	//for _, node := range s.Nodes {
	//if err != nil {
	//return fmt.Errorf("Collections.Delete() %s", err)
	//}
	//}
	return nil
}
