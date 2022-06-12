package models

import (
	"database/sql"
	"fmt"
	"github.com/pborman/uuid"
	"log"
)

type Loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type Validater interface {
	Validate() (err error)
}

type Upserter interface {
	Upsert(db *sql.DB, userId string) (err error)
}

type LoaderValidaterUpserter interface {
	Loader
	Validater
	Upserter
}

func JsonLoadValidateUpsert(
	data LoaderValidaterUpserter, fileBuffer []byte, db *sql.DB, userId string) (err error) {

	err = data.Load(&fileBuffer)
	err = data.Validate()
	err = data.Upsert(db, userId)
	return err
}

type RelatedTableUpserter interface {
	RelatedTableUpsert(db *sql.DB, userId string) (err error)
}

func RelatedTableUpsert(data RelatedTableUpserter, db *sql.DB, userId string) (err error) {
	err = data.RelatedTableUpsert(db, userId)
	return err
}

type ImageSizer interface {
	ImageSize(db *sql.DB, date, userId, userDir string) (err error)
}

func ImageSizeUpsert(data ImageSizer, db *sql.DB, date, userId, userDir string) (err error) {
	data.ImageSize(db, date, userId, userDir)
	return err
}

func JoinTableUpsert(db *sql.DB, q, idArray []string) (err error) {
	for _, varId := range idArray {
		// q {0=table, 1=userId, 2=groupId, 3=col1, 4=col2, 5=col3}
		qstr := fmt.Sprintf(`
            INSERT INTO %s (%s, %s, %s)
            VALUES ($1, $2, $3)
            ON CONFLICT (%s, %s)
            DO UPDATE SET %s = $1, %s = $2, %s = $3;`,
			q[0], q[3], q[4], q[5], q[4], q[5], q[3], q[4], q[5],
		)
		_, err = db.Exec(qstr, uuid.Parse(q[1]), uuid.Parse(q[2]), uuid.Parse(varId))
		if err != nil {
			log.Print("Err func JoinTableUpsert() ", err)
		}
	}
	return err
}
