package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"

    "github.com/tidwall/gjson"
)

type BaseData struct {
	Id       string `json:"id" validate:"required,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
    Type string `json:"type" validate:"omitempty,lte=20"`
	AccountId string
	Document string
}

type BaseIdData struct {
    Id string `json:"id" validate:"uuid4"`
    Attributes string `json:"attributes" validate:"omitempty,json"`
}

type CollectionIds struct {
    BaseIdData
}

type CollectionIdNodes struct {
    Nodes []CollectionIds `json:"collectionIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    collectionJson gjson.Result
}

func (s *CollectionIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "collection_id"
    for i, node := range s.Nodes {
        node.Attributes = s.collectionJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("CollectionIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type ContentIds struct {
    BaseIdData
}

type ContentIdNodes struct {
    Nodes []ContentIds `json:"contentIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    contentJson gjson.Result
}

func (s *ContentIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "content_id"
    for i, node := range s.Nodes {
        node.Attributes = s.contentJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("ContentIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type ImageIds struct {
    Id string `json:"id" validate:"omitempty"`
    Attributes string `json:"attributes" validate:"omitempty,json"`
}

type ImageIdNodes struct {
    Nodes []ImageIds `json:"imageIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    imageJson gjson.Result
}

func (s *ImageIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "image_id"
    for i, node := range s.Nodes {
        node.Attributes = s.imageJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("ImageIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type ItemIds struct {
    BaseIdData
}

type ItemIdNodes struct {
    Nodes []ItemIds `json:"itemIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    itemJson gjson.Result
}

func (s *ItemIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "item_id"
    for i, node := range s.Nodes {
        node.Attributes = s.itemJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("ItemIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type PersonIds struct {
    BaseIdData
}

type PersonIdNodes struct {
    Nodes []PersonIds `json:"personIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    personJson gjson.Result
}

func (s *PersonIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "person_id"
    for i, node := range s.Nodes {
        node.Attributes = s.personJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("PersonIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type PlaceIds struct {
    BaseIdData
}

type PlaceIdNodes struct {
    Nodes []PlaceIds `json:"placeIdNodes" validate:"dive,omitempty"`
    ascendentColumn string  // corresponding columns eg. "content_id"
    ascendentNodeId string  // parent object
    placeJson gjson.Result
}

func (s *PlaceIdNodes) Upsert(accountId string, db *sql.DB) (err error) {
    ascendentColumn := "place_id"
    for i, node := range s.Nodes {
        node.Attributes = s.placeJson.Array()[i].Get("attributes").String()
        j := NewJoinTable(
            s.ascendentColumn,
            s.ascendentNodeId,
            ascendentColumn,
            node.Id,
            node.Attributes,
        )
        err = UpsertHandler(j, accountId, db)
        if err != nil {
            return fmt.Errorf("PersonIdNodes.Upsert %s", err)
        }
    }
    return nil
}

type AllIdNodes struct {
    CollectionIdNodes
    ContentIdNodes
    ItemIdNodes
    ImageIdNodes
    PlaceIdNodes
    PersonIdNodes
}

type loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type validater interface {
	Validate() (err error)
}

type LoaderValidator interface {
	loader
	validater
}

type PreProcessor interface {
    PreProcess(accountId string, db *sql.DB) (err error)
}

type Upserter interface {
	Upsert(accountId string, db *sql.DB) (err error)
}

type ForeignKeyUpdater interface {
	ForeignKeyUpdate(db *sql.DB) (err error)
}

type RelatedTableUpserter interface {
	RelatedTableUpsert(accountId string, db *sql.DB) (err error)
}

func LoadValidateHandler(data LoaderValidator, fileBuffer *[]byte) (err error) {
	err = data.Load(fileBuffer)
	if err != nil {
		return fmt.Errorf("LoadValidateHandler 01 %s", err)
	}
	err = data.Validate()
	if err != nil {
		return fmt.Errorf("LoadValidateHandler 01 %s", err)
	}
	return nil
}

func PreProcessHandler(data PreProcessor, accountId string, db *sql.DB) (err error) {
	err = data.PreProcess(accountId, db)
	if err != nil {
		return fmt.Errorf("PreProcessHandler 01 %s", err)
	}
	return nil
}

func UpsertHandler(data Upserter, accountId string, db *sql.DB) (err error) {
	err = data.Upsert(accountId, db)
	if err != nil {
		return fmt.Errorf("UpsertHandler 01 %s", err)
	}
	return nil
}

func ForeignKeyUpdateHandler(data ForeignKeyUpdater, db *sql.DB) (err error) {
	err = data.ForeignKeyUpdate(db)
	if err != nil {
		return fmt.Errorf("ForeignKeyUpdateHandler 01 %s", err)
	}
	return nil
}

type Deleter interface {
    Delete(db *sql.DB) (err error)
}

func DeleteHandler(data Deleter, db *sql.DB) (err error) {
	err = data.Delete(db)
	if err != nil {
		return fmt.Errorf("DeleteHandler 01 %s", err)
	}
	return nil
}


func RelatedTableUpsertHandler(data RelatedTableUpserter, accountId string, db *sql.DB) (err error) {
	err = data.RelatedTableUpsert(accountId, db)
	if err != nil {
		return fmt.Errorf("RelatedTableUpsertHandler 01 %s", err)
	}
	return nil
}

type JoinTable struct {
    col1 []string  // {colName, id}
    col2 []string  // {colName, id}
    attributes string // json document
}

func NewJoinTable(ascCol, ascId, nodeCol, nodeId, nodeAttr string) *JoinTable {
    j := new(JoinTable)
    j.col1 = []string{ascCol, ascId}
    j.col2 = []string{nodeCol, nodeId}
    j.attributes = nodeAttr
    return j
}

func (s *JoinTable) Upsert(accountId string, db *sql.DB) (err error) {
	var strArr []string
	strArr = append(strArr, accountId, s.col1[1], s.col2[1])
    tid := Md5Hasher(strArr)
    if s.attributes == "" {
        s.attributes = "{}"
    }

	qstr := fmt.Sprintf(`
        INSERT INTO joins (id, account_id, %s, %s, attributes)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET account_id = $2, %s = $3, %s = $4, attributes = $5;
        `, s.col1[0], s.col2[0], s.col1[0], s.col2[0])
	_, err = db.Exec(qstr, tid, accountId, s.col1[1], s.col2[1], s.attributes)
	if err != nil {
		return fmt.Errorf("JoinTableUpsert() %s", err)
	}
	return nil
}

func JoinCollectionContentUpsert(db *sql.DB, svUserId, collectionId, contentId string) (err error) {
	var data []string

	data = append(data, svUserId, collectionId, contentId)
    jid := Md5Hasher(data)

	qstr := `
        INSERT INTO join_collection_content (id, sv_account_id, collection_id, content_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id)
        DO UPDATE SET sv_account_id = $2, collection_id = $3, content_id = $4;
        `
	_, err = db.Exec(qstr, jid, svUserId, collectionId, contentId)
	if err != nil {
		return fmt.Errorf("JoinCollectionContentUpsert 01 %s", err)
	}
	return nil
}

func Md5Hasher(data []string) (out string) {
	h := md5.New()
	for _, d := range data {
		io.WriteString(h, d)
	}
	out = fmt.Sprintf("%x", h.Sum(nil))
	return out
}
