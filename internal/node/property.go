package node

import (
	"errors"
	"time"
)

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

func NewProperty(name string, value any, readOnly bool) *NodeProperty {
	return &NodeProperty{
		Type:        String,
		Name:        name,
		Description: "",
		Optional:    false,
		Key:         "",
		Value:       value,
		Options:     nil,
		ReadOnly:    readOnly,
	}
}

// isNumber checks if the value is any numeric type
func isNumber(v any) bool {
	switch v.(type) {
	case float64, float32:
		return true
	case int64, int32, int16, int8, int:
		return true
	case uint64, uint32, uint16, uint8, uint:
		return true
	default:
		return false
	}
}

// isWholeNumber checks if a float64 value represents a whole number
func isWholeNumber(f float64) bool {
	return f == float64(int64(f))
}

// isPositiveWholeNumber checks if a float64 value represents a positive whole number
func isPositiveWholeNumber(f float64) bool {
	return f >= 0 && f == float64(uint64(f))
}

func (p *NodeProperty) ValidateValue(v any) error {
	if v == nil {
		if p.Optional {
			return nil
		}
		return errors.New("value is required")
	}
	
	switch p.Type {
	case Bool:
		_, ok := v.(bool)
		if !ok {
			return errors.New("invalid boolean value")
		}
	case Int64:
		// Accept various number types that can be converted to int64
		switch val := v.(type) {
		case int64, int32, int16, int8, int:
			// Valid integer types
		case float64:
			// JSON unmarshals numbers as float64, check if it's a whole number
			if !isWholeNumber(val) {
				return errors.New("invalid int64 value: not a whole number")
			}
		default:
			return errors.New("invalid int64 value")
		}
	case Uint64:
		// Accept various number types that can be converted to uint64
		switch val := v.(type) {
		case uint64, uint32, uint16, uint8, uint:
			// Valid unsigned integer types
		case float64:
			// JSON unmarshals numbers as float64, check if it's a positive whole number
			if !isPositiveWholeNumber(val) {
				return errors.New("invalid uint64 value: must be positive whole number")
			}
		default:
			return errors.New("invalid uint64 value")
		}
	case Float64:
		if !isNumber(v) {
			return errors.New("invalid float64 value")
		}
	case String:
		_, ok := v.(string)
		if !ok {
			return errors.New("invalid string value")
		}
	case JSON:
		switch v.(type) {
		case map[string]any, []any:
			// Valid JSON types
		default:
			return errors.New("invalid json value")
		}
	case XML:
		_, ok := v.(string)
		if !ok {
			return errors.New("invalid xml value")
		}
	case Date:
		switch v.(type) {
		case time.Time:
			// Valid time type
		case string:
			// Allow string representation that can be parsed
			// This is common when dates come from JSON
		default:
			return errors.New("invalid date value")
		}
	case StringArray:
		_, ok := v.([]string)
		if !ok {
			// Check if it's []any with all strings (common from JSON)
			if arr, ok := v.([]any); ok {
				for _, item := range arr {
					if _, ok := item.(string); !ok {
						return errors.New("invalid string array value: all elements must be strings")
					}
				}
			} else {
				return errors.New("invalid string array value")
			}
		}
	case NumberArray:
		_, ok := v.([]float64)
		if !ok {
			// Check if it's []any with all numbers (common from JSON)
			if arr, ok := v.([]any); ok {
				for _, item := range arr {
					if !isNumber(item) {
						return errors.New("invalid number array value: all elements must be numbers")
					}
				}
			} else {
				return errors.New("invalid number array value")
			}
		}
	case BooleanArray:
		_, ok := v.([]bool)
		if !ok {
			// Check if it's []any with all booleans (common from JSON)
			if arr, ok := v.([]any); ok {
				for _, item := range arr {
					if _, ok := item.(bool); !ok {
						return errors.New("invalid boolean array value: all elements must be booleans")
					}
				}
			} else {
				return errors.New("invalid boolean array value")
			}
		}
	default:
		return errors.New("invalid property type")
	}
	return nil
}

