package node

import (
	"fmt"
	"testing"
)

// Example demonstrates various ways to create properties using the builder pattern
func Example_propertyBuilder() {
	// Simple string property
	name := StringProp("name", "Name",
		RequiredWithDefault("John Doe"),
		WithPlaceholder("Enter your name"),
	)
	fmt.Printf("Name property: %s (required=%v)\n", name.Name, !name.Optional)
	
	// Integer with percentage slider
	progress := IntProp("progress", "Progress",
		Percentage(),
		OptionalWithDefault(int64(50)),
	)
	fmt.Printf("Progress: min=%v, max=%v\n", *progress.Input.Min, *progress.Input.Max)
	
	// Boolean with toggle
	enabled := BoolProp("enabled", "Enable Feature",
		Toggle(),
		WithDefault(true),
	)
	fmt.Printf("Enabled: type=%v, default=%v\n", enabled.Input.Type, enabled.Value)
	
	// Text area with preset
	description := StringProp("desc", "Description",
		TextArea(10),
		Optional(),
	)
	fmt.Printf("Description: rows=%v\n", description.Input.Rows)
	
	// Output:
	// Name property: Name (required=true)
	// Progress: min=0, max=100
	// Enabled: type=7, default=true
	// Description: rows=10
}

// TestPropertyBuilderUsagePatterns shows common usage patterns
func TestPropertyBuilderUsagePatterns(t *testing.T) {
	t.Run("Compact property creation", func(t *testing.T) {
		// Very compact for simple cases
		prop1 := StringProp("key", "Label", Optional())
		prop2 := IntProp("count", "Count", RequiredWithDefault(int64(0)))
		prop3 := BoolProp("flag", "Flag", Toggle())
		
		if prop1.Optional != true {
			t.Error("prop1 should be optional")
		}
		if prop2.Value != int64(0) {
			t.Error("prop2 should have default value 0")
		}
		if prop3.Input.Type != InputTypeSwitch {
			t.Error("prop3 should be a switch")
		}
	})
	
	t.Run("Range properties", func(t *testing.T) {
		// Volume with custom range
		volume := IntProp("volume", "Volume",
			RangeSlider(0, 100),
			WithDefault(int64(50)),
		)
		
		if volume.Input.Type != InputTypeRange {
			t.Error("Should be a range input")
		}
		if *volume.Input.Step != 1 { // Auto-calculated for Int64
			t.Error("Step should be 1 for integer")
		}
		
		// Opacity with float range
		opacity := FloatProp("opacity", "Opacity",
			RangeSlider(0, 1),
			WithDefault(0.8),
		)
		
		if *opacity.Input.Step == 1 {
			t.Error("Step should not be 1 for float")
		}
	})
	
	t.Run("Select with options", func(t *testing.T) {
		sizes := []string{"small", "medium", "large"}
		
		// Single select
		size := SelectProp("size", "Size", sizes,
			RequiredWithDefault("medium"),
		)
		
		if size.Optional {
			t.Error("Should be required")
		}
		if size.Value != "medium" {
			t.Error("Should have default value 'medium'")
		}
		
		// Multi select
		tags := MultiSelectProp("tags", "Tags", sizes,
			OptionalWithDefault([]string{"small", "large"}),
		)
		
		if !tags.Optional {
			t.Error("Should be optional")
		}
	})
	
	t.Run("Complex property with many options", func(t *testing.T) {
		prop := FloatProp("temperature", "Temperature",
			WithDescription("Set the temperature in Celsius"),
			RangeSlider(-20, 50),
			RequiredWithDefault(20.0),
			WithPlaceholder("20.0"),
		)
		
		if prop.Description != "Set the temperature in Celsius" {
			t.Error("Description not set")
		}
		if *prop.Input.Min != -20 || *prop.Input.Max != 50 {
			t.Error("Range not set correctly")
		}
		if prop.Optional {
			t.Error("Should be required")
		}
		if prop.Value != 20.0 {
			t.Error("Default value not set")
		}
	})
}