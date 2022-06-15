package models

import "encoding/json"

type Unit struct {
	Id              string `json:"id" validate:"omitempty,uuid4"`
	Singular        string `json:"singular" validate:"omitempty,lte=50"`
	Plural          string `json:"plural" validate:"omitempty,lte=50"`
	DisplaySingular string `json:"displaySingular" validate:"omitempty,lte=50"`
	DisplayPlural   string `json:"displayPlural" validate:"omitempty,lte=50"`
}
type UnitNodes struct {
	UnitNodes []Unit `json:"unitNodes" validate:"dive"`
}

func (s *UnitNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type PriceClass struct {
	Id string `json:"id"`
	// flat rate, gross margin, markup
	MultiplierType string `json:"multiplierType"`
	Amount         uint32 `json:"amount"`
}
type PriceClassNodes struct {
	PriceClassNodes []PriceClass `json:"priceClassNodes" validate:"dive"`
}

func (s *PriceClassNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Attribute struct {
	//Used for variable products
	Id          string `json:"id"`
	Name        string `json:"name"`
	IsVariation uint8  `json:"isVariation" validate:"omitempty,eq=0|eq=1"`
	IsDisplay   uint8  `json:"isDisplay" validate:"omitempty,eq=0|eq=1"`
	Order       uint8  `json:"order" validate:"omitempty,number"`
	Terms       []struct {
		Id          string `json:"id" validate:"omitempty,uuid4"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Order       uint8  `json:"order" validate:"omitempty,number"`
	} `json:"terms" validate:"dive"`
}
type AttributeNodes struct {
	AttributeNodes []Attribute `json:"attributeNodes" validate:"dive"`
}

func (s *AttributeNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}


type Item struct {
	Id string `json:"id" validate:"required,uuid4"`
	// user_id
	// identifiers_id
	// measurements_id
	// digital_assetts_id
	// groups -- m2m
	UnitId       string `json:"unitId" validate:"omitempty,uuid4"`
	ParentId     string `json:"parentId" validate:"omitempty,uuid4"`
	SKU          string `json:"sku" validate:"omitempty,lte=100"`
	Name         string `json:"name" validate:"omitempty,lte=200"`
	Description  string `json:"description" validate:"omitempty,lte=200"`
	Keywords     string `json:"keywords" validate:"omitempty,lte=200"`
	Cost         uint32 `json:"cost" validate:"omitempty,number"`
	CostOverride uint32 `json:"costOverride" validate:"omitempty,number"`
	PriceClassID  string `json:"priceClassId" validate:"omitempty,uuid4"`
	Price         uint32 `json:"price" validate:"omitempty,number"`
	PriceOverride uint32 `json:"priceOverride" validate:"omitempty,number"`
	Groups       struct {
		BrandIds      []string `json:"brandIds" validate:"dive,omitempty,uuid4"`
		CategoryIds   []string `json:"categoryIds" validate:"dive,omitempty,uuid4"`
		DepartmentIds []string `json:"departmentIds" validate:"dive,omitempty,uuid4"`
		TagIds        []string `json:"tagIds" validate:"dive,omitempty,uuid4"`
	} `json:"groups"`
	Measurements struct {
		Length uint32 `json:"length" validate:"omitempty,number"`
		Width  uint32 `json:"width" validate:"omitempty,number"`
		Height uint32 `json:"height" validate:"omitempty,number"`
		Weight uint32 `json:"weight" validate:"omitempty,number"`
	} `json:"measurements"`
	Identifiers struct {
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
	} `json:"identifiers"`
	DigitalAttributes struct {
		FileName     string `json:"fileName" validate:"omitempty,lte=100"`
		FilePath     string `json:"filePath" validate:"omitempty,url"`
		DownloadCode string `json:"downloadCode" validate:"omitempty,lte=100"`
		DownloadExp  string `json:"downloadExp" validate:"omitempty,datetime=2006-01-02"`
	} `json:"digitalAttributes"`
	Bundle struct {
		FixedPrice      uint32 `json:"fixedPrice" validate:"omitempty,number"`
		BundledItemsNode []struct {
			ProductId   string `json:"productId" validate:"omitempty,uuid4"`
			MinQuantity uint8  `json:"minQuantity" validate:"omitempty,number"`
			MaxQuantity uint8  `json:"maxQuantity" validate:"omitempty,number"`
			Discount    uint8  `json:"discount" validate:"omitempty,number"`
		} `json:"bundledProducts" validate:"dive"`
	} `json:"bundle"`
	Attributes []struct {
		Attribute
	} `json:"attributes" validate:"dive"`
	Variations []struct {
		ProductId   string `json:"productId"`
		AttributeId string `json:"attributeId"`
		TermId      string `json:"termId"`
		ImageId     string `json:"imageId"`
		MinQuantity uint8  `json:"minQuantity" validate:"omitempty,number"`
		MaxQuantity uint8  `json:"maxQuantity" validate:"omitempty,number"`
		Discount    uint8  `json:"discount" validate:"omitempty,number"`
	} `json:"variations" validate:"dive"`
    ImageNodes
}

type ItemNodes struct {
    Node []Item `json:"itemNode" vaidate:"dive"`
}

func Load(s *ItemNodes) (err error) {
}
