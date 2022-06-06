package loader

import (
	"encoding/json"
)

type SvData struct {
	Attribute    `json:"attribute"`
	Attributes   `json:"attributes"`
	Brand        `json:"brand"`
	Brands       `json:"brands"`
	Category     `json:"category"`
	Categories   `json:"categories"`
	Company      `json:"company"`
	Companies    `json:"companies"`
	Contact      `json:"contact"`
	Contacts     `json:"contacts"`
	Customer     `json:"customer"`
	Customers    `json:"customers"`
	Department   `json:"department"`
	Departments  `json:"departments"`
	Menu         `json:"menu"`
	Menus        `json:"menus"`
	Part         `json:"part"`
	Parts        `json:"parts"`
	PriceClass   `json:"priceClass"`
	PriceClasses `json:"priceClasses"`
	Product      `json:"product"`
	Products     `json:"products"`
	RawMaterial  `json:"rawMaterial"`
	RawMaterials `json:"rawMaterials"`
	Store        `json:"store"`
	Stores       `json:"stores"`
	Tag          `json:"tag"`
	Tags         `json:"tags"`
	Term         `json:"term"`
	Terms        `json:"terms"`
	Unit         `json:"unit"`
	Units        `json:"units"`
	Warehouse    `json:"warehouse"`
	Warehouses   `json:"warehouses"`
	Website      `json:"website"`
	Websites     `json:"websites"`
}

func (svd *SvData) LoadJson(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &svd)
	return err
}
