package models

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type Unit struct {
	Id              string `json:"id" validate:"omitempty,uuid4"`
	SvUserId        string `json:"svUserId" validate:"omitempty,uuid4"`
	Singular        string `json:"singular" validate:"omitempty,lte=50"`
	SingularDisplay string `json:"singularDisplay" validate:"omitempty,lte=50"`
	Plural          string `json:"plural" validate:"omitempty,lte=50"`
	PluralDisplay   string `json:"pluralDisplay" validate:"omitempty,lte=50"`
}
type UnitNodes struct {
	Nodes []Unit `json:"unitNodes" validate:"dive"`
}

func (s *UnitNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *UnitNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("UnitNodes.Validate() ", err)
	}
	return err
}

func (s *Unit) Process() {
	log.Print("Unit.Process() Not iplemented.")
}

func (s *Unit) Upsert(db *sql.DB) (err error) {

	qstr := `
        INSERT INTO unit (
            id, sv_user_id, singular, singular_display, plural, plural_display
        )
        VALUES (
            $1,  -- id
            $2,  -- sv_user_id
            $3,  -- singular
            $4,  -- singular_display
            $5,  -- plural
            $6   -- plural_display
        )
        ON CONFLICT (id) DO UPDATE
        SET singular=$3,
            singular_display=$4,
            plural=$5,
            plural_display=$5
        WHERE unit.id = $1;`

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Singular,
		s.SingularDisplay, s.Plural, s.PluralDisplay,
	)

	if err != nil {
		log.Print("Unit.Upsert()", err)
	}

	return err
}

func (s *Unit) ForeignKeyUpdate(db *sql.DB) (err error) {
	log.Print("Unit.ForeignKeyUpdate() Not implemented.")
	return err
}

func (s *Unit) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("Unit.RelatedTableUpsert() Not implemented.")
	return err
}

type PriceClass struct {
	Id       string `json:"id"`
	SvUserId string `json:"svUserId" validate:"omitempty,uuid4"`
	Type     string `json:"type" validate:"omitempty,lte=100,oneof=grossMargin markup fixed"`
	Name     string `json:"name" validate:"omitempty,lte=100"`
	Amount   uint32 `json:"amount" validate:"omitempty,number"`
	Note     string `json:"note" validate:"omitempty,lte=100"`
}
type PriceClassNodes struct {
	Nodes []PriceClass `json:"priceClassNodes" validate:"dive"`
}

func (s *PriceClassNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *PriceClassNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("PriceClassNodes.Validate() ", err)
	}
	return err
}

func (s *PriceClass) Process() {
	log.Print("PriceClass.Process() Not iplemented.")
}

func (s *PriceClass) Upsert(db *sql.DB) (err error) {

	qstr := `
        INSERT INTO price_class (
            id, sv_user_id, type, name, amount, note
        )
        VALUES (
            $1,  -- id
            $2,  -- sv_user_id
            $3,  -- type
            $4,  -- name
            $5,  -- amount
            $6  -- notes
        )
        ON CONFLICT (id) DO UPDATE
        SET type=$3, name=$4, amount=$5, note=$6
        WHERE price_class.id = $1;`

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Type,
		s.Name, s.Amount, s.Note,
	)

	if err != nil {
		log.Print("PriceClass.Upsert()", err)
	}
	return err
}

func (s *PriceClass) ForeignKeyUpdate(db *sql.DB) (err error) {
	log.Print("Unit.ForeignKeyUpdate() Not implemented.")
	return err
}

func (s *PriceClass) RelatedTableUpsert(db *sql.DB) (err error) {
	log.Print("Unit.RelatedTableUpsert() Not implemented.")
	return err
}

type Identifiers struct {
	MPN        string `json:"mpn" validate:"omitempty,lte=100"`
	CustomId1  string `json:"customId1" validate:"omitempty,lte=100"`
	CustomId2  string `json:"customId2" validate:"omitempty,lte=100"`
	Gtin08Code string `json:"gtin08Code" validate:"omitempty,eq=8"`
	Gtin08Img  string `json:"gtin08Img" validate:"omitempty,url"`
	Gtin12Code string `json:"gtin12Code" validate:"omitempty,len=12"`
	Gtin12Img  string `json:"gtin12Img" validate:"omitempty,url"`
	Gtin13Code string `json:"gtin13Code" validate:"omitempty,len=13"`
	Gtin13Img  string `json:"gtin13Img" validate:"omitempty,url"`
	Gtin14Code string `json:"gtin14Code" validate:"omitempty,len=14"`
	Gtin14Img  string `json:"gtin14Img" validate:"omitempty,url"`
	IsbnCode   string `json:"isbnCode" validate:"omitempty,isbn"`
	IsbnImg    string `json:"isbnImg" validate:"omitempty,url"`
	Isbn10Code string `json:"isbn10Code" validate:"omitempty,isbn10"`
	Isbn10Img  string `json:"isbn10Img" validate:"omitempty,url"`
	Isbn13Code string `json:"isbn13Code" validate:"omitempty,isbn13"`
	Isbn13Img  string `json:"isbn13Img" validate:"omitempty,url"`
}

type Item struct {
	ImageNodes
	Identifiers

	Id                string `json:"id" validate:"required,uuid4"`
	SvUserId          string `json:"svUserId" validate:"omitempty,uuid4"`
	ParentId          string `json:"parentId" validate:"omitempty,uuid4"`
	LocationId        string `json:"locationId" validate:"omitempty,uuid4"`
	PriceClassId      string `json:"priceClassId" validate:"omitempty,uuid4"`
	UnitId            string `json:"unitId" validate:"omitempty,uuid4"`
	Type              string `json:"type" validate:"omitempty,lte=100,oneof=rawMaterial part product"`
	IsVariable        uint8  `json:"isVariable" validate:"omitempty,number,oneof=0 1"`
	IsBundle          uint8  `json:"isBundle" validate:"omitempty,number,oneof=0 1"`
	Position          uint8  `json:"position" validate:"omitempty,number"`
	SKU               string `json:"sku" validate:"omitempty,lte=100"`
	Name              string `json:"name" validate:"omitempty,lte=200"`
	Description       string `json:"description" validate:"omitempty,lte=200"`
	Keywords          string `json:"keywords" validate:"omitempty,lte=200"`
	Cost              uint32 `json:"cost" validate:"omitempty,number"`
	CostOverride      uint32 `json:"costOverride" validate:"omitempty,number"`
	Price             uint32 `json:"price" validate:"omitempty,number"`
	PriceOverride     uint32 `json:"priceOverride" validate:"omitempty,number"`
	PriceDiscount     uint32 `json:"priceDiscount" validate:"omitempty,number"`
	PriceIsFixed      uint32 `json:"priceIsFixed" validate:"omitempty,number,oneof=0 1"`
	QuantityAvailable uint8  `json:"quantityAvailable" validate:"omitempty,number"`
	QuantityMin       uint8  `json:"quantityMin" validate:"omitempty,number"`
	QuantityMax       uint8  `json:"quantityMax" validate:"omitempty,number"`
	Discount          uint8  `json:"discount" validate:"omitempty,number"`
	// sale information ?

	// physical properties
	Length uint32 `json:"length" validate:"omitempty,number"`
	Width  uint32 `json:"width" validate:"omitempty,number"`
	Height uint32 `json:"height" validate:"omitempty,number"`
	Weight uint32 `json:"weight" validate:"omitempty,number"`
	// digital properties
	FileName           string `json:"fileName" validate:"omitempty,lte=100"`
	FilePath           string `json:"filePath" validate:"omitempty,url,lte=100"`
	DownloadCode       string `json:"downloadCode" validate:"omitempty,lte=100"`
	DownloadExpiration string `json:"downloadExpiration" validate:"omitempty,datetime=2006-01-02,lte=100"`

	ClusterIds []string `json:"clusterIds" validate:"dive,omitempty,uuid4"`
	VariationNodes
}

type VariationNodes struct {
	Nodes []Cluster `json:"variationNodes" vaidate:"dive"`
}

type ItemNodes struct {
	Nodes []Item `json:"itemNodes" vaidate:"dive"`
}

func (s *ItemNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}

func (s *ItemNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		log.Print("ItemNodes.Validate() ", err)
	}
	return err
}

func (s *Item) Process() {
	log.Print("Item.Process() Not iplemented.")
}

func (s *Item) Upsert(db *sql.DB) (err error) {
	// check if struct is empty

	qstr := `
        INSERT INTO item (
            id, type, is_variable, is_bundle, position, sku, name, description,
            keywords, cost, cost_override, price, price_override, price_discount,
            price_is_fixed, quantity_available, quantity_min, quantity_max,
            length, width, height, weight, file_name, file_path, download_code,
            download_expiration 
        )
        VALUES (
            $1,  -- id
            $2,  -- type
            $3,  -- isVariable
            $4,  -- isBundle
            $5,  -- position
            $6,  -- sku
            $7,  -- product name
            $8,  -- description
            $9,  -- keywords
            $10, -- cost
            $11, -- cost_override
            $12, -- price
            $13, -- price_ooverride
            $14, -- price_discount
            $15, -- price_is_fixed
            $16, -- quantity_avilable
            $17, -- quantity_min
            $18, -- quaantity_max
            $19, -- length
            $20, -- width
            $21, -- height
            $22, -- weight
            $23, -- file_name
            $24, -- file_path
            $25, -- download_code
            $26  -- download_expiration
        )
        ON CONFLICT (id) DO UPDATE
        SET type=$2,
            is_variable=$3,
            is_bundle=$4,
            position=$5,
            sku=$6,
            name=$7,
            description=$8,
            keywords=$9,
            cost=$10,
            cost_override=$11,
            price=$12,
            price_override=$13,
            price_discount=$14,
            price_is_fixed=$15,
            quantity_available=$16,
            quantity_min=$17,
            quantity_max=$18,
            length=$19,
            width=$20,
            height=$21,
            weight=$22,
            file_name=$23,
            file_path=$24,
            download_code=$25,
            download_expiration=$26
        WHERE item.id = $1;`
	// execute it

	_, err = db.Exec(
		qstr, FormatUUID(s.Id), s.Type, s.IsVariable, s.IsBundle, s.Position,
		s.SKU, s.Name, s.Description, s.Keywords, s.Cost, s.CostOverride,
		s.Price, s.PriceOverride, s.PriceDiscount, s.PriceIsFixed,
		s.QuantityAvailable, s.QuantityMin, s.QuantityMax,
		s.Length, s.Width, s.Height, s.Weight, s.FileName, s.FilePath,
		s.DownloadCode, s.DownloadExpiration,
	)

	if err != nil {
		log.Print("Item.Upsert()", err)
	}
	return err
}

func (s *Item) ForeignKeyUpdate(db *sql.DB) (err error) {
	qstr := `
        UPDATE item
        SET sv_user_id = $2,
            parent_id = $3,
            location_id = $4,
            unit_id = $5,
            price_class_id = $6
        WHERE id = $1;`

	_, err = db.Exec(qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId),
		FormatUUID(s.ParentId), FormatUUID(s.LocationId), FormatUUID(s.UnitId),
		FormatUUID(s.PriceClassId),
	)
	if err != nil {
		log.Print("Item.ForeignKeyUpdate()", err)
	}
	return err
}

func (s *Item) RelatedTableUpsert(db *sql.DB) (err error) {

	if s.ClusterIds != nil {
		for _, id := range s.ClusterIds {
			err = JoinClusterItemUpsert(
				db,
				s.SvUserId,
				id,   // cluster
				s.Id, // item
				s.Position,
			)
			if err != nil {
				log.Print("Item.RelatedTableUpsert() 01 ", err)
			}
		}
	}

	if s.VariationNodes.Nodes != nil {
		for _, node := range s.VariationNodes.Nodes {
			err = JoinClusterItemUpsert(
				db,
				s.SvUserId,
				node.Id, // cluster
				s.Id,    // item
				node.Position,
			)
			if err != nil {
				log.Print("Item.RelatedTableUpsert() 02 ", err)
			}
		}
	}

	return err
}
