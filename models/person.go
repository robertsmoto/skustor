package models

import (
	"database/sql"
    "encoding/json"
	"fmt"
	"errors"

    "github.com/go-playground/validator/v10"
)

type SvPerson struct {
	Id         string `json:"id" validate:"required,uuid4"`
	Type        string   `json:"type" validate:"required,lte=20,oneof=customer contact"`
	SvUserId         string `json:"svUserId" validate:"omitempty,uuid4"`
	PlaceId         string `json:"placeId" validate:"omitempty,uuid4"`
	Salutation string `json:"salutation" validate:"omitempty,lte=20"`
	Firstname  string `json:"firstname" validate:"omitempty,lte=100"`
	Lastname   string `json:"lastname" validate:"omitempty,lte=100"`
	Nickname   string `json:"nickname" validate:"omitempty,lte=100"`
	Phone      string `json:"phone" validate:"omitempty,lte=50"`
	Mobile     string `json:"mobile" validate:"omitempty,lte=50"`
	Email      string `json:"email" validate:"omitempty,email"`
    Address
	Addresses
    AddressIds []string
}

func (s *SvPerson) Process(userId string) (err error) {
    if userId == "" {
        return errors.New("SvPerson.Process() requires userId.")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvPerson) Upsert(db *sql.DB) (err error) {
	if s == (&SvPerson{}) {
		return nil
	}
	qstr := `
        INSERT INTO person (
            id, sv_user_id, salutation, firstname, lastname, nickname,
            phone, mobile, email
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id = $2,
            salutation = $3,
            firstname = $4,
            lastname = $5,
            nickname = $6,
            phone = $7
            mobile = $8,
            email = $9
        WHERE sv_user.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId),
		s.Salutation, s.Firstname, s.Lastname, s.Nickname, s.Phone,
		s.Mobile, s.Email)
	if err != nil {
		return fmt.Errorf("SvPerson.Upsert() %s", err)
	}
    return nil
}

func (s *SvPerson) ForeignKeyUpdate(db *sql.DB) (err error) {
    // addresses
    if s == (&SvPerson{}) {
        return nil
    }
    qstr := `
        UPDATE person
        SET place_id = $2
        WHERE person.id = $1;`

    _, err = db.Exec(
        qstr, FormatUUID(s.Id), FormatUUID(s.PlaceId),
    )
    if err != nil {
        return fmt.Errorf("SvPerson.ForeignKeyUpdate() %s", err)
    }
    return nil
}

func (s *SvPerson) RelatedTableUpsert(db *sql.DB) (err error) {
    // addressIds
    fmt.Println("SvPerson.RelatedTableUpsert() Not implemented.")
    if err != nil {
        return fmt.Errorf("SvPerson.RelatedTableUpsert() 01 %s", err)
    }
    return nil
}

func (s *SvPerson) Delete(db *sql.DB) (err error) {
    fmt.Println("SvPerson.Delete() Not implemented.")
    if err != nil {
        return fmt.Errorf("SvPerson.Delete() %s", err)
    }
    return nil
}

type Person struct {
    SvPerson `json:"person"`
}

func (s *Person) Load(fileBuffer *[]byte) (err error) {
    if s == (&Person{}) {
        return nil
    }
    err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Person.Load() %s", err)
	}
    return nil
}

func (s *Person) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Person.Validate() %s", err)
	}
    return nil
}

type People struct {
    Nodes []SvPerson `json:"people" validate:"dive"`
}

func (s *People) Load(fileBuffer *[]byte) (err error) {
    err = json.Unmarshal(*fileBuffer, &s)
    if err != nil {
        return fmt.Errorf("People.Load() %s", err)
    }
    return nil
}

func (s *People) Validate() (err error) {
    validate := validator.New()
    err = validate.Struct(s)
    if err != nil {
        return fmt.Errorf("People.Validate() %s", err)
    }
    return nil
}
