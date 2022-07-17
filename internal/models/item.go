package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
)


//type Identifiers struct {
//Id       string `json:"id" validate:"omitempty,uuid4"`
//SvUserId string
//Document string

////MPN        string `json:"mpn" validate:"omitempty,lte=100"`
////CustomId1  string `json:"customId1" validate:"omitempty,lte=100"`
////CustomId2  string `json:"customId2" validate:"omitempty,lte=100"`
////Gtin08Code string `json:"gtin08Code" validate:"omitempty,eq=8"`
////Gtin08Img  string `json:"gtin08Img" validate:"omitempty,url"`
////Gtin12Code string `json:"gtin12Code" validate:"omitempty,len=12"`
////Gtin12Img  string `json:"gtin12Img" validate:"omitempty,url"`
////Gtin13Code string `json:"gtin13Code" validate:"omitempty,len=13"`
////Gtin13Img  string `json:"gtin13Img" validate:"omitempty,url"`
////Gtin14Code string `json:"gtin14Code" validate:"omitempty,len=14"`
////Gtin14Img  string `json:"gtin14Img" validate:"omitempty,url"`
////IsbnCode   string `json:"isbnCode" validate:"omitempty,isbn"`
////IsbnImg    string `json:"isbnImg" validate:"omitempty,url"`
////Isbn10Code string `json:"isbn10Code" validate:"omitempty,isbn10"`
////Isbn10Img  string `json:"isbn10Img" validate:"omitempty,url"`
////Isbn13Code string `json:"isbn13Code" validate:"omitempty,isbn13"`
////Isbn13Img  string `json:"isbn13Img" validate:"omitempty,url"`
//}

type Item struct {
	BaseData
    AllIdNodes
	PlaceId   string `json:"placeId" validate:"omitempty"`

	//Id                string `json:"id" validate:"omitempty,uuid4"`
	//Type              string `json:"type" validate:"omitempty,lte=100,oneof=rawMaterial part product"`
	//SvUserId          string `json:"svUserId" validate:"omitempty,uuid4"`
	//ParentId          string `json:"parentId" validate:"omitempty,uuid4"`
	//LocationId        string `json:"locationId" validate:"omitempty,uuid4"`
	//PriceClassId      string `json:"priceClassId" validate:"omitempty,uuid4"`
	//UnitId            string `json:"unitId" validate:"omitempty,uuid4"`
	//IsVariable        uint8  `json:"isVariable" validate:"omitempty,number,oneof=0 1"`
	//IsBundle          uint8  `json:"isBundle" validate:"omitempty,number,oneof=0 1"`
	//Position          uint8  `json:"position" validate:"omitempty,number"`
	//SKU               string `json:"sku" validate:"omitempty,lte=100"`
	//Name              string `json:"name" validate:"omitempty,lte=200"`
	//Description       string `json:"description" validate:"omitempty,lte=200"`
	//Keywords          string `json:"keywords" validate:"omitempty,lte=200"`
	//Cost              uint32 `json:"cost" validate:"omitempty,number"`
	//CostOverride      uint32 `json:"costOverride" validate:"omitempty,number"`
	//Price             uint32 `json:"price" validate:"omitempty,number"`
	//PriceOverride     uint32 `json:"priceOverride" validate:"omitempty,number"`
	//PriceDiscount     uint32 `json:"priceDiscount" validate:"omitempty,number"`
	//PriceIsFixed      uint32 `json:"priceIsFixed" validate:"omitempty,number,oneof=0 1"`
	//QuantityAvailable uint8  `json:"quantityAvailable" validate:"omitempty,number"`
	//QuantityMin       uint8  `json:"quantityMin" validate:"omitempty,number"`
	//QuantityMax       uint8  `json:"quantityMax" validate:"omitempty,number"`
	//Discount          uint8  `json:"discount" validate:"omitempty,number"`
	//// sale information ?

	//// physical properties
	//Length uint32 `json:"length" validate:"omitempty,number"`
	//Width  uint32 `json:"width" validate:"omitempty,number"`
	//Height uint32 `json:"height" validate:"omitempty,number"`
	//Weight uint32 `json:"weight" validate:"omitempty,number"`
	//// digital properties
	//FileName           string `json:"fileName" validate:"omitempty,lte=100"`
	//FilePath           string `json:"filePath" validate:"omitempty,url,lte=100"`
	//DownloadCode       string `json:"downloadCode" validate:"omitempty,lte=100"`
	//DownloadExpiration string `json:"downloadExpiration" validate:"omitempty,datetime=2006-01-02,lte=100"`

	//Identifiers
	//Image
	//Images
	//ImageIds []string
	//CollectionNodes
	//CollectionIds []string `json:"collectionIds" validate:"dive,omitempty,uuid4"`
}

type ItemNodes struct {
	Nodes []*Item `json:"itemNodes" validate:"dive"`
	Gjson gjson.Result
}

func (s *ItemNodes) Load(fileBuffer *[]byte) (err error) {
	value := gjson.Get(string(*fileBuffer), "itemNodes")
	s.Gjson = value
	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("ItemNodes.Load() %s", err)
	}
	return nil
}

func (s *ItemNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("ItemNodes.Validate() %s", err)
	}
	return nil
}

func (s *ItemNodes) Upsert(accountId string, db *sql.DB) (err error) {
	for i, node := range s.Nodes {
		node.Document = s.Gjson.Array()[i].String()
		qstr := `
            INSERT INTO item (id, account_id, type, document)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET account_id = $2,
                type = $3,
                document = $4
            WHERE item.id = $1;`
		_, err = db.Exec(qstr, node.Id, accountId, node.Type, node.Document)
		if err != nil {
			return fmt.Errorf("ItemNodes.Upsert() %s", err)
		}
	}
	return nil
}

func (s *ItemNodes) ForeignKeyUpdate(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
        var qstr string
        if node.ParentId != "" {
            qstr = `
                UPDATE item
                SET parent_id = $2
                WHERE item.id = $1;`

            _, err = db.Exec(qstr, node.Id, node.ParentId)

		} else {
            // need this query to set fk uuid to null
            qstr = `
                UPDATE item
                SET parent_id = null
                WHERE item.id = $1;`
            _, err = db.Exec(qstr, node.Id)
        }

        if err != nil {
            return fmt.Errorf("ItemNodes.ForeignKeyUpdate() 01 %s", err)
        }
    }
    return nil
}

func (s *ItemNodes) RelatedTableUpsert(accountId string, db *sql.DB) (err error) {
    for i, node := range s.Nodes {
        structArray := []Upserter{}
        ascendentColumn := "item_id"
        if node.CollectionIdNodes.Nodes != nil {
            node.collectionJson = s.Gjson.Array()[i].Get("collectionIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        if node.ContentIdNodes.Nodes != nil {
            node.contentJson = s.Gjson.Array()[i].Get("contentIdNodes")
            node.ContentIdNodes.ascendentColumn = ascendentColumn
            node.ContentIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ContentIdNodes)
        }
        if node.ImageIdNodes.Nodes != nil {
            node.imageJson = s.Gjson.Array()[i].Get("imageIdNodes")
            node.ImageIdNodes.ascendentColumn = ascendentColumn
            node.ImageIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.ImageIdNodes)
        }
        //if node.ItemIdNodes.Nodes != nil {
            //node.itemJson = s.Gjson.Array()[i].Get("itemIdNodes")
            //node.ItemIdNodes.ascendentColumn = ascendentColumn
            //node.ItemIdNodes.ascendentNodeId = node.Id
            //structArray = append(structArray, &node.ItemIdNodes)
        //}
        if node.PlaceIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("placeIdNodes")
            node.PlaceIdNodes.ascendentColumn = ascendentColumn
            node.PlaceIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PlaceIdNodes)
        }
        if node.PersonIdNodes.Nodes != nil {
            node.placeJson = s.Gjson.Array()[i].Get("personIdNodes")
            node.PersonIdNodes.ascendentColumn = ascendentColumn
            node.PersonIdNodes.ascendentNodeId = node.Id
            structArray = append(structArray, &node.PersonIdNodes)
        }
        for _, sa := range structArray {
            err = UpsertHandler(sa, accountId, db)
            if err != nil {
                return fmt.Errorf("ItemNodes.RelatedTableUpsert %s", err)
            }

        }
    }
	return nil
}

func (s *ItemNodes) Delete(db *sql.DB) (err error) {
	for _, node := range s.Nodes {
		qstr := `
            DELETE FROM item
            WHERE item.id = $1;
            `
		_, err = db.Exec(qstr, node.Id)
		if err != nil {
			return fmt.Errorf("ItemNodes.Delete %s", err)
		}
	}
	return nil
}
