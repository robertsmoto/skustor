package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Account struct {
	Id        string `json:"id" validate:"required,uuid4"`
	ParentId        string `json:"parentId" validate:"omitempty,uuid4"`
	Auth      string `json:"auth" validate:"required,uuid4"`
	Key       string `json:"key" validate:"required,uuid4"`
	Username  string `json:"username" validate:"omitempty,gte=8,lte=100"`
	Firstname string `json:"firstname" validate:"omitempty,lte=100"`
	Lastname  string `json:"lastname" validate:"omitempty,lte=100"`
	Nickname  string `json:"nickname" validate:"omitempty,lte=100"`
	Document  string
}

type AccountNodes struct {
	Nodes []*Account `json:"accountNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *AccountNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "accountNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("AccountNodes.Load() %s", err)
	}
	return nil
}

func (s *AccountNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("AccountNodes.Validate() %s", err)
	}
	return nil
}

func (s *AccountNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO account (
                id, auth, key, document
            )
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET auth = $2,
                key = $3,
                document = $4
            WHERE account.id = $1;`
		_, err = db.Exec(
			qstr, accountId, node.Auth, node.Key, node.Document,
		)
		if err != nil {
			return fmt.Errorf("AccountNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *AccountNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		if node.ParentId == "" {
			continue
		}
		qstr := `
            UPDATE account
            SET parent_id = $2
            WHERE account.id = $1;`

		_, err = db.Exec(
			qstr, node.Id, node.ParentId,
		)
		if err != nil {
			return fmt.Errorf("AccountNodes.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *AccountNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
	fmt.Println("AccountNodes.RelatedTableUpsert Not implemented.")
	return nil
}

func (s *AccountNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("AccountNodes.Delete() Not implemented.")
	return nil
}
