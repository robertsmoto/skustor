package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SvAddress struct {
	Id         string `json:"id" validate:"omitempty,uuid4"`
	SvUserId   string `json:"userId" validate:"omitempty,uuid4"`
	PlaceId string `json:"placeId" validate:"omitempty,uuid4"`
	Type       string `json:"type" validate:"omitempty,lte=100,oneof=billing main mailing shipping"`
	Street1    string `json:"street1" validate:"omitempty,lte=100"`
	Street2    string `json:"street2" validate:"omitempty,lte=100"`
	City       string `json:"city" validate:"omitempty,lte=100"`
	State      string `json:"state" validate:"omitempty,lte=50"`
	ZipCode    string `json:"zipCode" validate:"omitempty,lte=20"`
	Country    string `json:"country" validate:"omitempty,lte=50"`
}


func (s *SvAddress) Process(userId string) (err error) {
	fmt.Println("SvAddress.Process() Not iplemented.")
    if err != nil {
		return fmt.Errorf("SvAddress.Process() %s", err)
    }
    return nil
}

func (s *SvAddress) Upsert(db *sql.DB) (err error) {
	if s == (&SvAddress{}) {
        return nil
	}
	qstr := `
        INSERT INTO address (
            id, sv_user_id, type, street1, street2, city,
            state, zipcode, country
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id=$2,
            type = $3,
            street1 = $4,
            street2 = $5,
            city = $6,
            state = $7,
            zipcode = $8,
            country = $9
        WHERE address.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Type, s.Street1, s.Street2,
		s.City, s.State, s.ZipCode, s.Country,
	)
	if err != nil {
		return fmt.Errorf("SvAddress.Upsert() %s", err)
	}
    return nil
}

func (s *SvAddress) ForeignKeyUpdate(db *sql.DB) (err error) {
	if s == (&SvAddress{}) {
		return nil
	}
	qstr := `
        UPDATE address
        SET place_id = $2
        WHERE address.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.PlaceId),
	)
	if err != nil {
		return fmt.Errorf("Address.ForeignKeyUpdate() %s", err)
	}
    return nil
}

func (s *SvAddress) RelatedTableUpsert(db *sql.DB) (err error) {
	fmt.Println("Address.RelatedTableUpsert() Not implemented.")
    if err != nil {
        return fmt.Errorf("Address.RelatedTableUpsert() %s", err)
    }
    return nil
}

func (s *SvAddress) Delete(db *sql.DB) (err error) {
	fmt.Println("Address.Delete() Not implemented.")
    if err != nil {
        return fmt.Errorf("Address.Delete() %s", err)
    }
    return nil
}

type Address struct {
    SvAddress `json:"address"`
}

func (s *Address) Load(fileBuffer *[]byte) (err error) {
    err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Address.Load() %s", err)
	}
    return nil
}

func (s *Address) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Address.Validate() %s", err)
	}
    return nil
}

type Addresses struct {
	Nodes []SvAddress `json:"addresses" validate:"dive"`
}

func (s *Addresses) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Addresses.Load() %s", err)
	}
    return nil
}

func (s *Addresses) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Addresses.Validate() %s", err)
	}
	return nil
}

func (s *Addresses) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("Addresses.Process() %s", err)
        }
    }
    return nil
}

func (s *Addresses) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Addresses.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Addresses) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Addresses.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Addresses) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Addresses.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Addresses) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Units.Delete() %s", err)
        }
    }
    return nil
}

type SvPlace struct {
	Id          string   `json:"id" validate:"omitempty,uuid4"`
	SvUserId    string   `json:"svUserId" validate:"omitempty,uuid4"`
	Type        string   `json:"type" valdidate:"omitempty,lte=100,oneof=company store warehouse website"`
	Name        string   `json:"name" validate:"omitempty,lte=100"`
	Description string   `json:"description" validate:"omitempty,lte=200"`
	Phone       string   `json:"phone" validate:"omitempty,lte=20"`
	Email       string   `json:"email" validate:"omitempty,lte=100,email"`
	Website     string   `json:"website" validate:"omitempty,lte=100,url"`
	Domain      string   `json:"domain" validate:"omitempty,lte=100"`
	AddressIds  []string `json:"addressIds" validate:"dive"`
    Address
	Addresses
}

func (s *SvPlace) Process(userId string) (err error) {
	fmt.Println("Place.Process() Not iplemented.")
    if err != nil {
        return fmt.Errorf("Place.Process() %s", err)
    }
    return nil
}

func (s *SvPlace) Upsert(db *sql.DB) (err error) {
	if s == (&SvPlace{}) {
        return nil
	}
	qstr := `
        INSERT INTO place (
            id, sv_user_id, type, name, description, phone,
            email, website, domain
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET sv_user_id=$2,
            type = $3,
            name = $4,
            description = $5,
            phone = $6,
            email = $7,
            website = $8,
            domain = $9
        WHERE place.id = $1;`
	_, err = db.Exec(
		qstr, FormatUUID(s.Id), FormatUUID(s.SvUserId), s.Type, s.Name, s.Description,
		s.Phone, s.Email, s.Website, s.Domain,
	)
	if err != nil {
		return fmt.Errorf("Place.Upsert() %s", err)
	}
    return nil
}

func (s *SvPlace) ForeignKeyUpdate(db *sql.DB) (err error) {
	fmt.Println("Loation.ForeignKeyUpdate() Not implemented.")
	if err != nil {
		return fmt.Errorf("Place.ForeignKeyUpdate() %s", err)
	}
    return nil
}

func (s *SvPlace) RelatedTableUpsert(db *sql.DB) (err error) {
	fmt.Println("Loation.RelatedTableUpsert() Not implemented.")
	if err != nil {
		return fmt.Errorf("Place.RelatedTableUpsert() %s", err)
	}
    return nil
}

func (s *SvPlace) Delete(db *sql.DB) (err error) {
	fmt.Println("Place.Delete() Not implemented.")
	if err != nil {
		return fmt.Errorf("Place.Delete() %s", err)
	}
    return nil
}

type Place struct {
    SvPlace `json:"place"`
}

func (s *Place) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Place.Load() %s", err)
	}
    return nil
}

func (s *Place) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Place.Validate() %s", err)
	}
    return nil
}

type Places struct {
	Nodes []SvPlace `json:"places" validate:"dive"`
}

func (s *Places) Load(fileBuffer *[]byte) (err error) {
	err = json.Unmarshal(*fileBuffer, &s)
	if err != nil {
		return fmt.Errorf("Places.Load() %s", err)
	}
    return nil
}

func (s *Places) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("Places.Validate() %s", err)
	}
    return nil
}

func (s *Places) Process(userId string) (err error) {
    for _, node := range s.Nodes {
        err = node.Process(userId)
        if err != nil {
            return fmt.Errorf("Places.Process() %s", err)
        }
    }
    return nil
}

func (s *Places) Upsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Upsert(db)
        if err != nil {
            return fmt.Errorf("Places.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Places) ForeignKeyUpdate(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.ForeignKeyUpdate(db)
        if err != nil {
            return fmt.Errorf("Places.Upsert() %s", err)
        }
    }
    return nil
}

func (s *Places) RelatedTableUpsert(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.RelatedTableUpsert(db)
        if err != nil {
            return fmt.Errorf("Places.RelatedTableUpsert() %s", err)
        }
    }
    return nil
}

func (s *Places) Delete(db *sql.DB) (err error) {
    for _, node := range s.Nodes {
        err = node.Delete(db)
        if err != nil {
            return fmt.Errorf("Places.Delete() %s", err)
        }
    }
    return nil
}
