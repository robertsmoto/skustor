package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
    "path/filepath"

	//_ "github.com/lib/pq"
	"github.com/go-playground/validator/v10"
	"github.com/pborman/uuid"
	"github.com/robertsmoto/skustor/configs"
	"github.com/robertsmoto/skustor/images"
)

type Group struct {
	ImageNodes
	Id       string `json:"id" validate:"required,uuid4"`
	ParentId string `json:"parentId" validate:"omitempty,uuid4"`
	// Type in the database eg. rawMaterialTag, partTag, productTag, productBrand
	// will need join table m2m relationship
	Type           string   `json:"type" validate:"omitempty,lte=200"`
	Name           string   `json:"name" validate:"omitempty,lte=200"`
	Description    string   `json:"description" validate:"omitempty,lte=200"`
	Keywords       string   `json:"keywords" validate:"omitempty,lte=200"`
	LinkUrl        string   `json:"linkUrl" validate:"omitempty,url,lte=200"`
	LinkText       string   `json:"linkText" validate:"omitempty,lte=200"`
	RawMaterialIds []string `json:"rawMaterialIds" validate:"dive,omitempty,uuid4"`
	PartIds        []string `json:"partIds" validate:"dive,omitempty,uuid4"`
	ProductIds     []string `json:"productIds" validate:"dive,omitempty,uuid4"`
}
type Brand struct {
	Group
}
type BrandNodes struct {
	Nodes []Brand `json:"brandNodes" validate:"dive"`
}

func (s *BrandNodes) Load(fileBuffer *[]byte) (err error) {
	json.Unmarshal(*fileBuffer, &s)
	return err
}
func (s *BrandNodes) Validate() (err error) {
	validate := validator.New()
	err = validate.Struct(s)
	return err
}
func (s *BrandNodes) Upsert(db *sql.DB, userId string) (err error) {
	// check if struct is empty
	if s.Nodes == nil {
		fmt.Println("BrandNodes struct == nil")
		return err
	}
	// construct the sql upsert statement
	qstr := `
        INSERT INTO groups (
            id, user_id, parent_id, type, name, description, keywords,
            link_url, link_text
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE
        SET user_id = $2,
            parent_id = $3,
            type = $4,
            name = $5,
            description = $6,
            keywords = $7,
            link_url = $8,
            link_text = $9
        WHERE groups.id = $1;`
	// execute it
	for _, node := range s.Nodes {
		_, err = db.Exec(
			qstr, uuid.Parse(node.Id), uuid.Parse(userId), uuid.Parse(node.ParentId),
			"brand", node.Name, node.Description, node.Keywords, node.LinkUrl,
			node.LinkText)
	}
	return err
}

func (s *BrandNodes) RelatedTableUpsert(db *sql.DB, userId string) (err error) {

	for _, node := range s.Nodes {
		// check m2m relationships for each node
		// retrieve userId from the request

		if node.RawMaterialIds != nil {
			qstrSpecs := []string{"join_group_item", userId, node.Id,
				"user_id", "group_id", "item_id"}
			err = JoinTableUpsert(db, qstrSpecs, node.RawMaterialIds)
			if err != nil {
				fmt.Println(err)
			}
		}
		if node.PartIds != nil {
			qstrSpecs := []string{"join_group_item", userId, node.Id,
				"user_id", "group_id", "item_id"}
			err = JoinTableUpsert(db, qstrSpecs, node.PartIds)
			if err != nil {
				fmt.Println(err)
			}
		}
		if node.ProductIds != nil {
			qstrSpecs := []string{"join_group_item", userId, node.Id,
				"user_id", "group_id", "item_id"}
			err = JoinTableUpsert(db, qstrSpecs, node.ProductIds)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return err
}

func (s *BrandNodes) ImageSize(db *sql.DB, date, userId, userDir string) (err error) {
	// Process image nodes, resize images.
	for _, node := range s.Nodes {
		if node.ImageNodes.Nodes == nil {
			log.Print("##no image node")
			continue
		}

		config := configs.Config{}
		err := configs.Load(&config)
		if err != nil {
			log.Print("Error loading config.", err)
		}

		resizeImg := images.WebImage{
			UserDir:      userDir,
			Date:         date,
			TempFileDir:  config.TempFileDir,
			UploadPrefix: "media",
			DoBucket:       config.DoSpaces.BucketName,
			DoEndpointUrl:  config.DoSpaces.EndpointUrl,
			DoAccessKey:    config.DoSpaces.AccessKey,
			DoSecret:       config.DoSpaces.Secret,
			DoRegionName:   config.DoSpaces.RegionName,
			DoCacheControl: "max-age=604800",
			DoContentType:  "image/webp",
		}
		// special cases for resizeImage struct
		for _, imgNode := range node.ImageNodes.Nodes {
			log.Print("##node --> ", imgNode)
			if imgNode.Id == "" || imgNode.Url == "" {
				log.Print("Image must have id and url to process.")
				continue
			}
			// resize and upload new images
			// yeilds NewSizes map[string][]string
			//localFilePath[uploadFilePath, height, width]
			resizeImg.Url = imgNode.Url
			images.Resize(&resizeImg)

			// now add the resizeImg.NewSizes to the images table
			qstr := `
                INSERT INTO images (
                    id, user_id, group_id, url, height, width, title, alt,
                    caption, position, featured 
                )
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                ON CONFLICT (id) DO UPDATE
                SET user_id = $2,
                    group_id = $3,
                    url = $4,
                    height = $5,
                    width = $6,
                    title = $7,
                    alt = $8,
                    caption = $9,
                    position = $10,
                    featured = $11
                WHERE images.id = $1;`
			// execute qstr
			log.Print("##qstr images --> ", qstr)
			for _, i := range resizeImg.NewSizes {
				_, err = db.Exec(
					qstr, uuid.Parse(imgNode.Id), uuid.Parse(userId),
					uuid.Parse(node.Id), filepath.Join(config.DoSpaces.VanityUrl, i[0]),
                    i[1], i[2], imgNode.Title, imgNode.Alt, imgNode.Caption,
                    imgNode.Position, imgNode.Featured,
				)
			}
		}
	}
	return err
}

func (s *BrandNodes) Delete(db sql.DB) (err error) {
	log.Print("Not implemented.")

	//// check if struct is empty
	//if s.Nodes == nil {
	//fmt.Println("BrandNodes struct == nil")
	//return err
	//}
	//// construct the sql upsert statement

	//// execute it

	return err
}

type Category struct {
	Group
}
type CategoryNodes struct {
	CategoryNodes []Category `json:"categoryNodes" validate:"dive"`
}

func (s *CategoryNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Department struct {
	Group
}
type DepartmentNodes struct {
	DepartmentNodes []Department `json:"departmentNodes" validate:"dive"`
}

func (s *DepartmentNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Tag struct {
	Group
}
type TagNodes struct {
	TagNodes []Tag `json:"tagNodes" validate:"dive"`
}

func (s *TagNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
