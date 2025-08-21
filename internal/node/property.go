package node

type PropertyType int

const (
	// Basic Types
	Invalid PropertyType = iota
	Bool
	Int64
	Uint64
	Float64
	String
	// Special Types
	JSON
	XML
	Date
	// Arrays
	StringArray
	NumberArray
	BooleanArray
)

// InputType defines how the property should be displayed in the UI
type InputType int

const (
	InputTypeText        InputType = iota + 1 // Single line text input
	InputTypeTextarea                         // Multi-line text input
	InputTypeNumber                           // Number input
	InputTypeRange                            // Slider input
	InputTypeSelect                           // Dropdown select
	InputTypeMultiSelect                      // Multiple selection
	InputTypeSwitch                           // Toggle switch
	InputTypeCheckbox                         // Checkbox
	InputTypeExpression                       // Expression input
	InputTypeKeyValue                         // Key-value pairs input
)

// InputConfig defines the configuration for input UI
type InputConfig struct {
	Type        InputType `json:"type"`
	Min         *float64  `json:"min,omitempty"`         // For range and number inputs
	Max         *float64  `json:"max,omitempty"`         // For range and number inputs
	Step        *float64  `json:"step,omitempty"`        // For range and number inputs
	Placeholder string    `json:"placeholder,omitempty"` // For text and textarea inputs
	Rows        int       `json:"rows,omitempty"`        // For textarea input
	Multiple    bool      `json:"multiple,omitempty"`    // For select inputs
}

func (p PropertyType) String() string {
	switch p {
	case Invalid:
		return "Invalid"
	case Bool:
		return "Bool"
	case Int64:
		return "Int64"
	case Uint64:
		return "Uint64"
	case Float64:
		return "Float64"
	case String:
		return "String"
	case JSON:
		return "JSON"
	case XML:
		return "XML"
	case Date:
		return "Date"
	case StringArray:
		return "StringArray"
	case NumberArray:
		return "NumberArray"
	case BooleanArray:
		return "BooleanArray"
	default:
		return "Unknown"
	}
}

type NodeProperty struct {
	// Type of the property (string, number, boolean, array, json, xml)
	Type PropertyType `json:"type"`
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
	Options []string `json:"options"`
	// not editable
	ReadOnly bool `json:"readOnly"`
}

// Helper functions for creating float64 pointers
func Float64Ptr(v float64) *float64 {
	return &v
}

func (p *NodeProperty) ValidateValue(v any) error {
	return ValidatePropertyValue(p.Type, v, p.Optional)
}
