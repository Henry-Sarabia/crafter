package crafter

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
)

type crafter struct {
	Recipes []recipe
	Types   []propType
	Groups  []propTypeGroup
}

// recipe is a set of information used to generate a unique item.
type recipe struct {
	Name       string      `json:"name"`
	BaseValue  float64     `json:"base_value"`
	BaseWeight float64     `json:"base_weight"`
	Comps      []component `json:"components"`
}

// component is an orthogonal section of an item. Each component for an item
// contains its own list of properties from which to choose from during item
// generation.
type component struct {
	Name     string     `json:"name"`
	Required bool       `json:"required"`
	Props    []property `json:"properties"`
}

// property is an orthogonal characteristic of a component. Each property
// describes some physical aspect of an item component. E.g. material, shape, engraving
type property struct {
	Name          string   `json:"name"`
	Required      bool     `json:"required"`
	TypeRefs      []string `json:"types"`
	TypeGroupRefs []string `json:"type_groups"`
	Types         []*propType
	TypeGroups    []*propTypeGroup
	// type groups? or concat to Types
}

// propType is a specific type of property. For example, the material property
// has several possible propTypes: base metal, precious metal, wood, gem, etc.
type propType struct {
	Name               string   `json:"name"`
	WeightFactor       factor   `json:"weight_factor"`
	MinorValueFactor   factor   `json:"minor_value_factor"`
	MinorValueVariants []string `json:"minor_value_variants"`
	AvgValueFactor     factor   `json:"avg_value_factor"`
	AvgValueDetails    []string `json:"avg_value_variants"`
	MajorValueFactor   factor   `json:"major_value_factor"`
	MajorValueDetails  []string `json:"major_value_variants"`
	Prefixes           []string `json:"prefix_references"`
}

// propTypeGroup are a group of propTypes. They're primarily used
// to include many types of propTypes that are similar to one another.
// For example, the propTypeGroup 'object' may contain every propType that
// can be represented as a physical form. The propTypeGroup 'weapon' may
// contain every type of weapon.
type propTypeGroup struct {
	Name     string   `json:"name"`
	TypeRefs []string `json:"types"`
	Types    []*propType
}

// propTypeVariant is a specific variant of a propertyType. For example, the
// material property type 'wood' has several possible propTypeVariants: ash,
// bamboo, oak, maple, etc.
type propTypeVariant struct {
	name  string
	value factor
}

// factor is a float used primarily for multiplying other values
// by a scalar.
type factor float64

func New() (*crafter, error) {
	c := &crafter{}

	r, err := loadRecipes()
	if err != nil {
		return nil, err
	}
	c.Recipes = append(c.Recipes, r...)

	g, err := loadGroups()
	if err != nil {
		return nil, err
	}
	c.Groups = append(c.Groups, g...)

	m, err := loadMaterials()
	if err != nil {
		return nil, err
	}
	c.Types = append(c.Types, m...)

	t, err := loadTypes()
	if err != nil {
		return nil, err
	}
	c.Types = append(c.Types, t...)

	err = c.linkGroups()
	if err != nil {
		return nil, err
	}

	err = c.linkRecipes()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *crafter) linkGroups() error {
	for i := range c.Groups {
		for _, ref := range c.Groups[i].TypeRefs {
			t, err := c.find(ref)
			if err != nil {
				return err
			}

			c.Groups[i].Types = append(c.Groups[i].Types, t)
		}
	}

	return nil
}

func (c *crafter) linkRecipes() error {
	for _, r := range c.Recipes {
		for _, comp := range r.Comps {
			for i := range comp.Props {
				for _, ref := range comp.Props[i].TypeRefs {
					t, err := c.find(ref)
					if err != nil {
						return err
					}

					comp.Props[i].Types = append(comp.Props[i].Types, t)
				}

				for _, ref := range comp.Props[i].TypeGroupRefs {
					t, err := c.findGroup(ref)
					if err != nil {
						return err
					}

					comp.Props[i].TypeGroups = append(comp.Props[i].TypeGroups, t)
				}
			}
		}
	}

	return nil
}

func (c *crafter) find(s string) (*propType, error) {
	for i := range c.Types {
		if c.Types[i].Name == s {
			fmt.Println("found type! ", c.Types[i].Name)
			return &c.Types[i], nil
		}
	}

	return nil, errors.Errorf("cannot find propType '%s' in memory", s)
}

func (c *crafter) findGroup(s string) (*propTypeGroup, error) {
	for i := range c.Groups {
		if c.Groups[i].Name == s {
			fmt.Println("found typeGroup!" , c.Groups[i].Name)
			return &c.Groups[i], nil
		}
	}

	return nil, errors.Errorf("cannot find propTypeGroup '%s' in memory", s)
}

func loadRecipes() ([]recipe, error) {
	names, err := filepath.Glob("testdata/recipes/*.json")
	if err != nil {
		return nil, err
	}

	var r []recipe

	for _, n := range names {
		f, err := ioutil.ReadFile(n)
		if err != nil {
			return nil, err
		}

		var dst recipe
		if err := json.Unmarshal(f, &dst); err != nil {
			return nil, err
		}

		r = append(r, dst)
	}

	return r, nil
}

// loadGroups reads the propTypeGroups JSON file and returns the data as a
// slice of propTypeGroups.
func loadGroups() ([]propTypeGroup, error) {
	f, err := ioutil.ReadFile("testdata/properties/typegroups.json")
	if err != nil {
		return nil, err
	}

	var dst []propTypeGroup
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

func loadMaterials() ([]propType, error) {
	f, err := ioutil.ReadFile("testdata/properties/materials.json")
	if err != nil {
		return nil, err
	}

	var dst []propType
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

func loadTypes() ([]propType, error) {
	f, err := ioutil.ReadFile("testdata/properties/types.json")
	if err != nil {
		return nil, err
	}

	var dst []propType
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}
