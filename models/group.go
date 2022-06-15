package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	//"sync"
	//"path/filepath"

	//_ "github.com/lib/pq"
	"github.com/go-playground/validator/v10"
	"github.com/pborman/uuid"
	//"github.com/robertsmoto/skustor/configs"
	//"github.com/robertsmoto/skustor/images"
)

type Group struct {
	ImageNodes
	userId   string
	Id       string `json:"id" validate:"required,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	// Type in the database eg. rawMaterialTag, partTag, productTag, productBrand
	// will need join table m2m relationship
	Type           string   `json:"type" validate:"required,lte=20,oneof=brand department rawMaterialCategory partCategory productCategory postCategory pageCategory docCategory rawMaterialTag partTag productTag postTag pageTag docTag"`
	Name           string   `json:"name" validate:"omitempty,lte=200"`
	Description    string   `json:"description" validate:"omitempty,lte=200"`
	Keywords       string   `json:"keywords" validate:"omitempty,lte=200"`
	LinkUrl        string   `json:"linkUrl" validate:"omitempty,url,lte=200"`
	LinkText       string   `json:"linkText" validate:"omitempty,lte=200"`
	RawMaterialIds []string `json:"rawMaterialIds" validate:"dive,omitempty,uuid4"`
	PartIds        []string `json:"partIds" validate:"dive,omitempty,uuid4"`
	ProductIds     []string `json:"productIds" validate:"dive,omitempty,uuid4"`
}

type GroupNodes struct {
	Nodes []Group `json:"groupNodes" validate:"dive"`
}

func (s *GroupNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *GroupNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("GroupNodes.Validate() ", err)
	}
	return err
}

func (s *Group) Upsert(db *sql.DB) (err error) {
	// check if struct is empty
	if s.Nodes == nil {
		log.Print("GroupNodes.Upsert() ", err)
		return err
	}
	// construct the sql upsert statement
	qstr := `
        INSERT INTO groups (
            id, user_id, parent_id, type, name, description, keywords,
           link_url, link_text
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET user_id = $2,
            parent_id = $3,
            type = $4,
            name = $5,
            description = $6,
            keywords = $7,
            link_url = $8,
            link_text = $9
        WHERE groups.id = $1;`
	// execute it

	fmt.Println("node.Id ", s.Id)
	fmt.Println("userId ", s.userId)
	fmt.Println("parentId ", s.ParentId)
	var pid uuid.UUID
	if s.ParentId == "" {
		pid = uuid.Parse("00000000-0000-0000-0000-000000000000")
	} else {
		pid = uuid.Parse(s.ParentId)
	}
	_, err = db.Exec(
		qstr, uuid.Parse(s.Id), uuid.Parse(s.userId), pid,
		s.Type, s.Name, s.Description, s.Keywords, s.LinkUrl,
		s.LinkText)

	return err
}

func (s *Group) RelatedTableUpsert(db *sql.DB) (err error) {

	if s.RawMaterialIds != nil {
		for _, id := range s.RawMaterialIds {
			err = JoinGroupItemUpsert(
				db,
				s.userId,
				s.Id,
				id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if s.PartIds != nil {
		for _, id := range s.PartIds {
			err = JoinGroupItemUpsert(
				db,
				s.userId,
				s.Id,
				id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if s.ProductIds != nil {
		for _, id := range s.ProductIds {
			err = JoinGroupItemUpsert(
				db,
				s.userId,
				s.Id,
				id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return err
}

func (s *GroupNodes) Delete(db sql.DB) (err error) {
	log.Print("brand.Delete() Not implemented.")
	return err
}
