package crafter

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
	Name       string   `json:"name"`
	Required   bool     `json:"required"`
	Types      []string `json:"types"`
	TypeGroups []string `json:"type_groups"`
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
	Name  string   `json:"name"`
	Types []string `json:"types"`
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
