package types

import (
	"testing"
)

func TestInputTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected InputType
		wantErr  bool
	}{
		{"text input", "text", InputTypeText, false},
		{"textarea input", "textarea", InputTypeTextarea, false},
		{"expression input", "expression", InputTypeExpression, false},
		{"date input", "date", InputTypeDate, false},
		{"number input", "number", InputTypeNumber, false},
		{"range input", "range", InputTypeRange, false},
		{"select input", "select", InputTypeSelect, false},
		{"multi-select input", "multi-select", InputTypeMultiSelect, false},
		{"switch input", "switch", InputTypeSwitch, false},
		{"list input", "list", InputTypeList, false},
		{"json input", "json", InputTypeJson, false},
		{"key-value input", "kv", InputTypeKeyValue, false},
		{"unknown input", "unknown", 0, true},
		{"empty input", "", 0, true},
		{"invalid input", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := InputTypeFromString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("InputTypeFromString(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("InputTypeFromString(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("InputTypeFromString(%q) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestInputType_String(t *testing.T) {
	tests := []struct {
		name     string
		input    InputType
		expected string
	}{
		{"text input", InputTypeText, "text"},
		{"textarea input", InputTypeTextarea, "textarea"},
		{"expression input", InputTypeExpression, "expression"},
		{"date input", InputTypeDate, "date"},
		{"number input", InputTypeNumber, "number"},
		{"range input", InputTypeRange, "range"},
		{"select input", InputTypeSelect, "select"},
		{"multi-select input", InputTypeMultiSelect, "multi-select"},
		{"switch input", InputTypeSwitch, "switch"},
		{"list input", InputTypeList, "list"},
		{"json input", InputTypeJson, "json"},
		{"key-value input", InputTypeKeyValue, "kv"},
		{"invalid input type 0", InputType(0), "unknown"},
		{"invalid input type -1", InputType(-1), "unknown"},
		{"invalid input type 999", InputType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("InputType(%d).String() = %q, expected %q", int(tt.input), result, tt.expected)
			}
		})
	}
}

func TestInputConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  InputConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:   "valid text input",
			config: InputConfig{Type: InputTypeText, Placeholder: "Enter text"},
			wantErr: false,
		},
		{
			name:   "valid textarea input",
			config: InputConfig{Type: InputTypeTextarea, Placeholder: "Enter description"},
			wantErr: false,
		},
		{
			name:   "valid number input",
			config: InputConfig{Type: InputTypeNumber},
			wantErr: false,
		},
		{
			name: "valid range input",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Max:  floatPtr(100),
				Step: floatPtr(1),
			},
			wantErr: false,
		},
		{
			name:   "valid select input",
			config: InputConfig{Type: InputTypeSelect, Multiple: true},
			wantErr: false,
		},
		{
			name:   "valid switch input",
			config: InputConfig{Type: InputTypeSwitch},
			wantErr: false,
		},
		{
			name:   "valid list input",
			config: InputConfig{Type: InputTypeList, Unique: true},
			wantErr: false,
		},
		{
			name:   "valid json input",
			config: InputConfig{Type: InputTypeJson},
			wantErr: false,
		},
		{
			name:   "valid key-value input",
			config: InputConfig{Type: InputTypeKeyValue},
			wantErr: false,
		},
		{
			name:    "invalid input type 0",
			config:  InputConfig{Type: InputType(0)},
			wantErr: true,
			errMsg:  "type is required",
		},
		{
			name:    "invalid input type -1",
			config:  InputConfig{Type: InputType(-1)},
			wantErr: true,
			errMsg:  "invalid input type",
		},
		{
			name:    "invalid input type 999",
			config:  InputConfig{Type: InputType(999)},
			wantErr: true,
			errMsg:  "invalid input type",
		},
		{
			name: "range input missing min",
			config: InputConfig{
				Type: InputTypeRange,
				Max:  floatPtr(100),
				Step: floatPtr(1),
			},
			wantErr: true,
			errMsg:  "min, max and step are required for range input",
		},
		{
			name: "range input missing max",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Step: floatPtr(1),
			},
			wantErr: true,
			errMsg:  "min, max and step are required for range input",
		},
		{
			name: "range input missing step",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Max:  floatPtr(100),
			},
			wantErr: true,
			errMsg:  "min, max and step are required for range input",
		},
		{
			name: "range input max < min",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(100),
				Max:  floatPtr(0),
				Step: floatPtr(1),
			},
			wantErr: true,
			errMsg:  "max must be greater than min",
		},
		{
			name: "range input step too large",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Max:  floatPtr(10),
				Step: floatPtr(20),
			},
			wantErr: true,
			errMsg:  "step must be less than max-min",
		},
		{
			name: "text input with range options",
			config: InputConfig{
				Type: InputTypeText,
				Min:  floatPtr(0),
				Max:  floatPtr(100),
				Step: floatPtr(1),
			},
			wantErr: true,
			errMsg:  "min, max and step are not allowed for text input",
		},
		{
			name:    "switch input with placeholder",
			config:  InputConfig{Type: InputTypeSwitch, Placeholder: "Toggle me"},
			wantErr: true,
			errMsg:  "placeholder is not allowed for switch input",
		},
		{
			name:    "key-value input with placeholder",
			config:  InputConfig{Type: InputTypeKeyValue, Placeholder: "Enter key-value"},
			wantErr: true,
			errMsg:  "placeholder is not allowed for kv input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("InputConfig.Validate() expected error, got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					// For partial error message matching
					found := false
					if len(tt.errMsg) < len(err.Error()) {
						found = err.Error()[:len(tt.errMsg)] == tt.errMsg ||
								err.Error()[len(err.Error())-len(tt.errMsg):] == tt.errMsg
					}
					for i := 0; i <= len(err.Error())-len(tt.errMsg); i++ {
						if err.Error()[i:i+len(tt.errMsg)] == tt.errMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("InputConfig.Validate() error = %q, expected to contain %q", err.Error(), tt.errMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("InputConfig.Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestInputConfig_rangeCheck(t *testing.T) {
	tests := []struct {
		name    string
		config  InputConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "non-range input with no range options",
			config:  InputConfig{Type: InputTypeText},
			wantErr: false,
		},
		{
			name: "non-range input with range options",
			config: InputConfig{
				Type: InputTypeText,
				Min:  floatPtr(0),
				Max:  floatPtr(100),
			},
			wantErr: true,
			errMsg:  "min, max and step are not allowed for text input",
		},
		{
			name: "range input with all options",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Max:  floatPtr(100),
				Step: floatPtr(1),
			},
			wantErr: false,
		},
		{
			name: "range input missing options",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
			},
			wantErr: true,
			errMsg:  "min, max and step are required for range input",
		},
		{
			name: "range input with invalid range",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(50),
				Max:  floatPtr(10),
				Step: floatPtr(1),
			},
			wantErr: true,
			errMsg:  "max must be greater than min",
		},
		{
			name: "range input with step too large",
			config: InputConfig{
				Type: InputTypeRange,
				Min:  floatPtr(0),
				Max:  floatPtr(5),
				Step: floatPtr(10),
			},
			wantErr: true,
			errMsg:  "step must be less than max-min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.rangeCheck()
			if tt.wantErr {
				if err == nil {
					t.Errorf("InputConfig.rangeCheck() expected error, got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("InputConfig.rangeCheck() error = %q, expected %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("InputConfig.rangeCheck() unexpected error: %v", err)
				}
			}
		})
	}
}

// Helper function to create float64 pointers
func floatPtr(f float64) *float64 {
	return &f
}

func TestInputType_Coverage(t *testing.T) {
	// Test that all InputType constants have corresponding string representations
	allInputTypes := []InputType{
		InputTypeText,
		InputTypeTextarea,
		InputTypeExpression,
		InputTypeDate,
		InputTypeNumber,
		InputTypeRange,
		InputTypeSelect,
		InputTypeMultiSelect,
		InputTypeSwitch,
		InputTypeList,
		InputTypeJson,
		InputTypeKeyValue,
	}

	for _, inputType := range allInputTypes {
		t.Run(inputType.String(), func(t *testing.T) {
			// Test String() method
			str := inputType.String()
			if str == "unknown" {
				t.Errorf("InputType(%d).String() returned 'unknown', expected valid string", int(inputType))
			}

			// Test round-trip conversion
			parsed, err := InputTypeFromString(str)
			if err != nil {
				t.Errorf("InputTypeFromString(%q) error: %v", str, err)
			}
			if parsed != inputType {
				t.Errorf("Round-trip conversion failed: %v -> %q -> %v", inputType, str, parsed)
			}
		})
	}

	// Test that inputTypeNames and inputTypeMap are consistent
	if len(inputTypeNames) != len(inputTypeMap) {
		t.Errorf("inputTypeNames length (%d) != inputTypeMap length (%d)", len(inputTypeNames), len(inputTypeMap))
	}

	for i, name := range inputTypeNames {
		expectedType := InputType(i + 1)
		if mapType, exists := inputTypeMap[name]; !exists {
			t.Errorf("inputTypeMap missing entry for %q", name)
		} else if mapType != expectedType {
			t.Errorf("inputTypeMap[%q] = %v, expected %v", name, mapType, expectedType)
		}
	}
}

func TestInputConfig_ValidateType(t *testing.T) {
	tests := []struct {
		name        string
		inputType   InputType
		propertyType PropertyType
		wantErr     bool
		errMsg      string
	}{
		// Valid combinations
		{
			name:         "switch with bool",
			inputType:    InputTypeSwitch,
			propertyType: Bool,
			wantErr:      false,
		},
		{
			name:         "list with array",
			inputType:    InputTypeList,
			propertyType: Array,
			wantErr:      false,
		},
		{
			name:         "json with json",
			inputType:    InputTypeJson,
			propertyType: JSON,
			wantErr:      false,
		},
		{
			name:         "key-value with map",
			inputType:    InputTypeKeyValue,
			propertyType: Map,
			wantErr:      false,
		},
		{
			name:         "date with date",
			inputType:    InputTypeDate,
			propertyType: Date,
			wantErr:      false,
		},
		{
			name:         "number with int64",
			inputType:    InputTypeNumber,
			propertyType: Int64,
			wantErr:      false,
		},
		{
			name:         "number with float64",
			inputType:    InputTypeNumber,
			propertyType: Float64,
			wantErr:      false,
		},
		{
			name:         "range with int64",
			inputType:    InputTypeRange,
			propertyType: Int64,
			wantErr:      false,
		},
		{
			name:         "range with float64",
			inputType:    InputTypeRange,
			propertyType: Float64,
			wantErr:      false,
		},
		{
			name:         "text with string",
			inputType:    InputTypeText,
			propertyType: String,
			wantErr:      false,
		},
		{
			name:         "textarea with string",
			inputType:    InputTypeTextarea,
			propertyType: String,
			wantErr:      false,
		},
		{
			name:         "expression with string",
			inputType:    InputTypeExpression,
			propertyType: String,
			wantErr:      false,
		},
		{
			name:         "expression with any type",
			inputType:    InputTypeExpression,
			propertyType: Int64,
			wantErr:      false,
		},
		{
			name:         "select with string",
			inputType:    InputTypeSelect,
			propertyType: String,
			wantErr:      false,
		},
		{
			name:         "select with int64",
			inputType:    InputTypeSelect,
			propertyType: Int64,
			wantErr:      false,
		},
		{
			name:         "multi-select with array",
			inputType:    InputTypeMultiSelect,
			propertyType: Array,
			wantErr:      false,
		},
		// Invalid combinations
		{
			name:         "switch with string",
			inputType:    InputTypeSwitch,
			propertyType: String,
			wantErr:      true,
			errMsg:       "switch input is only compatible with boolean properties",
		},
		{
			name:         "select with array",
			inputType:    InputTypeSelect,
			propertyType: Array,
			wantErr:      true,
			errMsg:       "select input is not compatible with array, map or json properties",
		},
		{
			name:         "select with map",
			inputType:    InputTypeSelect,
			propertyType: Map,
			wantErr:      true,
			errMsg:       "select input is not compatible with array, map or json properties",
		},
		{
			name:         "select with json",
			inputType:    InputTypeSelect,
			propertyType: JSON,
			wantErr:      true,
			errMsg:       "select input is not compatible with array, map or json properties",
		},
		{
			name:         "list with string",
			inputType:    InputTypeList,
			propertyType: String,
			wantErr:      true,
			errMsg:       "list input is only compatible with array properties",
		},
		{
			name:         "multi-select with string",
			inputType:    InputTypeMultiSelect,
			propertyType: String,
			wantErr:      true,
			errMsg:       "multi-select input is only compatible with array properties",
		},
		{
			name:         "json with string",
			inputType:    InputTypeJson,
			propertyType: String,
			wantErr:      true,
			errMsg:       "json input is only compatible with json properties",
		},
		{
			name:         "key-value with string",
			inputType:    InputTypeKeyValue,
			propertyType: String,
			wantErr:      true,
			errMsg:       "key-value input is only compatible with map properties",
		},
		{
			name:         "date with string",
			inputType:    InputTypeDate,
			propertyType: String,
			wantErr:      true,
			errMsg:       "date input is only compatible with date properties",
		},
		{
			name:         "number with string",
			inputType:    InputTypeNumber,
			propertyType: String,
			wantErr:      true,
			errMsg:       "number input is only compatible with numeric properties",
		},
		{
			name:         "range with bool",
			inputType:    InputTypeRange,
			propertyType: Bool,
			wantErr:      true,
			errMsg:       "range input is only compatible with numeric properties",
		},
		{
			name:         "text with bool",
			inputType:    InputTypeText,
			propertyType: Bool,
			wantErr:      true,
			errMsg:       "text input is only compatible with string properties",
		},
		{
			name:         "textarea with int64",
			inputType:    InputTypeTextarea,
			propertyType: Int64,
			wantErr:      true,
			errMsg:       "textarea input is only compatible with string properties",
		},
		{
			name:         "unknown input type",
			inputType:    InputType(999),
			propertyType: String,
			wantErr:      true,
			errMsg:       "unknown input type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := InputConfig{Type: tt.inputType}
			err := config.ValidateType(tt.propertyType)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateType() expected error, got nil")
				} else if tt.errMsg != "" {
					// Check if error message contains expected substring
					found := false
					if len(tt.errMsg) <= len(err.Error()) {
						for i := 0; i <= len(err.Error())-len(tt.errMsg); i++ {
							if err.Error()[i:i+len(tt.errMsg)] == tt.errMsg {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("ValidateType() error = %q, expected to contain %q", err.Error(), tt.errMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateType() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestInputConfig_ValidateType_EdgeCases(t *testing.T) {
	// Test expression input with all PropertyType combinations
	allPropertyTypes := []PropertyType{
		Bool, Int64, Float64, String, Date, Array, Map, JSON,
	}
	
	for _, propType := range allPropertyTypes {
		t.Run("expression_with_"+propType.String(), func(t *testing.T) {
			config := InputConfig{Type: InputTypeExpression}
			err := config.ValidateType(propType)
			if err != nil {
				t.Errorf("ValidateType() with expression input should accept all property types, got error: %v", err)
			}
		})
	}

	// Test select input with valid property types
	validSelectTypes := []PropertyType{Bool, Int64, Float64, String, Date}
	for _, propType := range validSelectTypes {
		t.Run("select_with_"+propType.String(), func(t *testing.T) {
			config := InputConfig{Type: InputTypeSelect}
			err := config.ValidateType(propType)
			if err != nil {
				t.Errorf("ValidateType() with select input should accept %s property type, got error: %v", propType.String(), err)
			}
		})
	}

	// Test select input with invalid property types
	invalidSelectTypes := []PropertyType{Array, Map, JSON}
	for _, propType := range invalidSelectTypes {
		t.Run("select_invalid_with_"+propType.String(), func(t *testing.T) {
			config := InputConfig{Type: InputTypeSelect}
			err := config.ValidateType(propType)
			if err == nil {
				t.Errorf("ValidateType() with select input should reject %s property type", propType.String())
			}
		})
	}
}