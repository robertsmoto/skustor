package models

import "encoding/json"

type Content struct {
	WebsiteIds    []string `json:"websiteIds" validate:"dive,omitempty,uuid4"`
	AuthorIds     []string `json:"authorIds" validate:"dive,omitempty,uuid4"`
	CategoryIds   []string `json:"categoryIds" validate:"dive,omitempty,uuid4"`
	TagIds        []string `json:"tagIds" validate:"dive,omitempty,uuid4"`
	ImageIds      []string `json:"imageIds" validate:"dive,omitempty,uuid4"`
	ParentId      string   `json:"parentId" validate:"omitempty,uuid4"`
	Id            string   `json:"id" validate:"omitempty,uuid4"`
	PublishedTime string   `json:"publishedTime" validate:"omitempty,datetime=15:04 MST"`
	Published     string   `json:"published" validate:"omitempty,datetime=2006-01-02"`
	Modified      string   `json:"modified" validate:"omitempty,datetime=2006-01-02"`
	Keywords      string   `json:"keywords"`
	Title         string   `json:"title"`
	Excerpt       string   `json:"excerpt"`
	Body          string   `json:"body" validate:"omitempty"`
	Footer        string   `json:"footer"`
}
type Document struct {
	Content
}
type DocumentNodes struct {
	DocumentNodes []Document `json:"documentNodes" validate:"dive"`
}

func (s *DocumentNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Page struct {
	Content
}
type PageNodes struct {
	PageNodes []Page `json:"pageNodes" validate:"dive"`
}

func (s *PageNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}

type Post struct {
	Content
}
type PostNodes struct {
	PostNodes []Page `json:"postNodes" validate:"dive"`
}

func (s *PostNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
