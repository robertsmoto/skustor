package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type Address struct {
	Id         string `json:"id" validate:"omitempty,uuid4"`
	UserId     string `json:"userId" validate:"omitempty,uuid4"`
	LocationId string `json:"locationId" validate:"omitempty,uuid4"`
	Type       string `json:"type" validate:"omitempty,lte=100,oneof=billing main mailing shipping"`
	Street1    string `json:"street1" validate:"omitempty,lte=100"`
	Street2    string `json:"street2" validate:"omitempty,lte=100"`
	City       string `json:"city" validate:"omitempty,lte=100"`
	State      string `json:"state" validate:"omitempty,lte=50"`
	ZipCode    string `json:"zipCode" validate:"omitempty,lte=20"`
	Country    string `json:"country" validate:"omitempty,lte=50"`
}

type AddressNodes struct {
	Nodes []Address `json:"addressNodes" validate:"dive"`
}

func (s *AddressNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *AddressNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("GroupNodes.Validate() ", err)
	}
	return err
}

type Location struct {
	Id          string `json:"id" validate:"omitempty,uuid4"`
	UserId      string `json:"id" validate:"omitempty,uuid4"`
	Name        string `json:"name" validate:"omitempty,lte=100"`
	Type        string `json:"type" valdidate:"omitempty,lte=100,oneof=company store warehouse website"`
	Description string `json:"description" validate:"omitempty,lte=200"`
	Phone       string `json:"phone" validate:"omitempty,lte=20"`
	Email       string `json:"email" validate:"omitempty,lte=100,email"`
	Website     string `json:"website" validate:"omitempty,lte=100,url"`
	Domain      string `json:"domain" validate:"omitempty,lte=100"`
	AddressNodes
}

type LocationNodes struct {
	Nodes []Location `json:"locationNodes" validate:"dive"`
}

func (s *LocationNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *LocationNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("GroupNodes.Validate() ", err)
	}
	return err
}

func (s *Location) Upsert(db *sql.DB) (err error) {
	log.Print("Loation.Upsert() Not implemented.")

	//// check if struct is empty
	//if s == nil {
	//return err
	//}
	//// construct the sql upsert statement
	//qstr := `
	//INSERT INTO groups (
	//id, position, type, name, description, keywords,
	//link_url, link_text
	//)
	//VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	//ON CONFLICT (id) DO UPDATE
	//SET position=$2,
	//type = $3,
	//name = $4,
	//description = $5,
	//keywords = $6,
	//link_url = $7,
	//link_text = $8
	//WHERE groups.id = $1;`

	//_, err = db.Exec(
	//qstr, FormatUUID(s.Id), s.Position, s.Type, s.Name, s.Description,
	//s.Keywords, s.LinkUrl, s.LinkText,
	//)
	return err
}

func (s *Location) ForeignKeyUpdate(db *sql.DB) (err error) {

	log.Print("Loation.ForeignKeyUpadte() Not implemented.")
	//qstr := `
	//UPDATE groups
	//SET user_id = $2, parent_id = $3
	//WHERE id = $1;`

	//_, err = db.Exec(qstr, FormatUUID(s.Id), FormatUUID(s.UserId),
	//FormatUUID(s.ParentId),
	//)
	return err
}

func (s *Location) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("Loation.RelatedTableUpsert() Not implemented.")

	//if s.ItemIds != nil {
	//for _, id := range s.ItemIds {
	//err = JoinGroupItemUpsert(
	//db,
	//s.Position,
	//s.UserId, // user
	//s.Id,     // group id
	//id,       // item id
	//)
	//if err != nil {
	//log.Print(err)
	//}
	//}
	//}
	return err
}

func (s *LocationNodes) Delete(db sql.DB) (err error) {
	log.Print("Location.Delete() Not implemented.")
	return err
}
