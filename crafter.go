package crafter

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
)

const (
	recipeDirPath string = "testdata/recipes"
	groupDirPath  string = "testdata/properties/groups"
	typeDirPath   string = "testdata/properties/types"
)

type crafter struct {
	Recipes map[string]*recipe
	Types   map[string]*propType
	Groups  map[string]*propTypeGroup
}

// recipe is a set of information used to generate a unique item.
type recipe struct {
	Name       string       `json:"name"`
	BaseValue  float64      `json:"base_value"`
	BaseWeight float64      `json:"base_weight"`
	Comps      []*component `json:"components"`
}

// component is an orthogonal section of an item. Each component for an item
// contains its own list of properties from which to choose from during item
// generation.
type component struct {
	Name     string      `json:"name"`
	Required bool        `json:"required"`
	Props    []*property `json:"properties"`
}

// property is an orthogonal characteristic of a component. Each property
// describes some physical aspect of an item component. E.g. material, shape, engraving
type property struct {
	Name          string   `json:"name"`
	Required      bool     `json:"required"`
	TypeRefs      []string `json:"type_refs"`
	TypeGroupRefs []string `json:"type_group_refs"`
	Types         []*propType
	TypeGroups    []*propTypeGroup
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
	TypeRefs []string `json:"type_refs"`
	Types    []*propType
}

// factor is a float used primarily for multiplying other values
// by a scalar.
type factor float64

func New() (*crafter, error) {
	c := &crafter{}

	r, err := loadRecipeDir(recipeDirPath)
	if err != nil {
		return nil, err
	}
	c.Recipes = r

	g, err := loadGroupDir(groupDirPath)
	if err != nil {
		return nil, err
	}
	c.Groups = g

	t, err := loadTypeDir(typeDirPath)
	if err != nil {
		return nil, err
	}
	c.Types = t

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

// linkGroups iterates through every group's TypeRefs and adds the
// corresponding propType addresses to the group's Types slice.
func (c *crafter) linkGroups() error {
	for _, g := range c.Groups {
		for _, ref := range g.TypeRefs {
			t, ok := c.Types[ref]
			if !ok {
				return errors.Errorf("cannot find propType '%s' in c.Types", ref)
			}

			g.Types = append(g.Types, t)
		}
	}

	return nil
}

// linkRecipes links every recipe's typeRefs and typeGroupRefs to their
// respective propTypes and propTypeGroups.
func (c *crafter) linkRecipes() error {
	for _, rec := range c.Recipes {
		for _, comp := range rec.Comps {
			for _, prop := range comp.Props {
				for _, ref := range prop.TypeRefs {
					t, ok := c.Types[ref]
					if !ok {
						return errors.Errorf("cannot find propType '%s' in c.Types", ref)
					}

					prop.Types = append(prop.Types, t)
				}

				for _, ref := range prop.TypeGroupRefs {
					g, ok := c.Groups[ref]
					if !ok {
						return errors.Errorf("cannot find propTypeGroup '%s' in c.Groups", ref)
					}

					prop.TypeGroups = append(prop.TypeGroups, g)
				}
			}
		}
	}

	return nil
}

// loadRecipeDir loads the recipes from the given directory. The directory path
// should not include a trailing slash. Each JSON file should contain an array of
// recipes.
func loadRecipeDir(path string) (map[string]*recipe, error) {
	names, err := filepath.Glob(path + "/*.json")
	if err != nil {
		return nil, err
	}

	m := make(map[string]*recipe)

	for _, n := range names {
		r, err := unmarshalRecipes(n)
		if err != nil {
			return nil, err
		}

		for _, v := range r {
			// TODO: add duplicate checking?
			m[v.Name] = &v
		}
	}

	return m, nil
}

// unmarshalRecipes returns the recipes from the given JSON file.
func unmarshalRecipes(filename string) ([]recipe, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dst []recipe
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

// loadGroupDir loads the propTypeGroups from the given directory. The
// directory path should not include a trailing slash. Each JSON file should
// contain an array of propTypeGroups.
func loadGroupDir(path string) (map[string]*propTypeGroup, error) {
	names, err := filepath.Glob(path + "/*.json")
	if err != nil {
		return nil, err
	}

	m := make(map[string]*propTypeGroup)

	for _, n := range names {
		g, err := unmarshalGroups(n)
		if err != nil {
			return nil, err
		}

		for _, v := range g {
			// TODO: add duplicate checking?
			m[v.Name] = &v
		}
	}

	return m, nil
}

// unmarshalGroups returns the propTypeGroups from the given JSON file.
func unmarshalGroups(filename string) ([]propTypeGroup, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dst []propTypeGroup
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

// loadTypeDir loads the propTypes from the given directory. The directory path
// should not include a trailing slash. Each JSON file should contain an array
// of propTypes.
func loadTypeDir(path string) (map[string]*propType, error) {
	names, err := filepath.Glob(path + "/*.json")
	if err != nil {
		return nil, err
	}

	m := make(map[string]*propType)

	for _, n := range names {
		t, err := unmarshalTypes(n)
		if err != nil {
			return nil, err
		}

		for _, v := range t {
			// TODO: add duplicate checking?
			m[v.Name] = &v
		}
	}

	return m, nil
}

// unmarshalTypes returns the propTypes from the given JSON file.
func unmarshalTypes(filename string) ([]propType, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dst []propType
	if err := json.Unmarshal(f, &dst); err != nil {
		return nil, err
	}

	return dst, nil
}
