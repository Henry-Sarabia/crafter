package crafter

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"testing"
)

func TestNew(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(c)
}

func TestLoadRecipes(t *testing.T) {
	r, err := loadRecipes()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(r)
}

func TestRecipe(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/ring.json")
	if err != nil {
		t.Fatal(err)
	}

	var r recipe
	if err := json.Unmarshal(f, &r); err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
}

func TestPropType(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/wood.json")
	if err != nil {
		t.Fatal(err)
	}

	var p propType
	if err := json.Unmarshal(f, &p); err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}

func TestLoadGroups(t *testing.T) {
	g, err := loadGroups()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range g {
		fmt.Println(v)
	}
}

func TestLoadMaterials(t *testing.T) {
	m, err := loadMaterials()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range m {
		fmt.Println(v)
	}
}

func TestLoadTypes(t *testing.T) {
	tp, err := loadTypes()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range tp {
		fmt.Println(v)
	}
}