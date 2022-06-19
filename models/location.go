package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type Address struct {
	Id         string `json:"id" validate:"omitempty,uuid4"`
	SvUserId   string `json:"userId" validate:"omitempty,uuid4"`
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

func (s *Address) Upsert(db *sql.DB) (err error) {
	if s == (&Address{}) {
		return err
	}
	qstr := `
        INSERT INTO address (
            id, sv_user_id, type, street1, street2, city,
            state, zipcode, country
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id=$2,
            type = $3,
            street1 = $4,
            street2 = $5,
            city = $6,
            state = $7,
            zipcode = $8,
            country = $9
        WHERE address.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Type, s.Street1, s.Street2,
		s.City, s.State, s.ZipCode, s.Country,
	)
	if err != nil {
		log.Print("Address.Upsert() ", err)
	}
	return err
}

func (s *Address) ForeignKeyUpdate(db *sql.DB) (err error) {
	if s.LocationId == "" {
		return err
	}
	qstr := `
        UPDATE address
        SET location_id = $2
        WHERE address.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.LocationId),
	)
	if err != nil {
		log.Print("Address.ForeignKeyUpdate() ", err)
	}
	return err
}

func (s *Address) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("Address.RelatedTableUpsert() Not implemented.")
	return err
}

func (s *Address) Delete(db sql.DB) (err error) {
	log.Print("Address.Delete() Not implemented.")
	return err
}

type Location struct {
	Id          string   `json:"id" validate:"omitempty,uuid4"`
	SvUserId    string   `json:"svUserId" validate:"omitempty,uuid4"`
	Type        string   `json:"type" valdidate:"omitempty,lte=100,oneof=company store warehouse website"`
	Name        string   `json:"name" validate:"omitempty,lte=100"`
	Description string   `json:"description" validate:"omitempty,lte=200"`
	Phone       string   `json:"phone" validate:"omitempty,lte=20"`
	Email       string   `json:"email" validate:"omitempty,lte=100,email"`
	Website     string   `json:"website" validate:"omitempty,lte=100,url"`
	Domain      string   `json:"domain" validate:"omitempty,lte=100"`
	AddressIds  []string `json:"addressIds" validate:"dive"`
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
	if s == (&Location{}) {
		return err
	}
	qstr := `
        INSERT INTO location (
            id, sv_user_id, type, name, description, phone,
            email, website, domain
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id=$2,
            type = $3,
            name = $4,
            description = $5,
            phone = $6,
            email = $7,
            website = $8,
            domain = $9
        WHERE location.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Type, s.Name, s.Description,
		s.Phone, s.Email, s.Website, s.Domain,
	)
	if err != nil {
		log.Print("Location.Upsert() ", err)
	}
	return err
}

func (s *Location) ForeignKeyUpdate(db *sql.DB) (err error) {
	log.Print("Loation.ForeignKeyUpdate() Not implemented.")
	return err
}

func (s *Location) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("Loation.RelatedTableUpsert() Not implemented.")
	return err
}

func (s *LocationNodes) Delete(db sql.DB) (err error) {
	log.Print("Location.Delete() Not implemented.")
	return err
}
