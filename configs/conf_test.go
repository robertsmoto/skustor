package configs

import (
	//"encoding/json"
	"fmt"
	"testing"
	//"os"
)

func TestLoadConfig(t *testing.T) {
	// tests a file exists, can be read and is in the correct json format

	config := Config{}
	err := ConfigOpenLoad(&config, "./conf.json")
	if err != nil {
		t.Error(err)
	}
	// output struct
	fmt.Println(config)
	//check variables loaded correctly
	if config.DbDevelopment.Dnam != "skustor_development" {
		t.Errorf("Did not parse json variable correctly %s", config.DbDevelopment.Dnam)
	}
	if config.Var01 != 0 {
		t.Errorf("Did not parse json variable correctly %d", config.Var01)
	}
}
