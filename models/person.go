package models

import "encoding/json"

type Person struct {
	Id         string `json:"id" validate:"omitempty,uuid4"`
	Salutation string `json:"salutation" validate:"omitempty,lte=20"`
	Firstname  string `json:"firstname" validate:"omitempty,lte=100"`
	Lastname   string `json:"lastname" validate:"omitempty,lte=100"`
	Nickname   string `json:"nickname" validate:"omitempty,lte=100"`
	Phone      string `json:"phone" validate:"omitempty,lte=50"`
	Mobile     string `json:"mobile" validate:"omitempty,lte=50"`
	Email      string `json:"email" validate:"omitempty,email"`
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

// people
type Customer struct {
	Person
}
type CustomerNodes struct {
	CustomerNodes []Customer `json:"customerNodes" validate:"dive"`
}

func (s *CustomerNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Contact struct {
	Person
}
type ContactNodes struct {
	ContactNodes []Contact `json:"contactNodes" validate:"dive"`
}

func (s *ContactNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
