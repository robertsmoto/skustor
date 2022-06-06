package conf

import (
	//"encoding/json"
	"fmt"
	"testing"
    //"os"
)

func TestLoadConfig(t *testing.T) {
	// tests a file exists, can be read and is in the correct json format

    SvConf := Config{}
    err := SvConf.LoadJson("./test_data/config.json")
	if err != nil {
		t.Error(err)
	}
    // output struct
	fmt.Println(SvConf)
	//check variables loaded correctly
	if SvConf.DbSodaVault.Dnam != "sodavault" {
		t.Errorf("Did not parse json variable correctly %s", SvConf.DbSodaVault.Dnam)
	}
	if SvConf.DbSodaVault.Pass != "svpassword" {
		t.Errorf("Did not parse json variable correctly %s", SvConf.DbSodaVault.Pass)
	}
	if SvConf.DoSpaces.Secret != "secret" {
		t.Errorf("Did not parse json variable correctly %s", SvConf.DoSpaces.Secret)
	}
	if SvConf.Var01 != 0 {
		t.Errorf("Did not parse json variable correctly %d", SvConf.Var01)
	}
}
