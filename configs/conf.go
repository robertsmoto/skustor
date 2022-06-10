package configs

import (
	"encoding/json"
	"errors"
	"os"
)

type db struct {
	Dnam string `json:"dnam"`
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Pkey string `json:"pkey"`
	Sslm string `json:"sslm"`
}

type Config struct {
	DbDevelopment struct {
		db 
	} `json:"dbDevelopment"`

	DbStaging struct {
		db
	} `json:"dbStaging"`

	DbProduction struct {
		db
	} `json:"dbProduction"`
	DoSpaces struct {
		UseSpaces    bool   `json:"useSpaces"`
		AccessKey    string `json:"accessKey"`
		Secret       string `json:"secret"`
		BucketName   string `json:"bucketName"`
		CustomDomain string `json:"customDomain"`
		RegionName   string `json:"regionName"`
		EndpointUrl  string `json:"endpointUrl"`
	} `json:"doSpaces"`

	TempFileDir string `json:"tempFileDir"`
	Var01       int8   `json:"var01"`
	Var02       int8   `json:"var02"`
}

func (c *Config) LoadJson(filePath string) (err error) {
	// will read from the delete or post request

	fileExists := exists(filePath)
	if fileExists == false {
		err = errors.New("The config file does not exist.")
		return err
	}

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	isJson := isJson(configFile)
	if isJson == false {
		err = errors.New("The config file is not in a valid json format.")
		return err
	}

	json.Unmarshal(configFile, &c)

	return err

}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func isJson(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
