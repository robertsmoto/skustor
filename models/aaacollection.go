package models

import (
	"database/sql"
	"encoding/json"
    "errors"
    "fmt"

	"github.com/go-playground/validator/v10"
    "github.com/tidwall/gjson"
)

type SvAaaCollection struct { 
	Id       string `json:"id" validate:"omitempty,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
    ImageIds []string `json:"imageIds" validate:"dive,omitempty"`
	ItemIds     []string `json:"itemIds" validate:"dive,omitempty,uuid4"`
	SvUserId string
    Gjson gjson.Result
    Document string
}

func (s *SvAaaCollection) Process(userId string) (err error) {
    if userId == "" {
        return errors.New("Collection.Process() userId is required.")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvAaaCollection) Upsert(db *sql.DB) (err error) {
    fmt.Println(s.Id, s.Document)
    fmt.Println(s)
     //is this right?
    if s == (&SvAaaCollection{}) {
        return nil
    }
    qstr := `
        INSERT INTO aaa_collection (
            id, sv_user_id, document
        )
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id = $2,
            document = $3
        WHERE aaa_collection.id = $1;`

    _, err = db.Exec(
        qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Document,
    )

    if err != nil {
        return fmt.Errorf("Collection.Upsert() %s", err)
    }
    return nil
}

func (s *SvAaaCollection) ForeignKeyUpdate(db *sql.DB) (err error) {
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

func (s *SvAaaCollection) RelatedTableUpsert(db *sql.DB) (err error) {
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
    return nil
}

func (s *SvAaaCollection) Delete(db *sql.DB) (err error) {
	fmt.Print("Collection.Delete() Not implemented.")
    if err != nil {
        return fmt.Errorf("Collection.Delete() %s", err)
    }
    return nil
}

type AaaCollection struct {
    SvAaaCollection `json:"aaaCollection"`
}

func (s *AaaCollection) Load(fileBuffer *[]byte) (err error) {

    // unmarshal JSON
	value := gjson.Get(string(*fileBuffer), "aaaCollection")
    fmt.Printf("## value %T", value)
	fmt.Println("one collection ### --> ", value.String())
    s.Gjson = value


    err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Collection.Load() %s", err)
	}
    return nil
}

func (s *AaaCollection) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Collection.Validate() %s", err)
	}
    return nil
}

type AaaCollections struct {
    Nodes []SvAaaCollection `json:"aaaCollections" validate:"dive"`
    Gjson gjson.Result 
}

func (s *AaaCollections) Load(fileBuffer *[]byte) (err error) {
    value := gjson.Get(string(*fileBuffer), "aaaCollections")
    s.Gjson = value

    json.Unmarshal(*fileBuffer, &s)
    if err != nil {
        return fmt.Errorf("Collections.Load() %s", err)
    }
    return nil
}

func (s *AaaCollections) Validate() (err error) {
    validate := validator.New()
    err = validate.Struct(s)
    if err != nil {
        return fmt.Errorf("Collections.Validate() %s", err)
    }
    return nil
}

func (s *AaaCollections) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        node.Process(userId)
        if err != nil {
            return fmt.Errorf("Collections.Process() %s", err)
        }
    }
    return nil
}

func (s *AaaCollections) Upsert(db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        node.Document = s.Gjson.Array()[i].String()
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Collections.Upsert() %s", err)
        }
    }
    return nil
}

func (s *AaaCollections) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Collections.ForeignKeyUpdate() %s", err)
        }
    }
    return nil
}

func (s *AaaCollections) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Collections.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *AaaCollections) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Collections.Delete() %s", err)
        }
    }
    return nil
}
