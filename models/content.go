package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)

type Content struct {
	SvUserId string
	Id       string `json:"id" validate:"omitempty,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	Document string

	//AuthorIds   []string `json:"authorIds" validate:"dive,omitempty,uuid4"`
	//WebsiteIds  []string `json:"websiteIds" validate:"dive,omitempty,uuid4"`
	//CategoryIds []string `json:"categoryIds" validate:"dive,omitempty,uuid4"`
	//TagIds      []string `json:"tagIds" validate:"dive,omitempty,uuid4"`
	//ImageIds    []string `json:"imageIds" validate:"dive,omitempty,uuid4"`

	//type: article, page, docs
	//PublishedTime string   `json:"publishedTime" validate:"omitempty,datetime=15:04 MST"`

	//Published     string   `json:"published" validate:"omitempty,datetime=2006-01-02"`
	//Modified      string   `json:"modified" validate:"omitempty,datetime=2006-01-02"`
	//Keywords      string   `json:"keywords"`
	//Title         string   `json:"title"`
	//Excerpt       string   `json:"excerpt"`
	//Body          string   `json:"body" validate:"omitempty"`
	//Footer        string   `json:"footer"`
}

type ContentNodes struct {
	Nodes []Content `json:"contentNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *ContentNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "contentNodes")
	s.Gjson = value

	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("ContentNodes.Load() %s", err)
	}
	return nil
}

func (s *ContentNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("ContentNodes.Validate() %s", err)
	}
	return nil
}

func (s *ContentNodes) Upsert(userId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO content (
                id, sv_user_id, document
            )
            VALUES ($1, $2, $3)
            ON CONFLICT (id) DO UPDATE
            SET sv_user_id = $2,
                document = $3
            WHERE content.id = $1;`

		_, err = db.Exec(
			qstr, FormatUUID(node.Id), FormatUUID(userId), node.Document,
		)
		if err != nil {
			return fmt.Errorf("ContentNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *ContentNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		qstr := `
            UPDATE content
            SET parent_id = $2
            WHERE content.id = $1;`

		_, err = db.Exec(
			qstr, FormatUUID(node.Id), FormatUUID(node.ParentId),
		)
		if err != nil {
			return fmt.Errorf("ContentNodes.ForeignKeyUpdate() %s", err)
		}
	}
	return nil
}

func (s *ContentNodes) RelatedTableUpsert(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		fmt.Println("node", node)
		//if s.ItemIds != nil {
		//for _, id := range s.ItemIds {
		//err = JoinContentItemUpsert(
		//db,
		//s.SvUserId,
		//s.Id,
		//id,
		//s.Position,
		//)
		//}
		//if err != nil {
		//return fmt.Errorf("Content.RelatedTableUpsert() 01 %s", err)
		//}
		//}
		if err != nil {
			return fmt.Errorf("ContentNodes.RelatedTableUpsert() %s", err)
		}
	}
	return nil
}

func (s *ContentNodes) Delete(db *sql.DB) (err error) {
	fmt.Println("ContentNodes.Delete() Not implemented.")
	//for _, node := range s.Nodes {
	//if err != nil {
	//return fmt.Errorf("Contents.Delete() %s", err)
	//}
	//}
	return nil
}
