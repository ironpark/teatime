package node

import (
	"errors"
	"fmt"
	"time"
)

// Validator is a function that validates a value
type Validator func(v any) error

// validators map for each PropertyType
var validators = map[PropertyType]Validator{
	Bool:         validateBool,
	Int64:        validateInt64,
	Uint64:       validateUint64,
	Float64:      validateFloat64,
	String:       validateString,
	JSON:         validateJSON,
	XML:          validateXML,
	Date:         validateDate,
	StringArray:  validateStringArray,
	NumberArray:  validateNumberArray,
	BooleanArray: validateBooleanArray,
}

func validateBool(v any) error {
	if _, ok := v.(bool); !ok {
		return errors.New("invalid boolean value")
	}
	return nil
}

func validateInt64(v any) error {
	switch val := v.(type) {
	case int64, int32, int16, int8, int:
		return nil
	case float64:
		if !isWholeNumber(val) {
			return errors.New("invalid int64 value: not a whole number")
		}
		return nil
	default:
		return errors.New("invalid int64 value")
	}
}

func validateUint64(v any) error {
	switch val := v.(type) {
	case uint64, uint32, uint16, uint8, uint:
		return nil
	case float64:
		if !isPositiveWholeNumber(val) {
			return errors.New("invalid uint64 value: must be positive whole number")
		}
		return nil
	default:
		return errors.New("invalid uint64 value")
	}
}

func validateFloat64(v any) error {
	if !isNumber(v) {
		return errors.New("invalid float64 value")
	}
	return nil
}

func validateString(v any) error {
	if _, ok := v.(string); !ok {
		return errors.New("invalid string value")
	}
	return nil
}

func validateJSON(v any) error {
	switch v.(type) {
	case map[string]any, []any:
		return nil
	default:
		return errors.New("invalid json value")
	}
}

func validateXML(v any) error {
	if _, ok := v.(string); !ok {
		return errors.New("invalid xml value")
	}
	return nil
}

func validateDate(v any) error {
	switch v.(type) {
	case time.Time:
		return nil
	case string:
		// Allow string representation that can be parsed
		return nil
	default:
		return errors.New("invalid date value")
	}
}

func validateStringArray(v any) error {
	if err := validateTypedArray[[]string](v, isString); err != nil {
		return errors.New("invalid string array value: all elements must be strings")
	}
	return nil
}

func validateNumberArray(v any) error {
	if err := validateTypedArray[[]float64](v, isNumber); err != nil {
		return errors.New("invalid number array value: all elements must be numbers")
	}
	return nil
}

func validateBooleanArray(v any) error {
	if err := validateTypedArray[[]bool](v, isBool); err != nil {
		return errors.New("invalid boolean array value: all elements must be booleans")
	}
	return nil
}

// ValidatePropertyValue validates a value against a property type
func ValidatePropertyValue(propertyType PropertyType, value any, optional bool) error {
	if value == nil {
		if optional {
			return nil
		}
		return errors.New("value is required")
	}
	
	validator, ok := validators[propertyType]
	if !ok {
		return fmt.Errorf("invalid property type: %v", propertyType)
	}
	
	return validator(value)
}