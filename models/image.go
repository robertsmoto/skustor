package models

import (
	"database/sql"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	//"github.com/pborman/uuid"
)

type Image struct {
	Id       string `json:"id" validate:"omitempty,uuid4"`
	Url      string `json:"url" validate:"omitempty,url"`
	Title    string `json:"title" validate:"omitempty,lte=200"`
	Alt      string `json:"alt" validate:"omitempty,lte=100"`
	Caption  string `json:"caption" validate:"omitempty,lte=200"`
	Position uint8  `json:"position" validate:"omitempty,number"`
	Featured uint8  `json:"featured" validate:"omitempty,number"`
}

type ImageNodes struct {
	Nodes []Image `json:"imageNodes" validate:"dive"`
}

func (s *ImageNodes) Load(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

func (s *ImageNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	return err
}

func Upsert(db *sql.DB, userId string) (err error) {
	return err
}
