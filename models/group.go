package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	//_ "github.com/lib/pq"
	"github.com/go-playground/validator/v10"
	"github.com/pborman/uuid"
)

type Loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type Validater interface {
	Validate() (err error)
}

type Upserter interface {
	Upsert(db *sql.DB, userId uuid.UUID) (err error)
}

type LoaderValidaterUpserter interface {
	Loader
	Validater
	Upserter
}

func JsonHandler(data LoaderValidaterUpserter, fileBuffer []byte, db *sql.DB, userId uuid.UUID) (err error){
    err = data.Load(&fileBuffer)
    err = data.Validate()
    err = data.Upsert(db, userId)
    return err
}

func JoinTableUpsert(
	db *sql.DB,
	idArray []uuid.UUID,
	table, col1, col2 string,
	staticId uuid.UUID,
) (err error) {

	for _, varId := range idArray {
		qstr := fmt.Sprintf(`
            INSERT INTO %s (%s, %s)
            VALUES ($1, $2)
            ON CONFLICT (%s, %s)
            DO UPDATE SET %s = $1, %s = $2;`,
			table, col1, col2, col1, col2, col1, col2)

		_, err = db.Exec(qstr, staticId, varId)
	}
	return err
}

type Group struct {
	Id       uuid.UUID `json:"id" validate:"required"`
	ParentId uuid.UUID `json:"parentId" validate:"omitempty"`
	// Type in the database eg. rawMaterialTag, partTag, productTag, productBrand
	// will need join table m2m relationship
	Type           string      `json:"type" validate:"omitempty,lte=200"`
	Name           string      `json:"name" validate:"omitempty,lte=200"`
	Description    string      `json:"description" validate:"omitempty,lte=200"`
	Keywords       string      `json:"keywords" validate:"omitempty,lte=200"`
	ImageUrl       string      `json:"imageUrl" validate:"omitempty,url,lte=200"`
	ImageAlt       string      `json:"imageAlt" validate:"omitempty,lte=200"`
	LinkUrl        string      `json:"linkUrl" validate:"omitempty,url,lte=200"`
	LinkText       string      `json:"linkText" validate:"omitempty,lte=200"`
	RawMaterialIds []uuid.UUID `json:"rawMaterialIds" validate:"dive,omitempty"`
	PartIds        []uuid.UUID `json:"partIds" validate:"dive,omitempty"`
	ProductIds     []uuid.UUID `json:"productIds" validate:"dive,omitempty"`
}
type Brand struct {
	Group
}
type BrandNodes struct {
	Nodes []Brand `json:"brandNodes" validate:"dive"`
}

func (s *BrandNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}
func (s *BrandNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	return err
}
func (s *BrandNodes) Upsert(db *sql.DB, userId uuid.UUID) (err error) {
	// check if struct is empty
	if s.Nodes == nil {
		fmt.Println("BrandNodes struct == nil")
		return err
	}
	// construct the sql upsert statement
	qstr := `
        INSERT INTO groups (
            id, user_id, parent_id, type, name, description, keywords, image_url,
            image_alt, link_url, link_text
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (id) DO UPDATE
        SET user_id = $2,
            parent_id = $3,
            type = $4,
            name = $5,
            description = $6,
            keywords = $7,
            image_url = $8,
            image_alt = $9,
            link_url = $10,
            link_text = $11
        WHERE groups.id = $1;`
	// execute it
	for _, node := range s.Nodes {
		_, err = db.Exec(
			qstr, node.Id, userId, node.ParentId, "brand", node.Name,
			node.Description, node.Keywords, node.ImageUrl, node.ImageAlt,
			node.LinkUrl, node.LinkText)
		// check m2m relationships for each node
		// retrieve userId from the request
		if node.RawMaterialIds != nil {
			err = JoinTableUpsert(
				db,
				node.RawMaterialIds,
				"join_group_item", "group_id", "item_id",
				node.Id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
		if node.PartIds != nil {
			err = JoinTableUpsert(
				db,
				node.PartIds,
				"join_group_item", "group_id", "item_id",
				node.Id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
		if node.ProductIds != nil {
			err = JoinTableUpsert(
				db,
				node.ProductIds,
				"join_group_item", "group_id", "item_id",
				node.Id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	return err
}

func (s *BrandNodes) PostgresDelete(db sql.DB) (err error) {
	fmt.Println("Not implemented.")

	//// check if struct is empty
	//if s.Nodes == nil {
	//fmt.Println("BrandNodes struct == nil")
	//return err
	//}
	//// construct the sql upsert statement

	//// execute it

	return err
}

type Category struct {
	Group
}
type CategoryNodes struct {
	CategoryNodes []Category `json:"categoryNodes" validate:"dive"`
}

func (s *CategoryNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Department struct {
	Group
}
type DepartmentNodes struct {
	DepartmentNodes []Department `json:"departmentNodes" validate:"dive"`
}

func (s *DepartmentNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Tag struct {
	Group
}
type TagNodes struct {
	TagNodes []Tag `json:"tagNodes" validate:"dive"`
}

func (s *TagNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
