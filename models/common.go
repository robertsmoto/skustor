package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
)

type loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type validater interface {
	Validate() (err error)
}

type loaderValidater interface {
	loader
	validater
}

func LoaderHandler(data loaderValidater, fileBuffer []byte) {
	// this loads and validates the Nodes or []structs
	var err error
	err = data.Load(&fileBuffer)
	if err != nil {
		log.Print("LoaderHandler.Load() ", err)
	}
	err = data.Validate()
	if err != nil {
		log.Print("LoaderHandler.Validate() ", err)
	}
}

type upserter interface {
	Upsert(db *sql.DB) (err error)
}

type relatedTableUpserter interface {
	RelatedTableUpsert(db *sql.DB) (err error)
}

type foreignKeyUpdater interface {
	ForeignKeyUpdate(db *sql.DB) (err error)
}

type upserterForeignKeysRelatedTableUpserter interface {
	upserter
	foreignKeyUpdater
	relatedTableUpserter
}

func UpsertHandler(data upserterForeignKeysRelatedTableUpserter, db *sql.DB) {
	var err error
	err = data.Upsert(db)
	if err != nil {
		log.Print("UpsertHandler.Upsert ", err)
	}
	err = data.ForeignKeyUpdate(db)
	if err != nil {
		log.Print("UpsertHandler.ForeignKeyUpdate ", err)
	}
	err = data.RelatedTableUpsert(db)
	if err != nil {
		log.Print("UpsertHandler.RelatedTableUpsert ", err)
	}
}

func JoinClusterItemUpsert(db *sql.DB, svUserId, clusterId, itemId string, position uint8) (err error) {
	var data []string
	data = append(data, svUserId, clusterId, itemId)
	jid := Md5Hasher(data)

	qstr := `
        INSERT INTO join_cluster_item (id, sv_user_id, cluster_id, item_id, position)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET sv_user_id = $2, cluster_id = $3, item_id = $4, position = $5;
        `
	_, err = db.Exec(
		qstr, jid, FormatUUID(svUserId), FormatUUID(clusterId),
		FormatUUID(itemId), position,
	)
	if err != nil {
		log.Print("JoinClusterItemUpsert() ", err)
	}
	return err
}

func FormatUUID(str string) *string {
	if str == "" {
		return nil
	} else {
		ret := str
		return &ret
	}
}

func Md5Hasher(data []string) (out string) {
	h := md5.New()
	for _, d := range data {
		io.WriteString(h, d)
	}
	out = fmt.Sprintf("%x", h.Sum(nil))
	return out
}
