package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type User struct {
	Id        string `json:"id" validate:"required,uuid4"`
	Auth      string `json:"auth" validate:"required,uuid4"`
	Key       string `json:"key" validate:"required,uuid4"`
	Username  string `json:"username" validate:"omitempty,gte=8,lte=100"`
	Firstname string `json:"firstname" validate:"omitempty,lte=100"`
	Lastname  string `json:"lastname" validate:"omitempty,lte=100"`
	Nickname  string `json:"nickname" validate:"omitempty,lte=100"`
	Document  string
}

type UserNodes struct {
	Nodes []*User `json:"userNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *UserNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "userNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("UserNodes.Load() %s", err)
	}
	return nil
}

func (s *UserNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("UserNodes.Validate() %s", err)
	}
	return nil
}

func (s *UserNodes) Upsert(userId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
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
		_, err = db.Exec(
			qstr, userId, node.Auth, node.Key,
			node.Username, node.Firstname, node.Lastname, node.Nickname,
		)
		if err != nil {
			return fmt.Errorf("UserNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *UserNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("UserNodes.ForeignKeyUpdate Not implemented.")
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

func (s *UserNodes) RelatedTableUpsert(userId string, db *sql.DB) (err error) {
	fmt.Println("UserNodes.ForeignKeyUpdate Not implemented.")
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

func (s *UserNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("UserNodes.Delete() Not implemented.")
	//for _, node := range s.Nodes {
	//if err != nil {
	//return fmt.Errorf("Collections.Delete() %s", err)
	//}
	//}
	return nil
}
