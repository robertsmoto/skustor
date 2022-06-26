package models

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
)

//type ModelNodes interface {
        //SvUserNodes
        //PlaceNodes
        //PriceClassNodes
        //UnitNodes
        //CollectionNodes
        //ItemNodes

//}

type loader interface {
	Load(fileBuffer *[]byte) (err error)
}

type validater interface {
	Validate() (err error)
}

type processer interface {
    Process(userId string) (err error)

}

type upserter interface {
	Upsert(db *sql.DB) (err error)
}

type foreignKeyUpdater interface {
	ForeignKeyUpdate(db *sql.DB) (err error)
}

type relatedTableUpserter interface {
	RelatedTableUpsert(db *sql.DB) (err error)
}

type LoaderProcesserUpserter interface {
    loader
    validater
    processer
	upserter
	foreignKeyUpdater
	relatedTableUpserter
}

func JsonLoaderUpserterHandler(data LoaderProcesserUpserter, userId string, fileBuffer *[]byte, db *sql.DB) (err error) {
    err = data.Load(fileBuffer)
    err = data.Validate()
    err = data.Process(userId)
	err = data.Upsert(db)
	err = data.ForeignKeyUpdate(db)
	err = data.RelatedTableUpsert(db)
    if err != nil {
        return fmt.Errorf("JsonLoaderUpserterHandler %s", err)
    }
    return nil
}

func JoinCollectionItemUpsert(db *sql.DB, svUserId, clusterId, itemId string, position uint8) (err error) {
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
		fmt.Println("JoinCollectionItemUpsert() ", err)
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
