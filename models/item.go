package models

import (
	"database/sql"
	"encoding/json"
    "errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SvUnit struct {
	Id              string `json:"id" validate:"omitempty,uuid4"`
	SvUserId        string `json:"svUserId" validate:"omitempty,uuid4"`
	Singular        string `json:"singular" validate:"omitempty,lte=50"`
	SingularDisplay string `json:"singularDisplay" validate:"omitempty,lte=50"`
	Plural          string `json:"plural" validate:"omitempty,lte=50"`
	PluralDisplay   string `json:"pluralDisplay" validate:"omitempty,lte=50"`
}

func (s *SvUnit) Process(userId string) (err error) {
    if userId == "" {
        return errors.New("SvUnit.Process() requires userId.")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvUnit) Upsert(db *sql.DB) (err error) {
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
		return fmt.Errorf("SvUnit.Upsert() %s", err)
	}
    return nil
}

func (s *SvUnit) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("SvUnit.ForeignKeyUpdate() Not implemented.")
    return nil
}

func (s *SvUnit) RelatedTableUpsert(db *sql.DB) (err error) {
	fmt.Print("SvUnit.RelatedTableUpsert() Not implemented.")
    return nil
}

func (s *SvUnit) Delete(db *sql.DB) (err error) {
	fmt.Print("SvUnit.Delete() Not implemented.")
    return nil
}

type Unit struct {
	SvUnit `json:"unit"`
}

func (s *Unit) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Unit.Load() %s", err)
	}
    return nil
}

func (s *Unit) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Unit.Validate() %s", err)
	}
    return nil
}

type Units struct {
	Nodes []Unit `json:"units" validate:"dive"`
}

func (s *Units) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Units.Load() %s", err)
	}
    return nil
}

func (s *Units) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Units.Validate() %s", err)
	}
    return nil
}

func (s *Units) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("Units.Process() %s", err)
        }
    }
    return nil
}

func (s *Units) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Units.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Units) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Units.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Units) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Units.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Units) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Units.Delete() %s", err)
        }
    }
    return nil
}

type SvPriceClass struct {
    Id       string `json:"id" validate:"omitempty,uuid4"`
	Type     string `json:"type" validate:"omitempty,lte=100,oneof=grossMargin markup fixed"`
	SvUserId string `json:"svUserId" validate:"omitempty,uuid4"`
	Name     string `json:"name" validate:"omitempty,lte=100"`
	Amount   uint32 `json:"amount" validate:"omitempty,number"`
	Note     string `json:"note" validate:"omitempty,lte=100"`
}

func (s *SvPriceClass) Process(userId string) (err error) {
    if userId == "" {
        return fmt.Errorf("SvPriceClass.Process() requires userId")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvPriceClass) Upsert(db *sql.DB) (err error) {
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
		return fmt.Errorf("SvPriceClass.Upsert() %s", err)
	}
    return nil
}

func (s *SvPriceClass) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("SvUnit.ForeignKeyUpdate() Not implemented.")
    return nil
}

func (s *SvPriceClass) RelatedTableUpsert(db *sql.DB) (err error) {
	fmt.Println("SvUnit.RelatedTableUpsert() Not implemented.")
    return nil
}

func (s *SvPriceClass) Delete(db *sql.DB) (err error) {
	fmt.Println("SvUnit.Delete() Not implemented.")
    return nil
}

type PriceClass struct {
	SvPriceClass `json:"priceClass"`
}

func (s *PriceClass) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("PriceClass.Load() %s", err)
	}
    return nil
}

func (s *PriceClass) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("PriceClass.Validate() %s", err)
	}
    return nil
}

type PriceClasses struct {
	Nodes []PriceClass `json:"priceClasses" validate:"dive"`
}

func (s *PriceClasses) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("PriceClasses.Load() %s", err)
	}
	return nil
}

func (s *PriceClasses) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("PriceClasses.Validate() %s", err)
	}
	return nil
}

func (s *PriceClasses) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("PriceClasses.Process() %s", err)
        }
    }
    return nil
}

func (s *PriceClasses) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("PriceClasses.Upsert() %s", err)
        }
    }
    return nil
}

func (s *PriceClasses) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("PriceClasses.ForeignKeyUpdate() %s", err)
        }
    }
    return nil
}

func (s *PriceClasses) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("PriceClasses.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *PriceClasses) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("PriceClasses.Delete() %s", err)
        }
    }
    return nil
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

type SvItem struct {

	Id                string `json:"id" validate:"omitempty,uuid4"`
	Type              string `json:"type" validate:"omitempty,lte=100,oneof=rawMaterial part product"`
	SvUserId          string `json:"svUserId" validate:"omitempty,uuid4"`
	ParentId          string `json:"parentId" validate:"omitempty,uuid4"`
	LocationId        string `json:"locationId" validate:"omitempty,uuid4"`
	PriceClassId      string `json:"priceClassId" validate:"omitempty,uuid4"`
	UnitId            string `json:"unitId" validate:"omitempty,uuid4"`
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

	Identifiers
    Image
    Images
    ImageIds []string
    Collection
    Collections
	CollectionIds []string `json:"collectionIds" validate:"dive,omitempty,uuid4"`
}

func (s *SvItem) Process(userId string) (err error) {
    if userId == "" {
        return fmt.Errorf("svItem.Process() requires userId")
    }
    s.SvUserId = userId
    return nil
}

func (s *SvItem) Upsert(db *sql.DB) (err error) {
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
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), s.Type, s.IsVariable, s.IsBundle, s.Position,
		s.SKU, s.Name, s.Description, s.Keywords, s.Cost, s.CostOverride,
		s.Price, s.PriceOverride, s.PriceDiscount, s.PriceIsFixed,
		s.QuantityAvailable, s.QuantityMin, s.QuantityMax,
		s.Length, s.Width, s.Height, s.Weight, s.FileName, s.FilePath,
		s.DownloadCode, s.DownloadExpiration,
	)
	if err != nil {
		return fmt.Errorf("SvItem.Upsert() %s", err)
	}
    return nil
}

func (s *SvItem) ForeignKeyUpdate(db *sql.DB) (err error) {
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
		return fmt.Errorf("SvItem.ForeignKeyUpdate() %s", err)
	}
    return nil
}

func (s *SvItem) RelatedTableUpsert(db *sql.DB) (err error) {
	if s.CollectionIds != nil {
		for _, id := range s.CollectionIds {
			err = JoinCollectionItemUpsert(
				db,
				s.SvUserId,
				id,   // collection
				s.Id, // item
				s.Position,
			)
			if err != nil {
				return fmt.Errorf("SvItem.RelatedTableUpsert() 01 %s", err)
			}
		}
	}
    return nil
}

type Item struct {
	SvItem `json:"item"`
}

func (s *Item) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Item.Load() %s", err)
	}
    return nil
}

func (s *Item) Validate() (err error){
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Item.Validate() %s", err)
	}
    return nil
}

type Items struct {
	Nodes []SvItem `json:"items" vaidate:"dive"`
}

func (s *Items) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Items.Load() %s", err)
	}
    return nil
}

func (s *Items) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Items.Validate() %s", err)
	}
    return nil
}

func (s *Items) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("Items.Process() %s", err)
        }
    }
    return nil
}

func (s *Items) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Items.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Items) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Items.ForeignKeyUpdate() %s", err)
        }
    }
    return nil
}

func (s *Items) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Items.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Items) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Items.Delete() %s", err)
        }
    }
    return nil
}
