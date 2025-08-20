package node

// Helper functions for creating NodeProperty instances

// StringProperty creates a string property
func StringProperty(key, name, description string, defaultValue string, optional bool) NodeProperty {
	return NodeProperty{
		Type:        String,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeText,
		},
	}
}

// TextAreaProperty creates a multi-line text property
func TextAreaProperty(key, name, description string, defaultValue string, rows int, optional bool) NodeProperty {
	return NodeProperty{
		Type:        String,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeTextarea,
			Rows: rows,
		},
	}
}

// IntProperty creates an integer property
func IntProperty(key, name, description string, defaultValue int64, optional bool) NodeProperty {
	return NodeProperty{
		Type:        Int64,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeNumber,
		},
	}
}

// IntRangeProperty creates an integer property with min/max range
func IntRangeProperty(key, name, description string, defaultValue, min, max int64, optional bool) NodeProperty {
	minFloat := float64(min)
	maxFloat := float64(max)
	return NodeProperty{
		Type:        Int64,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeRange,
			Min:  &minFloat,
			Max:  &maxFloat,
			Step: Float64Ptr(1),
		},
	}
}

// FloatProperty creates a float property
func FloatProperty(key, name, description string, defaultValue float64, optional bool) NodeProperty {
	return NodeProperty{
		Type:        Float64,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeNumber,
			Step: Float64Ptr(0.01),
		},
	}
}

// FloatRangeProperty creates a float property with min/max range
func FloatRangeProperty(key, name, description string, defaultValue, min, max, step float64, optional bool) NodeProperty {
	return NodeProperty{
		Type:        Float64,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeRange,
			Min:  &min,
			Max:  &max,
			Step: &step,
		},
	}
}

// BoolProperty creates a boolean property
func BoolProperty(key, name, description string, defaultValue bool, optional bool) NodeProperty {
	return NodeProperty{
		Type:        Bool,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeSwitch,
		},
	}
}

// SelectProperty creates a select property with options
func SelectProperty(key, name, description string, options []string, defaultValue string, optional bool) NodeProperty {
	return NodeProperty{
		Type:        String,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Options:     options,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeSelect,
		},
	}
}

// MultiSelectProperty creates a multi-select property with options
func MultiSelectProperty(key, name, description string, options []string, defaultValue []string, optional bool) NodeProperty {
	return NodeProperty{
		Type:        StringArray,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Options:     options,
		Optional:    optional,
		Input: &InputConfig{
			Type:     InputTypeMultiSelect,
			Multiple: true,
		},
	}
}

// JSONProperty creates a JSON property
func JSONProperty(key, name, description string, defaultValue map[string]any, optional bool) NodeProperty {
	return NodeProperty{
		Type:        JSON,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
		Input: &InputConfig{
			Type: InputTypeTextarea,
			Rows: 10,
		},
	}
}

// StringArrayProperty creates a string array property
func StringArrayProperty(key, name, description string, defaultValue []string, optional bool) NodeProperty {
	return NodeProperty{
		Type:        StringArray,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
	}
}

// NumberArrayProperty creates a number array property
func NumberArrayProperty(key, name, description string, defaultValue []float64, optional bool) NodeProperty {
	return NodeProperty{
		Type:        NumberArray,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
	}
}

// BooleanArrayProperty creates a boolean array property
func BooleanArrayProperty(key, name, description string, defaultValue []bool, optional bool) NodeProperty {
	return NodeProperty{
		Type:        BooleanArray,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       defaultValue,
		Optional:    optional,
	}
}

// ReadOnlyProperty creates a read-only property
func ReadOnlyProperty(key, name, description string, value any) NodeProperty {
	// Determine type from value
	propertyType := String
	switch value.(type) {
	case bool:
		propertyType = Bool
	case int, int8, int16, int32, int64:
		propertyType = Int64
	case uint, uint8, uint16, uint32, uint64:
		propertyType = Uint64
	case float32, float64:
		propertyType = Float64
	case []string:
		propertyType = StringArray
	case []float64, []float32, []int, []int64:
		propertyType = NumberArray
	case []bool:
		propertyType = BooleanArray
	case map[string]any, []any:
		propertyType = JSON
	}
	
	return NodeProperty{
		Type:        propertyType,
		Key:         key,
		Name:        name,
		Description: description,
		Value:       value,
		Optional:    true,
		ReadOnly:    true,
	}
}

// OutputProperty creates an output property (always read-only)
func OutputProperty(key, name, description string, propertyType PropertyType) NodeProperty {
	return NodeProperty{
		Type:        propertyType,
		Key:         key,
		Name:        name,
		Description: description,
		Optional:    true,
		ReadOnly:    true,
	}
}