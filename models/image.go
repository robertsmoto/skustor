package models

import "encoding/json"

type Image struct {
	Id       string `json:"id" validate:"omitempty,uuid4"`
	Url      string `json:"url" validate:"omitempty,url"`
	Title    string `json:"title" validate:"omitempty,lte=200"`
	Alt      string `json:"alt" validate:"omitempty,lte=100"`
	Caption  string `json:"caption" validate:"omitempty,lte=200"`
	Order    uint8  `json:"order" validate:"omitempty,number"`
	Featured uint8  `json:"featured" validate:"omitempty,number"`
}

type ImageNodes struct {
	ImageNodes []Image `json:"imageNodes" validate:"dive"`
}

func (s *ImageNodes) JsonLoad(fileBuffer []byte) (err error) {
	json.Unmarshal(fileBuffer, &s)
	return err
}
