package models

import (
	"database/sql"
	"encoding/json"
    "errors"
    "fmt"

	"github.com/go-playground/validator/v10"
)

type SvCollection struct { 
	Id       string `json:"id" validate:"omitempty,uuid4"`
	Type        string   `json:"type" validate:"omitempty,lte=20,oneof=attribute variation brand department rawMaterialCategory partCategory productCategory postCategory pageCategory docCategory rawMaterialTag partTag productTag postTag pageTag docTag"`
	SvUserId string
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	// Type in the database eg. rawMaterialTag, partTag, productTag, productBrand
	// will need join table m2m relationship
	Position uint8  `json:"position" validate:"omitempty,number"`
	Name        string   `json:"name" validate:"omitempty,lte=200"`
	Description string   `json:"description" validate:"omitempty,lte=200"`
	Keywords    string   `json:"keywords" validate:"omitempty,lte=200"`
	LinkUrl     string   `json:"linkUrl" validate:"omitempty,url,lte=200"`
	LinkText    string   `json:"linkText" validate:"omitempty,lte=200"`

    ImageIds []string `json:"imageIds" validate:"dive,omitempty"`
	ItemIds     []string `json:"itemIds" validate:"dive,omitempty,uuid4"`

}

func (s *SvCollection) Process(userId string) (err error) {
    if userId == "" {
        return errors.New("Collection.Process() userId is required.")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvCollection) Upsert(db *sql.DB) (err error) {
    // is this right?
	if s == (&SvCollection{}) {
        return nil
	}
	qstr := `
        INSERT INTO collection (
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
        WHERE collection.id = $1;`

    _, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId),
		s.Type, s.Name, s.Description, s.Keywords, s.LinkUrl,
		s.LinkText)

	if err != nil {
        return fmt.Errorf("Collection.Upsert() %s", err)
	}
    return nil
}

func (s *SvCollection) ForeignKeyUpdate(db *sql.DB) (err error) {
	if s.ParentId == "" {
		return nil
	}
	qstr := `
        UPDATE collection
        SET parent_id = $2
        WHERE collection.id = $1;`

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.ParentId),
	)
	if err != nil {
		return fmt.Errorf("Collection.ForeignKeyUpdate() %s", err)
	}
    return nil
}

func (s *SvCollection) RelatedTableUpsert(db *sql.DB) (err error) {
    if s.ItemIds != nil {
        for _, id := range s.ItemIds {
            err = JoinCollectionItemUpsert(
                db,
                s.SvUserId,
                s.Id,
                id,
                s.Position,
            )
        }
        if err != nil {
            return fmt.Errorf("Collection.RelatedTableUpsert() 01 %s", err)
        }
    }
    return nil
}

func (s *SvCollection) Delete(db *sql.DB) (err error) {
	fmt.Print("Collection.Delete() Not implemented.")
    if err != nil {
        return fmt.Errorf("Collection.Delete() %s", err)
    }
    return nil
}

type Collection struct {
    SvCollection `json:"collection"`
}

func (s *Collection) Load(fileBuffer *[]byte) (err error) {
    err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Collection.Load() %s", err)
	}
    return nil
}

func (s *Collection) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Collection.Validate() %s", err)
	}
    return nil
}

type Collections struct {
    Nodes []SvCollection `json:"collections" validate:"dive"`
}

func (s *Collections) Load(fileBuffer *[]byte) (err error) {
    json.Unmarshal(*fileBuffer, &s)
    if err != nil {
        return fmt.Errorf("Collections.Load() %s", err)
    }
    return nil
}

func (s *Collections) Validate() (err error) {
    validate := validator.New()
    err = validate.Struct(s)
    if err != nil {
        return fmt.Errorf("Collections.Validate() %s", err)
    }
    return nil
}

func (s *Collections) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        node.Process(userId)
        if err != nil {
            return fmt.Errorf("Collections.Process() %s", err)
        }
    }
    return nil
}

func (s *Collections) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Collections.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Collections) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Collections.ForeignKeyUpdate() %s", err)
        }
    }
    return nil
}

func (s *Collections) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Collections.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Collections) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Collections.Delete() %s", err)
        }
    }
    return nil
}
