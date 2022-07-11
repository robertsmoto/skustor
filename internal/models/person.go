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
	PlaceId  string `json:"placeId" validate:"omitempty,uuid4"`
	//Id         string `json:"id" validate:"required,uuid4"`
	//Type       string `json:"type" validate:"required,lte=20,oneof=customer contact"`
	//SvUserId   string `json:"svUserId" validate:"omitempty,uuid4"`
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

func (s *PersonNodes) Upsert(userId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO person (
                id, sv_user_id, document
            )
            VALUES ($1, $2, $3)
            ON CONFLICT (id) DO UPDATE
            SET sv_user_id = $2,
                document = $3
            WHERE person.id = $1;`
		_, err = db.Exec(qstr, node.Id, userId, node.Document)
		if err != nil {
			return fmt.Errorf("PersonNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *PersonNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("PersonNodes.ForeignKeyUpdate Not implemented.")
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

func (s *PersonNodes) RelatedTableUpsert(userId string, db *sql.DB) (err error) {
	fmt.Println("PersonNodes.ForeignKeyUpdate Not implemented.")
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

func (s *PersonNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("PersonNodes.Delete() Not implemented.")
	//for _, node := range s.Nodes {
	//if err != nil {
	//return fmt.Errorf("Collections.Delete() %s", err)
	//}
	//}
	return nil
}
