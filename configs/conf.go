package configs

import (
	"encoding/json"
	"errors"
	//"log"
	"os"
	"path/filepath"
)

type Opener interface {
	Open(filePath string) (fileStream []byte, err error)
}
type Loader interface {
	Load(fileStream []byte) (err error)
}
type OpenerLoader interface {
	Opener
	Loader
}

func Load(conf OpenerLoader) (err error) {
	filePath := os.Getenv("CONFPATH")
	if filePath == "" {
		os.Setenv("CONFPATH", "/home/robertsmoto/dev/skustor/configs/conf.json")
		filePath = os.Getenv("CONFPATH")
	}
	filePath, _ = filepath.Abs(filePath)

	file, err := conf.Open(filePath)
	err = conf.Load(file)
	return err
}

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
		VanityUrl    string `json:"vanityUrl"`
	} `json:"doSpaces"`

	TempFileDir string `json:"tempFileDir"`
	UploadPrefix       string   `json:"uploadPrefix"`
    RootDir string `json:"rootDir"`
	Var02       int8   `json:"var02"`
}

func (c *Config) Open(filePath string) (fileStream []byte, err error) {

	fileExists := exists(filePath)
	if fileExists == false {
		err = errors.New("The config file does not exist.")
		return fileStream, err
	}

	fileStream, err = os.ReadFile(filePath)
	return fileStream, err
}

func (c *Config) Load(fileStream []byte) (err error) {

	isJson := isJson(fileStream)
	if isJson == false {
		err = errors.New("Config file isn't in a valid json format.")
		return err
	}
	json.Unmarshal(fileStream, &c)
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
