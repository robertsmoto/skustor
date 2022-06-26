package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

func (s *SvUser) Process(userId string) (err error) {
    fmt.Println("SvUser.Process() Not iplemented.")
    if err != nil {
        return fmt.Errorf("SvUser.Process() %s", err)
    }
    return nil
}

func (s *SvUser) Upsert(db *sql.DB) (err error) {
	if *s == (SvUser{}) {
		return nil
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
		return fmt.Errorf("SvUser.Upsert() %s", err)
	}
    return nil
}

func (s *SvUser) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("SvUser.ForeignKeyUpsert Not implemented.")
	if err != nil {
		return fmt.Errorf("SvUser.ForeignKeyUpdate() %s", err)
	}
    return nil
}

func (s *SvUser) RelatedTableUpsert(db *sql.DB) (err error) {
	fmt.Println("SvUser.RelatedTableUpsert Not implemented.")
	if err != nil {
		return fmt.Errorf("SvUser.RelatedTableUpsert() %s", err)
	}
    return nil
}

func (s *SvUser) Delete(db *sql.DB) (err error) {
	fmt.Print("SvUser.Delete() Not implemented.")
    return nil
}

type User struct {
	SvUser `json:"user"`
}

func (s *User) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("User.Load() %s", err)
	}
    return err
}

func (s *User) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("User.Validate() %s", err)
	}
    return nil
}

type Users struct {
	Nodes []SvUser `json:"users" validate:"dive"`
}

func (s *Users) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Users.Load() %s", err)
	}
    return nil
}

func (s *Users) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Users.Validate() %s", err)
	}
    return nil
}

func (s *Users) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("Users.Process() %s", err)
        }
    }
    return nil
}

func (s *Users) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Users.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Users) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Users.ForeignKeyUpdate() %s", err)
        }
    }
    return nil
}

func (s *Users) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Users.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Users) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Users.Delete() %s", err)
        }
    }
    return nil
}
