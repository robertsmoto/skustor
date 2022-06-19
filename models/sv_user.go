package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/pborman/uuid"
)

type SvUser struct {
	Id        string `json:"id" validate:"omitempty,uuid4"`
	Auth      string `json:"auth" validate:"omitempty,uuid4"`
	Key       string `json:"key" validate:"omitempty,uuid4"`
	Username  string `json:"username" validate:"omitempty,gte=8,lte=100"`
	Firstname string `json:"firstname" validate:"omitempty,lte=100"`
	Lastname  string `json:"lastname" validate:"omitempty,lte=100"`
	Nickname  string `json:"nickname" validate:"omitempty,lte=100"`
}

type SvUserNodes struct {
	Nodes []SvUser `json:"svUserNodes" validate:"dive"`
}

func (s *SvUserNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *SvUserNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("GroupNodes.Validate() ", err)
	}
	return err
}

func (s *SvUser) Upsert(db *sql.DB) (err error) {
	if *s == (SvUser{}) {
		return err
	}
	// construct the sql upsert statement
	qstr := `
        INSERT INTO sv_user (
            id, auth, key, username, firstname, lastname, nickname
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE
        SET auth = $2,
            key = $3,
            username = $4,
            firstname = $5,
            lastname = $6,
            nickname = $7
        WHERE sv_user.id = $1;`
	// execute it

	_, err = db.Exec(
		qstr, uuid.Parse(s.Id), uuid.Parse(s.Auth), uuid.Parse(s.Key),
		s.Username, s.Firstname, s.Lastname, s.Nickname,
	)
    if err != nil {
        log.Print("SvUser.Upsert() ", err)
    }
	return err
}

func (s *SvUser) ForeignKeyUpdate(db *sql.DB) (err error) {
	log.Print("SvUser.ForeignKeyUpsert Not implemented.")
	return err
}

func (s *SvUser) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("SvUser.RelatedTableUpsert Not implemented.")
	return err
}
