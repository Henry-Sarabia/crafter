package crafter

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestNew(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(c)
}

func TestLoadRecipeDir(t *testing.T) {
	r, err := loadRecipeDir(recipeDirPath)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(r)
}

func TestLoadGroupDir(t *testing.T) {
	g, err := loadGroupDir(groupDirPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range g {
		fmt.Println(v)
	}
}

func TestLoadTypeDir(t *testing.T) {
	tp, err := loadTypeDir(typeDirPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range tp {
		fmt.Println(v)
	}
}