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

func Load(conf OpenerLoader) { 
	filePath := os.Getenv("CONFPATH")
	if filePath == "" {
		os.Setenv("CONFPATH", "/home/robertsmoto/dev/skustor/internal/configs/conf.json")
		filePath = os.Getenv("CONFPATH")
	}
	filePath, _ = filepath.Abs(filePath)

	file, err := conf.Open(filePath)
	err = conf.Load(file)
    if err != nil {
        panic("Configs.Load() error.")
    }
}

type Config struct {
	DbPostgres struct {
        Dnam string `json:"dnam"`
        Host string `json:"host"`
        Port string    `json:"port"`
        User string `json:"user"`
        Pass string `json:"pass"`
        Pkey string `json:"pkey"`
        Sslm string `json:"sslm"`
	} `json:"dbPostgres"`
	DoSpaces struct {
		UseSpaces    string   `json:"useSpaces"`
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
	Var02       string   `json:"var02"`
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
    // set env variables here
    os.Setenv("PGDNAM", c.DbPostgres.Dnam)
    os.Setenv("PGHOST", c.DbPostgres.Host)
    os.Setenv("PGPORT", string(c.DbPostgres.Port))
    os.Setenv("PGUSER", c.DbPostgres.User)
    os.Setenv("PGPASS", c.DbPostgres.Pass)
    os.Setenv("PGPKEY", c.DbPostgres.Pkey)
    os.Setenv("PGSSLM", c.DbPostgres.Sslm)
    os.Setenv("DOUSES", c.DoSpaces.UseSpaces)
    os.Setenv("DOAKEY", c.DoSpaces.AccessKey)
    os.Setenv("DOSECR", c.DoSpaces.Secret)
    os.Setenv("DOBCKT", c.DoSpaces.BucketName)
    os.Setenv("DOCDOM", c.DoSpaces.CustomDomain)
    os.Setenv("DOREGN", c.DoSpaces.RegionName)
    os.Setenv("DOENDU", c.DoSpaces.EndpointUrl)
    os.Setenv("DOVANU", c.DoSpaces.VanityUrl)
    os.Setenv("TMPDIR", c.TempFileDir)
    os.Setenv("ULOADP", c.UploadPrefix)
    os.Setenv("ROOTDR", c.RootDir)
    os.Setenv("VAR002", c.Var02)
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
