package models

import (
	"database/sql"
	//"fmt"
	"github.com/pborman/uuid"
	"log"
)

type Loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type Validater interface {
	Validate() (err error)
}

type LoaderValidater interface {
	Loader
	Validater
}

func LoaderHandler(data LoaderValidater, fileBuffer []byte) (err error) {
	// this loads and validates the Nodes or []structs
	err = data.Load(&fileBuffer)
	err = data.Validate()
	return err
}

type Upserter interface {
	Upsert(db *sql.DB) (err error)
}

type RelatedTableUpserter interface {
	RelatedTableUpsert(db *sql.DB) (err error)
}

type GroupUpserter interface {
	Upserter
	RelatedTableUpserter
}

func UpsertHandler(data GroupUpserter, db *sql.DB) (err error) {
	go data.Upsert(db)
	go data.RelatedTableUpsert(db)
	return err
}

func JoinGroupItemUpsert(db *sql.DB, userId, groupId, itemId string) (err error) {
	qstr := `
        INSERT INTO join_group_item (user_id, group_id, item_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (group_id, item_id)
        DO UPDATE SET user_id = $1, group_id = $2, item_id = $3;
        `
	_, err = db.Exec(qstr, uuid.Parse(userId), uuid.Parse(groupId), uuid.Parse(itemId))
	if err != nil {
		log.Print("Err func JoinTableUpsert() ", err)
	}
	return err
}
