package node

import (
	"reflect"
	"testing"
)

func TestPropertyBuilder(t *testing.T) {
	t.Run("String property with options", func(t *testing.T) {
		prop := StringProp("name", "Name",
			WithDescription("User's name"),
			WithDefault("John Doe"),
			WithPlaceholder("Enter your name"),
			Optional(),
		)
		
		if prop.Type != String {
			t.Errorf("Type = %v, want %v", prop.Type, String)
		}
		if prop.Description != "User's name" {
			t.Errorf("Description = %v, want %v", prop.Description, "User's name")
		}
		if prop.Value != "John Doe" {
			t.Errorf("Value = %v, want %v", prop.Value, "John Doe")
		}
		if !prop.Optional {
			t.Errorf("Optional = %v, want %v", prop.Optional, true)
		}
		if prop.Input.Placeholder != "Enter your name" {
			t.Errorf("Placeholder = %v, want %v", prop.Input.Placeholder, "Enter your name")
		}
	})
	
	t.Run("Text property with rows", func(t *testing.T) {
		prop := TextProp("description", "Description",
			WithDefault("Default text"),
			WithRows(10),
		)
		
		if prop.Input.Type != InputTypeTextarea {
			t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeTextarea)
		}
		if prop.Input.Rows != 10 {
			t.Errorf("Rows = %v, want %v", prop.Input.Rows, 10)
		}
	})
	
	t.Run("Int with range", func(t *testing.T) {
		prop := IntProp("volume", "Volume",
			WithDefault(int64(50)),
			WithRange(0, 100, 1),
			WithInput(InputTypeRange),
		)
		
		if prop.Type != Int64 {
			t.Errorf("Type = %v, want %v", prop.Type, Int64)
		}
		if *prop.Input.Min != 0 {
			t.Errorf("Min = %v, want %v", *prop.Input.Min, 0)
		}
		if *prop.Input.Max != 100 {
			t.Errorf("Max = %v, want %v", *prop.Input.Max, 100)
		}
		if prop.Input.Type != InputTypeRange {
			t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeRange)
		}
	})
	
	t.Run("Float with step", func(t *testing.T) {
		prop := FloatProp("price", "Price",
			WithDefault(19.99),
			WithStep(0.1),
		)
		
		if prop.Type != Float64 {
			t.Errorf("Type = %v, want %v", prop.Type, Float64)
		}
		if *prop.Input.Step != 0.1 {
			t.Errorf("Step = %v, want %v", *prop.Input.Step, 0.1)
		}
	})
	
	t.Run("Bool property", func(t *testing.T) {
		prop := BoolProp("enabled", "Enabled",
			WithDefault(true),
			WithDescription("Enable this feature"),
		)
		
		if prop.Type != Bool {
			t.Errorf("Type = %v, want %v", prop.Type, Bool)
		}
		if prop.Value != true {
			t.Errorf("Value = %v, want %v", prop.Value, true)
		}
		if prop.Input.Type != InputTypeSwitch {
			t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeSwitch)
		}
	})
	
	t.Run("Select property", func(t *testing.T) {
		options := []string{"small", "medium", "large"}
		prop := SelectProp("size", "Size", options,
			WithDefault("medium"),
		)
		
		if prop.Type != String {
			t.Errorf("Type = %v, want %v", prop.Type, String)
		}
		if !reflect.DeepEqual(prop.Options, options) {
			t.Errorf("Options = %v, want %v", prop.Options, options)
		}
		if prop.Input.Type != InputTypeSelect {
			t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeSelect)
		}
	})
	
	t.Run("MultiSelect property", func(t *testing.T) {
		options := []string{"tag1", "tag2", "tag3"}
		prop := MultiSelectProp("tags", "Tags", options,
			WithDefault([]string{"tag1"}),
			Optional(),
		)
		
		if prop.Type != StringArray {
			t.Errorf("Type = %v, want %v", prop.Type, StringArray)
		}
		if !prop.Input.Multiple {
			t.Errorf("Multiple = %v, want %v", prop.Input.Multiple, true)
		}
		if !prop.Optional {
			t.Errorf("Optional = %v, want %v", prop.Optional, true)
		}
	})
	
	t.Run("Output property", func(t *testing.T) {
		prop := OutputProp(String, "result", "Result",
			WithDescription("Operation result"),
		)
		
		if !prop.ReadOnly {
			t.Errorf("ReadOnly = %v, want %v", prop.ReadOnly, true)
		}
		if !prop.Optional {
			t.Errorf("Optional = %v, want %v", prop.Optional, true)
		}
	})
	
	t.Run("Required property", func(t *testing.T) {
		prop := StringProp("required", "Required Field",
			Required(), // Explicitly set as required
		)
		
		if prop.Optional {
			t.Errorf("Optional = %v, want %v", prop.Optional, false)
		}
	})
	
	t.Run("Chaining multiple options", func(t *testing.T) {
		prop := IntProp("complex", "Complex Property",
			WithDescription("A complex property"),
			WithDefault(int64(42)),
			WithMin(1),
			WithMax(100),
			WithStep(5),
			Optional(),
			WithInput(InputTypeRange),
		)
		
		if prop.Description != "A complex property" {
			t.Errorf("Description = %v, want %v", prop.Description, "A complex property")
		}
		if prop.Value != int64(42) {
			t.Errorf("Value = %v, want %v", prop.Value, int64(42))
		}
		if *prop.Input.Min != 1 {
			t.Errorf("Min = %v, want %v", *prop.Input.Min, 1)
		}
		if *prop.Input.Max != 100 {
			t.Errorf("Max = %v, want %v", *prop.Input.Max, 100)
		}
		if *prop.Input.Step != 5 {
			t.Errorf("Step = %v, want %v", *prop.Input.Step, 5)
		}
		if !prop.Optional {
			t.Errorf("Optional = %v, want %v", prop.Optional, true)
		}
		if prop.Input.Type != InputTypeRange {
			t.Errorf("Input.Type = %v, want %v", prop.Input.Type, InputTypeRange)
		}
	})
}