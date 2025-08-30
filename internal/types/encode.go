package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MarshalProp converts a struct to a list of properties based on struct tags.
//
// Supported tag formats:
//   - `prop:"name"` - Basic property with display name
//   - `prop:"name,optional"` - Optional property
//   - `prop:"name,hide"` - Hidden from node preview
//   - `prop:"name,enum(opt1,opt2,opt3)"` - Property with enum options
//   - `prop:"name,optional,hide"` - Combined options
//   - `description:"Property description"` - Property description
//   - `input:"type,option=value"` - Input configuration
//
// Input types: text, textarea, expression, number, range, select, multi-select
// Input options: min, max, step, placeholder, rows, multiple, unique
func MarshalProp(v any) ([]Property, error) {
	if v == nil {
		return nil, fmt.Errorf("value cannot be nil")
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("map and struct must be passed as pointers, got %v", val.Kind())
	}
	val = val.Elem()
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s", typ.Name())
	}

	var properties []Property
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		prop, err := marshalField(field, fieldValue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s: %w", field.Name, err)
		}
		if prop != nil {
			properties = append(properties, *prop)
		}
	}

	return properties, nil
}

func marshalField(field reflect.StructField, value reflect.Value) (*Property, error) {
	propTag := field.Tag.Get("prop")
	if propTag == "" || propTag == "-" {
		return nil, nil
	}

	prop := Property{
		Key:         field.Name,
		Description: field.Tag.Get("description"),
		Value:       value.Interface(),
	}

	// Parse prop tag - handle enum specially since it contains commas
	var cleanParts []string
	var currentPart strings.Builder
	inEnum := false

	for _, char := range propTag {
		switch char {
		case '(':
			if strings.HasSuffix(currentPart.String(), "enum") {
				inEnum = true
			}
			currentPart.WriteRune(char)
		case ')':
			currentPart.WriteRune(char)
			if inEnum {
				inEnum = false
			}
		case ',':
			if inEnum {
				currentPart.WriteRune(char)
			} else {
				cleanParts = append(cleanParts, currentPart.String())
				currentPart.Reset()
			}
		default:
			currentPart.WriteRune(char)
		}
	}
	if currentPart.Len() > 0 {
		cleanParts = append(cleanParts, currentPart.String())
	}

	if len(cleanParts) > 0 {
		prop.Name = strings.TrimSpace(cleanParts[0])
	}

	// Parse the remaining parts
	for i, part := range cleanParts[1:] {
		part = strings.TrimSpace(part)
		switch {
		case part == "optional":
			prop.Optional = true
		case part == "hide":
			prop.HideOnPreview = true
		case strings.HasPrefix(part, "enum(") && strings.HasSuffix(part, ")"):
			enumStr := part[5 : len(part)-1]
			if enumStr != "" {
				enums := strings.Split(enumStr, ",")
				for _, e := range enums {
					trimmed := strings.TrimSpace(e)
					if trimmed != "" {
						prop.Enums = append(prop.Enums, trimmed)
					}
				}
			}
		}
		_ = i // avoid unused variable warning
	}

	// Determine property type
	propType, err := getPropertyType(field.Type)
	if err != nil {
		return nil, err
	}
	prop.Type = propType

	// Set ValType for array and map typess

	switch propType {
	case Array:
		if field.Type.Elem() != nil {
			valType, err := getPropertyType(field.Type.Elem())
			if err == nil {
				prop.ValType = valType
			}
		}
	case Map:
		if field.Type.Elem() != nil {
			valType, err := getPropertyType(field.Type.Elem())
			if err == nil {
				prop.ValType = valType
			}
		}
	}

	// Parse input configuration
	inputTag := field.Tag.Get("input")
	if inputTag != "" {
		inputConfig, err := parseInputConfig(inputTag)
		if err != nil {
			return nil, fmt.Errorf("failed to parse input config: %w", err)
		}
		prop.Input = inputConfig
	} else {
		// Set fallback input type if no input config is provided
		prop.Input = &InputConfig{
			Type: propType.FallbackInputType(),
		}
	}

	// Validate that the input configuration is compatible with the property type
	if err := prop.Input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input config for field %s: %w", field.Name, err)
	}
	
	if err := prop.Input.ValidateType(propType); err != nil {
		return nil, fmt.Errorf("input type incompatible with property type for field %s: %w", field.Name, err)
	}

	return &prop, nil
}

func getPropertyType(t reflect.Type) (PropertyType, error) {
	switch t.Kind() {
	case reflect.Bool:
		return Bool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Int64, nil
	case reflect.Float32, reflect.Float64:
		return Float64, nil
	case reflect.String:
		return String, nil
	case reflect.Slice, reflect.Array:
		return Array, nil
	case reflect.Map:
		return Map, nil
	case reflect.Interface:
		return JSON, nil
	default:
		return 0, fmt.Errorf("unsupported type: %s", t.Kind())
	}
}

func parseInputConfig(inputTag string) (*InputConfig, error) {
	config := &InputConfig{}

	parts := strings.Split(inputTag, ",")
	if len(parts) == 0 {
		return config, nil
	}

	// First part is the input type
	inputTypeStr := strings.TrimSpace(parts[0])
	// Handle alias
	if inputTypeStr == "key-value" {
		inputTypeStr = "kv"
	}

	inputType, err := InputTypeFromString(inputTypeStr)
	if err != nil {
		// Default to text input if unknown type
		config.Type = InputTypeText
	} else {
		config.Type = inputType
	}

	// Parse additional options
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
			key, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			switch key {
			case "min":
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					config.Min = &val
				}
			case "max":
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					config.Max = &val
				}
			case "step":
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					config.Step = &val
				}
			case "placeholder":
				config.Placeholder = value
			}
		} else {
			switch part {
			case "multiple":
				config.Multiple = true
			case "unique":
				config.Unique = true
			}
		}
	}

	return config, nil
}