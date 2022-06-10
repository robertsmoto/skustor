package models

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

    "example.com/configs/conf"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/pborman/uuid"
)

func ConfigLoad() (SvConf conf.Config, err error) {
	SvConf = conf.Config{}
	err = SvConf.LoadJson("../internal/conf/test_data/test_config.json")
	return SvConf, err
}

func TestContent(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestContentLoad ...")
	testFile, err := os.ReadFile("./test_data/content.json")
	if err != nil {
		t.Errorf("TestContent %s", err)
		fmt.Println("error", err)
	}
	p := PageNodes{}
	_ = p.JsonLoad(testFile)
	//fmt.Println(p)

	validate := validator.New()
	err = validate.Struct(p)
	if err != nil {
		t.Errorf("TestPages %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}

func TestImage(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestImageLoad ...")
	testFile, err := os.ReadFile("./test_data/images.json")
	if err != nil {
		t.Errorf("TestImage %s", err)
		fmt.Println("error", err)
	}
	i := ImageNodes{}
	_ = i.JsonLoad(testFile)
	//fmt.Println(i)

	validate := validator.New()
	err = validate.Struct(i)
	if err != nil {
		t.Errorf("TestImages %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}
func TestPerson(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestPersonLoad ...")
	testFile, err := os.ReadFile("./test_data/contact.json")
	if err != nil {
		t.Errorf("TestPerson %s", err)
		fmt.Println("error", err)
	}
	c := ContactNodes{}
	_ = c.JsonLoad(testFile)
	//fmt.Println(c)

	validate := validator.New()
	err = validate.Struct(c)
	if err != nil {
		t.Errorf("TestCompanies %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}

func TestLocation(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestCompanyLoad ...")
	testFile, err := os.ReadFile("./test_data/companies.json")
	if err != nil {
		t.Errorf("TestLocations %s", err)
		fmt.Println("error", err)
	}
	c := CompanyNodes{}
	_ = c.JsonLoad(testFile)
	//fmt.Println(c)

	validate := validator.New()
	err = validate.Struct(c)
	if err != nil {
		t.Errorf("TestCompanies %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}

func TestItem(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestProductLoad ...")
	testFile, err := os.ReadFile("./test_data/products.json")
	if err != nil {
		t.Errorf("TestProducts %s", err)
		fmt.Println("error", err)
	}
	p := ProductNodes{}
	_ = p.JsonLoad(testFile)
	//fmt.Println(p)

	validate := validator.New()

	err = validate.Struct(p)
	if err != nil {
		t.Errorf("TestDepartments %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}

func TestBrand_JsonHandler(t *testing.T) {
    // read file (will eventually come from the request)
	testFile, err := os.ReadFile("./test_data/brands.json")
	if err != nil {
		t.Errorf("TestBrands %s", err)
		fmt.Println("error", err)
	}
    // need userId (will eventually come from the request)
    userId := uuid.Parse("d24f5e24-5368-417d-b730-f727577b8247")
    // make the db connection
	SvConf, err := ConfigLoad()
	if err != nil {
		fmt.Println("Error loading Config", err)
	}
	// open the db connection
	dbCred := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		SvConf.DbDevelopment.Host,
		SvConf.DbDevelopment.Port,
		SvConf.DbDevelopment.User,
		SvConf.DbDevelopment.Pass,
		SvConf.DbDevelopment.Dnam,
		SvConf.DbDevelopment.Sslm,
	)
	db, err := sql.Open("postgres", dbCred)

    brands := BrandNodes{}
    err = JsonHandler(&brands, testFile, db, userId)

    db.Close()
}
