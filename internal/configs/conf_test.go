package configs

import (
	//"encoding/json"
	"fmt"
	"testing"
    "os"
)



func Test_LoadConfigInterface(t *testing.T) {
	// tests a file exists, can be read and is in the correct json format

	config := Config{}
    Load(&config)
    
	// output struct
	fmt.Println(config)
	//check variables loaded correctly
	if os.Getenv("PGDNAM") != "skustor_development" {
		t.Errorf("Did not parse json variable correctly %s", os.Getenv("PGDNAM"))
	}
	if os.Getenv("VAR002") != "1" {
		t.Errorf("Did not parse json variable correctly %s", os.Getenv("VAR002"))
	}
}

func Test_LoadConfig(t *testing.T) {
	// tests a file exists, can be read and is in the correct json format

	config := Config{}
    fileStream, err := config.Open("./test_data/conf.json")
    err = config.Load(fileStream)
	if err != nil {
		t.Error(err)
	}
	// output struct
	fmt.Println(config)
	//check variables loaded correctly
	if os.Getenv("PGDNAM") != "skustor_development" {
		t.Errorf("Did not parse json variable correctly %s", os.Getenv("PGDNAM"))
	}
	if os.Getenv("VAR002") != "1" {
		t.Errorf("Did not parse json variable correctly %s", os.Getenv("VAR002"))
	}
}
