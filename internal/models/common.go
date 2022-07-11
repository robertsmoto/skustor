package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
)

type BaseData struct {
	Id       string `json:"id" validate:"required,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	SvUserId string
	Document string
}

type loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type validater interface {
	Validate() (err error)
}

type upserter interface {
	Upsert(userId string, db *sql.DB) (err error)
}

type foreignKeyUpdater interface {
	ForeignKeyUpdate(db *sql.DB) (err error)
}

type relatedTableUpserter interface {
	RelatedTableUpsert(userId string, db *sql.DB) (err error)
}

type LoaderProcesserUpserter interface {
	loader
	validater
	upserter
	foreignKeyUpdater
	relatedTableUpserter
}

func JsonLoaderUpserterHandler(
	data LoaderProcesserUpserter,
	userId string,
	fileBuffer *[]byte,
	db *sql.DB) (err error) {

	err = data.Load(fileBuffer)
	if err != nil {
		return err
	}
	err = data.Validate()
	if err != nil {
		return err
	}
	err = data.Upsert(userId, db)
	if err != nil {
		return err
	}
	err = data.ForeignKeyUpdate(db)
	if err != nil {
		return err
	}
	err = data.RelatedTableUpsert(userId, db)
	if err != nil {
		return err
	}
	return nil
}

func JoinCollectionItemUpsert(db *sql.DB, svUserId, collectionId, itemId string) (err error) {
	var data []string

	data = append(data, svUserId, collectionId, itemId)
    jid := Md5Hasher(data)

	qstr := `
        INSERT INTO join_collection_item (id, sv_user_id, collection_id, item_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id)
        DO UPDATE SET sv_user_id = $2, collection_id = $3, item_id = $4;
        `
	_, err = db.Exec(qstr, jid, svUserId, collectionId, itemId)
	if err != nil {
		fmt.Println("JoinCollectionItemUpsert() ", err)
	}
	return err
}

func JoinCollectionContentUpsert(db *sql.DB, svUserId, collectionId, contentId string) (err error) {
	var data []string

	data = append(data, svUserId, collectionId, contentId)
    jid := Md5Hasher(data)

	qstr := `
        INSERT INTO join_collection_content (id, sv_user_id, collection_id, content_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id)
        DO UPDATE SET sv_user_id = $2, collection_id = $3, content_id = $4;
        `
	_, err = db.Exec(qstr, jid, svUserId, collectionId, contentId)
	if err != nil {
		fmt.Println("JoinCollectionContentUpsert() ", err)
	}
	return err
}

func Md5Hasher(data []string) (out string) {
	h := md5.New()
	for _, d := range data {
		io.WriteString(h, d)
	}
	out = fmt.Sprintf("%x", h.Sum(nil))
	return out
}
