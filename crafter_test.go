package crafter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

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