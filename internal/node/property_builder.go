package node

// PropertyOption is a function that modifies a NodeProperty
type PropertyOption func(*NodeProperty)

// WithDescription sets the description
func WithDescription(description string) PropertyOption {
	return func(p *NodeProperty) {
		p.Description = description
	}
}

// WithDefault sets the default value
func WithDefault(value any) PropertyOption {
	return func(p *NodeProperty) {
		p.Value = value
	}
}

// Optional makes the property optional
func Optional() PropertyOption {
	return func(p *NodeProperty) {
		p.Optional = true
	}
}

// Required makes the property required (default)
func Required() PropertyOption {
	return func(p *NodeProperty) {
		p.Optional = false
	}
}

// ReadOnly makes the property read-only
func ReadOnly() PropertyOption {
	return func(p *NodeProperty) {
		p.ReadOnly = true
	}
}

// WithOptions sets the available options for select properties
func WithOptions(options ...string) PropertyOption {
	return func(p *NodeProperty) {
		p.Options = options
	}
}

// WithInput sets the input configuration
func WithInput(inputType InputType) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = inputType
	}
}

// WithPlaceholder sets the placeholder for text inputs
func WithPlaceholder(placeholder string) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Placeholder = placeholder
	}
}

// WithRows sets the number of rows for textarea
func WithRows(rows int) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Rows = rows
	}
}

// WithRange sets min, max, and step for range inputs
func WithRange(min, max, step float64) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeRange
		p.Input.Min = &min
		p.Input.Max = &max
		p.Input.Step = &step
	}
}

// WithMin sets the minimum value
func WithMin(min float64) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Min = &min
	}
}

// WithMax sets the maximum value
func WithMax(max float64) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Max = &max
	}
}

// WithStep sets the step value
func WithStep(step float64) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Step = &step
	}
}

// Multiple enables multiple selection
func Multiple() PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Multiple = true
	}
}

// Expression creates an expression input
func Expression() PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeExpression
	}
}

// KeyValue creates a key-value pairs input
func KeyValue() PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeKeyValue
	}
}

// Common preset options for convenience

// RequiredWithDefault sets a property as required with a default value
func RequiredWithDefault(value any) PropertyOption {
	return func(p *NodeProperty) {
		p.Optional = false
		p.Value = value
	}
}

// OptionalWithDefault sets a property as optional with a default value
func OptionalWithDefault(value any) PropertyOption {
	return func(p *NodeProperty) {
		p.Optional = true
		p.Value = value
	}
}

// RangeSlider creates a range slider input with common defaults
func RangeSlider(min, max float64) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeRange
		p.Input.Min = &min
		p.Input.Max = &max

		// Auto-calculate step based on range
		step := (max - min) / 100
		if p.Type == Int64 || p.Type == Uint64 {
			step = 1
		}
		p.Input.Step = &step
	}
}

// TextArea creates a textarea input with specified rows
func TextArea(rows int) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeTextarea
		p.Input.Rows = rows
	}
}

// Percentage creates a percentage input (0-100)
func Percentage() PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		min, max, step := 0.0, 100.0, 1.0
		p.Input.Type = InputTypeRange
		p.Input.Min = &min
		p.Input.Max = &max
		p.Input.Step = &step
	}
}

// Toggle creates a toggle switch input
func Toggle() PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeSwitch
	}
}

// DynamicList creates a dynamic list input
func DynamicList(uniqueSet bool) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		p.Input.Type = InputTypeDynamicList
		p.Input.Unique = uniqueSet
	}
}

// Validation helpers

// WithMinLength sets minimum length validation (for strings)
func WithMinLength(min int) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		minVal := float64(min)
		p.Input.Min = &minVal
	}
}

// WithMaxLength sets maximum length validation (for strings)
func WithMaxLength(max int) PropertyOption {
	return func(p *NodeProperty) {
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		maxVal := float64(max)
		p.Input.Max = &maxVal
	}
}

// WithPattern sets a regex pattern for validation
func WithPattern(pattern string) PropertyOption {
	return func(p *NodeProperty) {
		// Store pattern in placeholder for now (could be extended later)
		if p.Input == nil {
			p.Input = &InputConfig{}
		}
		// TODO: Add Pattern field to InputConfig
	}
}

// NewProp creates a new NodeProperty with the given options
func NewProp(propertyType PropertyType, key, name string, opts ...PropertyOption) NodeProperty {
	prop := NodeProperty{
		Type:     propertyType,
		Key:      key,
		Name:     name,
		Optional: false, // Required by default
		Input:    &InputConfig{},
	}

	// Set default input type based on property type
	switch propertyType {
	case Bool:
		prop.Input.Type = InputTypeSwitch
	case String:
		prop.Input.Type = InputTypeText
	case Int64, Uint64, Float64:
		prop.Input.Type = InputTypeNumber
	case JSON, XML:
		prop.Input.Type = InputTypeTextarea
		prop.Input.Rows = 10
	case StringArray:
		prop.Input.Type = InputTypeMultiSelect
		prop.Input.Multiple = true
	}

	// Apply options
	for _, opt := range opts {
		opt(&prop)
	}

	return prop
}

// Convenience constructors using the builder pattern

// StringProp creates a string property
func StringProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(String, key, name, opts...)
}

// TextProp creates a multi-line text property
func TextProp(key, name string, opts ...PropertyOption) NodeProperty {
	defaultOpts := []PropertyOption{
		WithInput(InputTypeTextarea),
		WithRows(5),
	}
	return NewProp(String, key, name, append(defaultOpts, opts...)...)
}

// IntProp creates an integer property
func IntProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(Int64, key, name, opts...)
}

// FloatProp creates a float property
func FloatProp(key, name string, opts ...PropertyOption) NodeProperty {
	defaultOpts := []PropertyOption{WithStep(0.01)}
	return NewProp(Float64, key, name, append(defaultOpts, opts...)...)
}

// BoolProp creates a boolean property
func BoolProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(Bool, key, name, opts...)
}

// SelectProp creates a select property
func SelectProp(key, name string, options []string, opts ...PropertyOption) NodeProperty {
	defaultOpts := []PropertyOption{
		WithInput(InputTypeSelect),
		WithOptions(options...),
	}
	return NewProp(String, key, name, append(defaultOpts, opts...)...)
}

// MultiSelectProp creates a multi-select property
func MultiSelectProp(key, name string, options []string, opts ...PropertyOption) NodeProperty {
	defaultOpts := []PropertyOption{
		WithOptions(options...),
		// Input type and Multiple are already set by default for StringArray
	}
	return NewProp(StringArray, key, name, append(defaultOpts, opts...)...)
}

// JSONProp creates a JSON property
func JSONProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(JSON, key, name, opts...)
}

// StringArrayProp creates a string array property
func StringArrayProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(StringArray, key, name, opts...)
}

// NumberArrayProp creates a number array property
func NumberArrayProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(NumberArray, key, name, opts...)
}

// BoolArrayProp creates a boolean array property
func BoolArrayProp(key, name string, opts ...PropertyOption) NodeProperty {
	return NewProp(BooleanArray, key, name, opts...)
}

// OutputProp creates an output property (always read-only and optional)
func OutputProp(propertyType PropertyType, key, name string, opts ...PropertyOption) NodeProperty {
	defaultOpts := []PropertyOption{ReadOnly(), Optional()}
	return NewProp(propertyType, key, name, append(defaultOpts, opts...)...)
}
