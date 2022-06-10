package models

import "encoding/json"

type Address struct {
	Id      string `json:"id" validate:"omitempty,uuid4"`
	Street1 string `json:"street1" validate:"omitempty,lte=100"`
	Street2 string `json:"street2" validate:"omitempty,lte=100"`
	City    string `json:"city" validate:"omitempty,lte=100"`
	State   string `json:"state" validate:"omitempty,lte=50"`
	ZipCode string `json:"zipCode" validate:"omitempty,lte=20"`
	Country string `json:"country" validate:"omitempty,lte=50"`
}
type Location struct {
	Id          string `json:"id" validate:"omitempty,uuid4"`
	Name        string `json:"name" validate:"omitempty,lte=100"`
	Description string `json:"description" validate:"omitempty,lte=200"`
	Phone       string `json:"phone" validate:"omitempty,lte=50"`
	Email       string `json:"email" validate:"omitempty,email"`
	Website     string `json:"website" validate:"omitempty,url"`
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
type Company struct {
	Location
	Contacts []struct {
		Person
	} `json:"contacts" validate:"dive"`
	ContactIds []string `json:"contactIds" validate:"dive,omitempty,uuid4"`
}
type CompanyNodes struct {
	CompanyNodes []Company `json:"companyNodes" validate:"dive"`
}

func (s *CompanyNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Store struct {
	Location
	Contacts []struct {
		Person
	} `json:"contacts"`
	ContactIds []string `json:"contactIds"`
}
type StoreNodes struct {
	StoreNodes []Store `json:"storeNodes" validate:"dive"`
}

func (s *StoreNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Warehouse struct {
	Location
}
type WarehouseNodes struct {
	WarehouseNodes []Warehouse `json:"warehouseNodes" validate:"dive"`
}

func (s *WarehouseNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Website struct {
	Location
	Domain string `json:"domain"`
}
type WebsiteNodes struct {
	WebsiteNodes []Website `json:"websiteNodes" validate:"dive"`
}

func (s *WebsiteNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
