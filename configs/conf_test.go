package configs

import (
	//"encoding/json"
	"fmt"
	"testing"
	//"os"
)

func TestLoadConfig(t *testing.T) {
	// tests a file exists, can be read and is in the correct json format

	SvConf := configs.Config{}
	err := SvConf.LoadJson("./test_data/test_config.json")
	if err != nil {
		t.Error(err)
	}
	// output struct
	fmt.Println(SvConf)
	//check variables loaded correctly
	if SvConf.DbDevelopment.Dnam != "skustor_development" {
		t.Errorf("Did not parse json variable correctly %s", SvConf.DbDevelopment.Dnam)
	}
	if SvConf.Var01 != 0 {
		t.Errorf("Did not parse json variable correctly %d", SvConf.Var01)
	}
}
