package node

import (
	"reflect"
	"testing"
)

func TestStringProperty(t *testing.T) {
	prop := StringProperty("test", "Test String", "A test string property", "default", true)
	
	if prop.Type != String {
		t.Errorf("Type = %v, want %v", prop.Type, String)
	}
	if prop.Key != "test" {
		t.Errorf("Key = %v, want %v", prop.Key, "test")
	}
	if prop.Value != "default" {
		t.Errorf("Value = %v, want %v", prop.Value, "default")
	}
	if prop.Input.Type != InputTypeText {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeText)
	}
}

func TestTextAreaProperty(t *testing.T) {
	prop := TextAreaProperty("desc", "Description", "Long description", "multi\nline", 5, false)
	
	if prop.Type != String {
		t.Errorf("Type = %v, want %v", prop.Type, String)
	}
	if prop.Input.Type != InputTypeTextarea {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeTextarea)
	}
	if prop.Input.Rows != 5 {
		t.Errorf("Input.Rows = %v, want %v", prop.Input.Rows, 5)
	}
	if prop.Optional != false {
		t.Errorf("Optional = %v, want %v", prop.Optional, false)
	}
}

func TestIntProperty(t *testing.T) {
	prop := IntProperty("count", "Count", "Item count", 42, false)
	
	if prop.Type != Int64 {
		t.Errorf("Type = %v, want %v", prop.Type, Int64)
	}
	if prop.Value != int64(42) {
		t.Errorf("Value = %v, want %v", prop.Value, int64(42))
	}
	if prop.Input.Type != InputTypeNumber {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeNumber)
	}
}

func TestIntRangeProperty(t *testing.T) {
	prop := IntRangeProperty("volume", "Volume", "Volume level", 50, 0, 100, true)
	
	if prop.Type != Int64 {
		t.Errorf("Type = %v, want %v", prop.Type, Int64)
	}
	if prop.Input.Type != InputTypeRange {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeRange)
	}
	if *prop.Input.Min != 0 {
		t.Errorf("Input.Min = %v, want %v", *prop.Input.Min, 0)
	}
	if *prop.Input.Max != 100 {
		t.Errorf("Input.Max = %v, want %v", *prop.Input.Max, 100)
	}
	if *prop.Input.Step != 1 {
		t.Errorf("Input.Step = %v, want %v", *prop.Input.Step, 1)
	}
}

func TestFloatProperty(t *testing.T) {
	prop := FloatProperty("price", "Price", "Product price", 19.99, false)
	
	if prop.Type != Float64 {
		t.Errorf("Type = %v, want %v", prop.Type, Float64)
	}
	if prop.Value != 19.99 {
		t.Errorf("Value = %v, want %v", prop.Value, 19.99)
	}
	if *prop.Input.Step != 0.01 {
		t.Errorf("Input.Step = %v, want %v", *prop.Input.Step, 0.01)
	}
}

func TestBoolProperty(t *testing.T) {
	prop := BoolProperty("enabled", "Enabled", "Enable feature", true, false)
	
	if prop.Type != Bool {
		t.Errorf("Type = %v, want %v", prop.Type, Bool)
	}
	if prop.Value != true {
		t.Errorf("Value = %v, want %v", prop.Value, true)
	}
	if prop.Input.Type != InputTypeSwitch {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeSwitch)
	}
}

func TestSelectProperty(t *testing.T) {
	options := []string{"option1", "option2", "option3"}
	prop := SelectProperty("choice", "Choice", "Select an option", options, "option1", false)
	
	if prop.Type != String {
		t.Errorf("Type = %v, want %v", prop.Type, String)
	}
	if !reflect.DeepEqual(prop.Options, options) {
		t.Errorf("Options = %v, want %v", prop.Options, options)
	}
	if prop.Input.Type != InputTypeSelect {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeSelect)
	}
}

func TestMultiSelectProperty(t *testing.T) {
	options := []string{"tag1", "tag2", "tag3"}
	defaultValue := []string{"tag1", "tag3"}
	prop := MultiSelectProperty("tags", "Tags", "Select tags", options, defaultValue, true)
	
	if prop.Type != StringArray {
		t.Errorf("Type = %v, want %v", prop.Type, StringArray)
	}
	if !reflect.DeepEqual(prop.Options, options) {
		t.Errorf("Options = %v, want %v", prop.Options, options)
	}
	if prop.Input.Type != InputTypeMultiSelect {
		t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeMultiSelect)
	}
	if !prop.Input.Multiple {
		t.Errorf("Input.Multiple = %v, want %v", prop.Input.Multiple, true)
	}
}

func TestJSONProperty(t *testing.T) {
	defaultValue := map[string]any{"key": "value"}
	prop := JSONProperty("config", "Config", "Configuration", defaultValue, true)
	
	if prop.Type != JSON {
		t.Errorf("Type = %v, want %v", prop.Type, JSON)
	}
	if !reflect.DeepEqual(prop.Value, defaultValue) {
		t.Errorf("Value = %v, want %v", prop.Value, defaultValue)
	}
	if prop.Input.Rows != 10 {
		t.Errorf("Input.Rows = %v, want %v", prop.Input.Rows, 10)
	}
}

func TestArrayProperties(t *testing.T) {
	t.Run("StringArray", func(t *testing.T) {
		defaultValue := []string{"a", "b", "c"}
		prop := StringArrayProperty("items", "Items", "Item list", defaultValue, false)
		
		if prop.Type != StringArray {
			t.Errorf("Type = %v, want %v", prop.Type, StringArray)
		}
		if !reflect.DeepEqual(prop.Value, defaultValue) {
			t.Errorf("Value = %v, want %v", prop.Value, defaultValue)
		}
	})
	
	t.Run("NumberArray", func(t *testing.T) {
		defaultValue := []float64{1.0, 2.0, 3.0}
		prop := NumberArrayProperty("values", "Values", "Number values", defaultValue, true)
		
		if prop.Type != NumberArray {
			t.Errorf("Type = %v, want %v", prop.Type, NumberArray)
		}
		if !reflect.DeepEqual(prop.Value, defaultValue) {
			t.Errorf("Value = %v, want %v", prop.Value, defaultValue)
		}
	})
	
	t.Run("BooleanArray", func(t *testing.T) {
		defaultValue := []bool{true, false, true}
		prop := BooleanArrayProperty("flags", "Flags", "Flag values", defaultValue, false)
		
		if prop.Type != BooleanArray {
			t.Errorf("Type = %v, want %v", prop.Type, BooleanArray)
		}
		if !reflect.DeepEqual(prop.Value, defaultValue) {
			t.Errorf("Value = %v, want %v", prop.Value, defaultValue)
		}
	})
}

func TestReadOnlyProperty(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		wantType PropertyType
	}{
		{"string", "test", String},
		{"bool", true, Bool},
		{"int", 42, Int64},
		{"float", 3.14, Float64},
		{"string array", []string{"a", "b"}, StringArray},
		{"number array", []float64{1, 2}, NumberArray},
		{"bool array", []bool{true, false}, BooleanArray},
		{"json map", map[string]any{"key": "value"}, JSON},
		{"json array", []any{1, 2, 3}, JSON},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prop := ReadOnlyProperty("key", "Name", "Description", tt.value)
			
			if prop.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", prop.Type, tt.wantType)
			}
			if !prop.ReadOnly {
				t.Errorf("ReadOnly = %v, want %v", prop.ReadOnly, true)
			}
			if !prop.Optional {
				t.Errorf("Optional = %v, want %v", prop.Optional, true)
			}
		})
	}
}

func TestOutputProperty(t *testing.T) {
	prop := OutputProperty("result", "Result", "Operation result", String)
	
	if prop.Type != String {
		t.Errorf("Type = %v, want %v", prop.Type, String)
	}
	if !prop.ReadOnly {
		t.Errorf("ReadOnly = %v, want %v", prop.ReadOnly, true)
	}
	if !prop.Optional {
		t.Errorf("Optional = %v, want %v", prop.Optional, true)
	}
}