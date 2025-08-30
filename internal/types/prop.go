package types

type PropertyType int

const (
	Bool PropertyType = iota + 1
	Int64
	Float64
	String

	Date
	Array
	Map
	JSON
)

var propertyTypeNames = []string{
	"bool",
	"int64",
	"float64",
	"string",
	"date",
	"array",
	"map",
	"json",
}

func (p PropertyType) FallbackInputType() InputType {
	switch p {
	case Bool:
		return InputTypeSwitch
	case Int64, Float64:
		return InputTypeNumber
	case String:
		return InputTypeText
	case Date:
		return InputTypeDate
	case Array:
		return InputTypeList
	case Map:
		return InputTypeKeyValue
	case JSON:
		return InputTypeJson
	}
	return InputTypeText
}

func (p PropertyType) String() string {
	if p < 1 || int(p) > len(propertyTypeNames) {
		return "unknown"
	}
	return propertyTypeNames[p-1]
}

type Property struct {
	// Type of the property (string, number, boolean, array, json, xml)
	Type PropertyType `json:"type"`
	// Type of the value of the property (bool, int64, float64, string, json)
	ValType PropertyType `json:"valType,omitempty"`
	// Input configuration for UI
	Input *InputConfig `json:"input,omitempty"`
	// Name of the property
	Name string `json:"name"`
	// Description of the property
	Description string `json:"description"`
	// Whether the property is required
	Optional bool `json:"optional"`
	// Key of the property
	Key string `json:"key"`
	// Value of the property
	Value any `json:"value"`
	// Options of the property (for select/multi-select)
	// enums must be a slice of bool, int64, float64, string
	// enum values must be unique
	Enums []any `json:"enums"`
	// hide on node preview
	HideOnPreview bool `json:"hideOnPreview"`
}
