package loader

// Consider userId -- will need to coordinate between state and persistent data

type Menu struct {
	// will likely need a locationId (website)
	ParentId    string `json:"parentId"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Device      string `json:"device"` // desktop, tablet, mobile
	IsPrimary   uint8  `json:"isPrimary"`
	IsSecondary uint8  `json:"isSecondary"`
	IsTertiary  uint8  `json:"isTertiary"`
	Order       uint8  `json:"order"`
	// may also need content such as pages or blogs
	Groups []struct {
		Id    string `json:"id"`
		Order uint8  `json:"order"`
	} `json:"groups"`
}
type Menus []struct {
	Menu
}
type Group struct {
	Parent string `json:"parent" validate:"omitempty,uuid4"`
	Id     string `json:"id" validate:"omitempty,uuid4"`
	// Type in the database eg. rawMaterialsTag, partTag, productTag
	// will need join table m2m relationship
	Name           string   `json:"name" validate:"omitempty,lte=200"`
	Description    string   `json:"description" validate:"omitempty,lte=200"`
	Keywords       string   `json:"keywords" validate:"omitempty,lte=200"`
    ImageUrl       string   `json:"imageUrl" validate:"omitempty,url"`
	ImageAlt       string   `json:"imageAlt" validate:"omitempty,lte=200"`
	LinkUrl        string   `json:"linkUrl" validate:"omitempty,url"`
	LinkText       string   `json:"linkText" validate:"omitempty,lte=200"`
	RawMaterialIds []string `json:"rawMaterialIds" validate:"dive,omitempty,uuid4"`
	PartIds        []string `json:"partIds" validate:"dive,omitempty,uuid4"`
	ProductIds     []string `json:"productIds" validate:"dive,omitempty,uuid4"`
}
type Unit struct {
	Id              string `json:"id"`
	Singular        string `json:"singular"`
	Plural          string `json:"plural"`
	DisplaySingular string `json:"displaySingular"`
	DisplayPlural   string `json:"displayPlural"`
}
type Units []struct {
	Unit
}

type Item struct {
	UnitId       string `json:"unitId"`
	ParentSKU    string `json:"Parentsku"`
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Keywords     string `json:"keywords"`
	Cost         uint32 `json:"cost"`
	CostOverride uint32 `json:"costOverride"`
	Groups       struct {
		BrandIds      []string `json:"brandIds"`
		CategoryIds   []string `json:"categoryIds"`
		DepartmentIds []string `json:"departmentIds"`
		TagIds        []string `json:"tagIds"`
	} `json:"groups"`
	Measurements struct {
		Length uint32 `json:"length"`
		Width  uint32 `json:"width"`
		Height uint32 `json:"height"`
		Weight uint32 `json:"weight"`
	} `json:"measurements"`
	Identifiers struct {
		ManufacturersPartNumber string `json:"manufacturersPartNumber"`
		CustomId1               string `json:"customId1"`
		CustomId2               string `json:"customId2"`
		Gtin12Code              string `json:"gtin12Code"`
		Gtin12Img               string `json:"gtin12Img"`
		Gtin13Code              string `json:"gtin13Code"`
		Gtin13Img               string `json:"gtin13Img"`
		IsbnCode                string `json:"isbnCode"`
		IsbnImg                 string `json:"isbnImg"`
	} `json:"identifiers"`
}
type Person struct {
	Id         string `json:"id"`
	Salutation string `json:"salutation"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
	Address    struct {
		Address
	} `json:"Address"`
	BillingAddress struct {
		Address
	} `json:"billingAddress"`
	ShippingAddress struct {
		Address
	} `json:"shippingAddress"`
}
type Address struct {
	Id      string `json:"id"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}
type Location struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Address     struct {
		Address
	} `json:"address"`
	BillingAddress struct {
		Address
	} `json:"billingAddress"`
	ShippingAddress struct {
		Address
	} `json:"shippingAddress"`
}

type Department struct {
	Group
}
type Departments []struct {
	Department
}
type Brand struct {
	Group
}
type Brands []struct {
	Brand
}
type Category struct {
	Group
}
type Categories []struct {
	Category
}
type Tag struct {
	Group
}
type Tags []struct {
	Tag
}
type RawMaterial struct {
	Item
}
type RawMaterials []struct {
	RawMaterial
}
type Part struct {
	Item
}
type Parts []struct {
	Part
}
type PriceClass struct {
	Id string `json:"id"`
	// flat rate, gross margin, markup
	MultiplierType string `json:"multiplierType"`
	Amount         uint32 `json:"amount"`
}
type PriceClasses []struct {
	PriceClass
}
type Attribute struct {
	Group
	Terms []struct {
		TermId string `json:"termId"`
		Order  uint8  `json:"order"`
	} `json:"terms"`
	TermIds []string `json:"termids"`
}
type Attributes []struct {
	Attribute
}
type Term struct {
	Group
}
type Terms []struct {
	Term
}
type Product struct {
	Item
	PriceClassID  string `json:"priceClassId"`
	Price         uint32 `json:"price"`
	PriceOverride uint32 `json:"priceOverride"`
	Bundle        struct {
		FixedPrice      uint32 `json:"fixedPrice"`
		BundledProducts []struct {
			ProductId   string `json:"productId"`
			MinQuantity uint8  `json:"minQuantity"`
			MaxQuantity uint8  `json:"maxQuantity"`
			Discount    uint8  `json:"discount"`
		} `json:"bundledProducts"`
	} `json:"bundle"`
	DigitalAttributes struct {
		Id                 string `json:"id"`
		FileName           string `json:"fileName"`
		FilePath           string `json:"filePath"`
		DownloadCode       string `json:"downloadCode"`
		DownloadExpiration string `json:"downloadExpiration"`
	} `json:"digitalAttributes"`
	Variable struct {
		Attributes []struct {
			Attribute
			IsVariation uint8 `json:"isVariation"`
			IsDisplay   uint8 `json:"isDisplay"`
			Order       uint8 `json:"order"`
			Terms       []struct {
				Term
			} `json:"terms"`
		}
		Variations []struct {
		} `json:"variations"`
	} `json:"variation"`
}
type Products []struct {
	Product
}

// people
type Customer struct {
	Person
}
type Customers []struct {
	Customer
}
type Contact struct {
	Person
}
type Contacts []struct {
	Contact
}
type Company struct {
	Location
	Contacts []struct {
		Person
	} `json:"contacts"`
	ContactIds []string `json:"contactIds"`
}
type Companies []struct {
	Company
}
type Store struct {
	Location
	Contacts []struct {
		Person
	} `json:"contacts"`
	ContactIds []string `json:"contactIds"`
}
type Stores []struct {
	Store
}
type Warehouse struct {
	Location
}
type Warehouses struct {
	Warehouse
}
type Website struct {
	Location
	Domain string `json:"domain"`
}
type Websites []struct {
	Website
}
