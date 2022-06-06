package loader

import (
	//"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestAttributes(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestAttributesLoad ...")
	testFile, err := os.ReadFile("./test_data/attributes.json")
	if err != nil {
		t.Errorf("TestAttributes %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)

}
func TestBrands(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestDepartmentsLoad ...")
	testFile, err := os.ReadFile("./test_data/brands.json")
	if err != nil {
		t.Errorf("TestDepartments %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	fmt.Println(SvData)
	fmt.Println("##Here -->", SvData.Brand.Id)
	validate := validator.New()

	err = validate.Struct(SvData)
	if err != nil {
		t.Errorf("TestDepartments %s", err)
		fmt.Println(err.(validator.ValidationErrors))
	}
}
func TestCategories(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestCategoriesLoad ...")
	testFile, err := os.ReadFile("./test_data/categories.json")
	if err != nil {
		t.Errorf("TestCategorys %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestCompanies(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestCompaniesLoad ...")
	testFile, err := os.ReadFile("./test_data/companies.json")
	if err != nil {
		t.Errorf("TestCategorys %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestDepartments(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestDepartmentsLoad ...")
	testFile, err := os.ReadFile("./test_data/departments.json")
	if err != nil {
		t.Errorf("TestDepartments %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestMenus(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestMenusLoad ...")
	testFile, err := os.ReadFile("./test_data/menus.json")
	if err != nil {
		t.Errorf("TestMenus %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestParts(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestPartsLoad ...")
	testFile, err := os.ReadFile("./test_data/parts.json")
	if err != nil {
		t.Errorf("TestParts %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
	//fmt.Println(SvData.Part.Item)
	//for _, p := range SvData.Parts {
	//fmt.Println(p.Item)
	//}
}
func TestPriceClasses(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestPriceClassesLoad ...")
	testFile, err := os.ReadFile("./test_data/priceClasses.json")
	if err != nil {
		t.Errorf("TestPriceClasses %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestProducts(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestProductsLoad ...")
	testFile, err := os.ReadFile("./test_data/products.json")
	if err != nil {
		t.Errorf("TestProducts %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestRawMaterials(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestRawMaterialsLoad ...")
	testFile, err := os.ReadFile("./test_data/rawMaterials.json")
	if err != nil {
		t.Errorf("TestRawMaterials %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
func TestStores(t *testing.T) {
	// tests load function with parser interface
	fmt.Println("TestRawStoresLoad ...")
	testFile, err := os.ReadFile("./test_data/stores.json")
	if err != nil {
		t.Errorf("TestRawStores %s", err)
		fmt.Println("error", err)
	}
	SvData := SvData{}
	_ = SvData.LoadJson(testFile)
	//fmt.Println(SvData)
}
