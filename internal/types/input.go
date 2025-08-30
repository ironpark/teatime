package types

import "fmt"

// InputType represents the type of UI input component for a property.
// Each input type corresponds to a specific UI widget with unique behavior.
type InputType int

const (
	InputTypeText        InputType = iota + 1 // Single line text input
	InputTypeTextarea                         // Multi-line text input
	InputTypeExpression                       // Expression input
	InputTypeDate                             // Date input (time.Time, int64, string)
	InputTypeNumber                           // Number input
	InputTypeRange                            // Slider input
	InputTypeSelect                           // Dropdown select
	InputTypeMultiSelect                      // Multiple selection
	InputTypeSwitch                           // Toggle switch (only for boolean properties)
	InputTypeList                             // List input (only for array properties)
	InputTypeJson                             // JSON input (only for json properties)
	InputTypeKeyValue                         // Key-value pairs input (only for map properties)
)

var (
	// inputTypeNames maps InputType constants to their string representations.
	// Used internally for String() method implementation.
	inputTypeNames = []string{
		"text",
		"textarea",
		"expression",
		"date",
		"number",
		"range",
		"select",
		"multi-select",
		"switch",
		"list",
		"json",
		"kv",
	}

	// inputTypeMap provides reverse lookup from string to InputType.
	// Used by InputTypeFromString for parsing input type strings.
	inputTypeMap = map[string]InputType{
		"text":         InputTypeText,
		"textarea":     InputTypeTextarea,
		"expression":   InputTypeExpression,
		"date":         InputTypeDate,
		"number":       InputTypeNumber,
		"range":        InputTypeRange,
		"select":       InputTypeSelect,
		"multi-select": InputTypeMultiSelect,
		"switch":       InputTypeSwitch,
		"list":         InputTypeList,
		"json":         InputTypeJson,
		"kv":           InputTypeKeyValue,
	}
)

// InputTypeFromString converts a string to InputType.
// Returns an error if the string doesn't match any known input type.
//
// Supported input type strings:
//   - "text", "textarea", "expression", "date"
//   - "number", "range", "select", "multi-select"
//   - "switch", "list", "json", "kv"
func InputTypeFromString(s string) (InputType, error) {
	if _, ok := inputTypeMap[s]; !ok {
		return 0, fmt.Errorf("unknown input type: %s", s)
	}
	return inputTypeMap[s], nil
}

// String returns the string representation of the InputType.
func (i InputType) String() string {
	if i < 1 || int(i) > len(inputTypeNames) {
		return "unknown"
	}
	return inputTypeNames[i-1]
}

// InputConfig holds configuration options for UI input components.
// Different input types support different configuration options.
type InputConfig struct {
	Type        InputType `json:"type"`                  // Required: type of input component
	Min         *float64  `json:"min,omitempty"`         // For range and number inputs: minimum value
	Max         *float64  `json:"max,omitempty"`         // For range and number inputs: maximum value
	Step        *float64  `json:"step,omitempty"`        // For range and number inputs: increment step
	Placeholder string    `json:"placeholder,omitempty"` // For text and textarea inputs: placeholder text
	Multiple    bool      `json:"multiple,omitempty"`    // For select inputs: allow multiple selection
	Unique      bool      `json:"unique,omitempty"`      // For list inputs: ensure unique elements
}

// rangeCheck validates range-specific configuration options.
// For InputTypeRange, min/max/step are required and must be logical.
// For InputTypeNumber, min/max/step are optional.
// For other input types, these options are not allowed.
func (i InputConfig) rangeCheck() error {
	if i.Type != InputTypeRange && i.Type != InputTypeNumber {
		if i.Min != nil || i.Max != nil || i.Step != nil {
			return fmt.Errorf("min, max and step are not allowed for %s input", i.Type.String())
		}
		return nil
	}

	// For range input, all parameters are required
	if i.Type == InputTypeRange {
		if i.Min == nil || i.Max == nil || i.Step == nil {
			return fmt.Errorf("min, max and step are required for range input")
		}
		min, max, step := *i.Min, *i.Max, *i.Step
		if max < min {
			return fmt.Errorf("max must be greater than min")
		}
		if step > max-min {
			return fmt.Errorf("step must be less than max-min")
		}
	} else {
		// For number input, parameters are optional but must be valid if provided
		if i.Min != nil && i.Max != nil {
			min, max := *i.Min, *i.Max
			if max < min {
				return fmt.Errorf("max must be greater than min")
			}
			if i.Step != nil {
				step := *i.Step
				if step > max-min {
					return fmt.Errorf("step must be less than max-min")
				}
			}
		}
	}
	return nil
}

// Validate checks if the InputConfig is valid for its input type.
// Returns an error if any configuration option is invalid or incompatible
// with the specified input type.
//
// Validation rules:
//   - Type is required and must be a valid InputType
//   - Min/Max/Step are only allowed for range inputs
//   - Placeholder is not allowed for switch and key-value inputs
//   - For range inputs, min/max/step are required and must be logical
func (i InputConfig) Validate() error {
	if i.Type == 0 {
		return fmt.Errorf("type is required")
	}

	if i.Type < 1 || int(i.Type) > len(inputTypeNames) {
		return fmt.Errorf("invalid input type: %d", i.Type)
	}

	if err := i.rangeCheck(); err != nil {
		return err
	}
	if i.Placeholder != "" {
		if i.Type == InputTypeSwitch || i.Type == InputTypeKeyValue {
			return fmt.Errorf("placeholder is not allowed for %s input", i.Type.String())
		}
	}
	return nil
}

// ValidateType checks if the InputType is compatible with the given PropertyType.
// Returns an error if the input type is not suitable for the property type.
//
// Compatibility rules:
//   - InputTypeSwitch: only for Bool properties
//   - InputTypeList: only for Array properties
//   - InputTypeJson: only for JSON properties
//   - InputTypeKeyValue: only for Map properties
//   - InputTypeDate: only for Date properties
//   - InputTypeNumber/Range: for Int64, Float64 properties
//   - InputTypeText/Textarea: for String properties
//   - InputTypeExpression: for any property type (flexible input)
//   - InputTypeSelect: for Bool, Int64, Float64, String, Date properties (with enums)
//   - InputTypeList/MultiSelect: for Array properties
func (i InputConfig) ValidateType(t PropertyType) error {
	// expression input is compatible with any property type
	if i.Type == InputTypeExpression {
		return nil
	}

	switch i.Type {
	case InputTypeSwitch:
		if t != Bool {
			return fmt.Errorf("switch input is only compatible with boolean properties, got %s", t.String())
		}
	case InputTypeJson:
		if t != JSON {
			return fmt.Errorf("json input is only compatible with json properties, got %s", t.String())
		}
	case InputTypeKeyValue:
		if t != Map {
			return fmt.Errorf("key-value input is only compatible with map properties, got %s", t.String())
		}
	case InputTypeDate:
		if t != Date {
			return fmt.Errorf("date input is only compatible with date properties, got %s", t.String())
		}
	case InputTypeNumber, InputTypeRange:
		if t != Int64 && t != Float64 {
			return fmt.Errorf("%s input is only compatible with numeric properties (int64, float64), got %s", i.Type.String(), t.String())
		}
	case InputTypeText, InputTypeTextarea:
		if t != String {
			return fmt.Errorf("%s input is only compatible with string properties, got %s", i.Type.String(), t.String())
		}
	case InputTypeSelect:
		if t == Array || t == Map || t == JSON {
			return fmt.Errorf("%s input is not compatible with array, map or json properties, got %s", i.Type.String(), t.String())
		}
	case InputTypeList, InputTypeMultiSelect:
		if t != Array {
			return fmt.Errorf("%s input is only compatible with array properties, got %s", i.Type.String(), t.String())
		}
	default:
		return fmt.Errorf("unknown input type: %s", i.Type.String())
	}
	return nil
}
