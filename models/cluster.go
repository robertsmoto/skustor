package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type Cluster struct {
	ImageNodes
	SvUserId string
	Id       string `json:"id" validate:"required,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	Position uint8  `json:"position" validate:"omitempty,number"`
	// Type in the database eg. rawMaterialTag, partTag, productTag, productBrand
	// will need join table m2m relationship
	Type        string   `json:"type" validate:"required,lte=20,oneof=attribute variation brand department rawMaterialCategory partCategory productCategory postCategory pageCategory docCategory rawMaterialTag partTag productTag postTag pageTag docTag"`
	Name        string   `json:"name" validate:"omitempty,lte=200"`
	Description string   `json:"description" validate:"omitempty,lte=200"`
	Keywords    string   `json:"keywords" validate:"omitempty,lte=200"`
	LinkUrl     string   `json:"linkUrl" validate:"omitempty,url,lte=200"`
	LinkText    string   `json:"linkText" validate:"omitempty,lte=200"`
	ItemIds     []string `json:"itemIds" validate:"dive,omitempty,uuid4"`
}

type ClusterNodes struct {
	Nodes []Cluster `json:"clusterNodes" validate:"dive"`
}

func (s *ClusterNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *ClusterNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("ClusterNodes.Validate() ", err)
	}
	return err
}

func (s *Cluster) Upsert(db *sql.DB) (err error) {
	// check if struct is empty
	if s == (&Cluster{}) {
		log.Print("Cluster.Upsert() ", err)
		return err
	}

	qstr := `
        INSERT INTO cluster (
            id, sv_user_id, type, name, description, keywords,
            link_url, link_text
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id = $2,
            type = $3,
            name = $4,
            description = $5,
            keywords = $6,
            link_url = $7,
            link_text = $8
        WHERE cluster.id = $1;`

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId),
		s.Type, s.Name, s.Description, s.Keywords, s.LinkUrl,
		s.LinkText)

	if err != nil {
		log.Print("Cluster.Upsert() ", err)
	}

	return err
}

func (s *Cluster) ForeignKeyUpdate(db *sql.DB) (err error) {
	if s.ParentId == "" {
		return err
	}
	qstr := `
        UPDATE cluster
        SET parent_id = $2
        WHERE cluster.id = $1;`

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.ParentId),
	)
	if err != nil {
		log.Print("Cluster.ForeignKeyUpdate() ", err)
	}
	return err
}

func (s *Cluster) RelatedTableUpsert(db *sql.DB) (err error) {

	if s.ItemIds != nil {
		for _, id := range s.ItemIds {
			err = JoinClusterItemUpsert(
				db,
				s.SvUserId,
				s.Id,
				id,
				s.Position,
			)
			if err != nil {
				log.Print("Cluster.RelatedTableUpsert() ", err)
			}
		}
	}
	return err
}

func (s *Cluster) Delete(db sql.DB) (err error) {
	log.Print("Cluster.Delete() Not implemented.")
	return err
}
