package node

import (
	"errors"
	
	"github.com/samber/lo"
)

// Type checking helper functions for property validation

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

// isString checks if the value is a string
func isString(v any) bool {
	_, ok := v.(string)
	return ok
}

// isBool checks if the value is a boolean
func isBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

// isWholeNumber checks if a float64 value represents a whole number
func isWholeNumber(f float64) bool {
	return f == float64(int64(f))
}

// isPositiveWholeNumber checks if a float64 value represents a positive whole number
func isPositiveWholeNumber(f float64) bool {
	return f >= 0 && f == float64(uint64(f))
}

// validateTypedArray validates if value is a specific typed array or []any with all elements of the checker type
func validateTypedArray[T any](v any, checker func(any) bool) error {
	// Check if it's already the correct typed array
	if _, ok := v.(T); ok {
		return nil
	}
	
	// Check if it's []any with all elements passing the checker
	if arr, ok := v.([]any); ok {
		if !lo.EveryBy(arr, checker) {
			return errors.New("all elements must be of the correct type")
		}
		return nil
	}
	
	return errors.New("invalid array type")
}