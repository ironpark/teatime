package types

import (
	"fmt"
	"reflect"
	"strconv"
)

// UnmarshalProp applies property values to a struct based on property key-value pairs.
// The target struct must be passed as a pointer and contain exported fields.
// Property keys are matched against struct field names (case-sensitive).
//
// Usage:
//   var config MyStruct
//   properties := []Property{...}
//   err := UnmarshalProp(&config, properties)
//
// The function performs type conversion and validation:
//   - Values are converted to match the target field's Go type
//   - Bool properties accept: bool, string ("true"/"false"), numeric (0/1)
//   - Numeric properties accept: numeric types and parseable strings
//   - String properties accept: any type (converted via fmt.Sprintf)
//   - Date properties accept: time.Time, int64 (unix timestamp), string (RFC3339)
//   - Array/Map/JSON properties accept: compatible slice/map/interface{} types
//
// Returns an error if:
//   - Target is not a pointer to struct
//   - Property key doesn't match any struct field
//   - Value type conversion fails
//   - Required validation fails
func UnmarshalProp(target any, properties []Property) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to struct, got %v", val.Kind())
	}
	
	val = val.Elem()
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct, got pointer to %v", typ.Kind())
	}

	// Create a map of field names to field info for efficient lookup
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			fieldMap[field.Name] = val.Field(i)
		}
	}

	// Apply each property to the corresponding struct field
	for _, prop := range properties {
		fieldValue, exists := fieldMap[prop.Key]
		if !exists {
			return fmt.Errorf("property key '%s' does not match any exported struct field", prop.Key)
		}
		
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field '%s'", prop.Key)
		}

		if err := setFieldValue(fieldValue, prop.Value); err != nil {
			return fmt.Errorf("failed to set field '%s': %w", prop.Key, err)
		}
	}

	return nil
}

// setFieldValue converts and assigns a property value to a struct field.
// Handles type conversion between property values and Go struct field types.
func setFieldValue(field reflect.Value, value any) error {
	if value == nil {
		// Set zero value for nil
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	valueReflect := reflect.ValueOf(value)
	fieldType := field.Type()

	// Direct assignment if types match
	if valueReflect.Type().AssignableTo(fieldType) {
		field.Set(valueReflect)
		return nil
	}

	// Type-specific conversion logic
	switch fieldType.Kind() {
	case reflect.Bool:
		return setBoolField(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntField(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintField(field, value)
	case reflect.Float32, reflect.Float64:
		return setFloatField(field, value)
	case reflect.String:
		return setStringField(field, value)
	case reflect.Slice, reflect.Array:
		return setSliceField(field, value)
	case reflect.Map:
		return setMapField(field, value)
	case reflect.Interface:
		// For interface{} or any, assign directly
		field.Set(valueReflect)
		return nil
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}
}

// setBoolField converts various types to bool
func setBoolField(field reflect.Value, value any) error {
	switch v := value.(type) {
	case bool:
		field.SetBool(v)
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			field.SetBool(b)
		} else {
			return fmt.Errorf("cannot parse '%s' as bool", v)
		}
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		field.SetBool(val != 0)
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(v).Uint()
		field.SetBool(val != 0)
	case float32, float64:
		val := reflect.ValueOf(v).Float()
		field.SetBool(val != 0.0)
	default:
		return fmt.Errorf("cannot convert %T to bool", v)
	}
	return nil
}

// setIntField converts various types to int
func setIntField(field reflect.Value, value any) error {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		if field.OverflowInt(val) {
			return fmt.Errorf("value %d overflows %s", val, field.Type())
		}
		field.SetInt(val)
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(v).Uint()
		intVal := int64(val)
		if val > ^uint64(0)>>1 || field.OverflowInt(intVal) {
			return fmt.Errorf("value %d overflows %s", val, field.Type())
		}
		field.SetInt(intVal)
	case float32, float64:
		val := reflect.ValueOf(v).Float()
		intVal := int64(val)
		if field.OverflowInt(intVal) {
			return fmt.Errorf("value %f overflows %s", val, field.Type())
		}
		field.SetInt(intVal)
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			if field.OverflowInt(val) {
				return fmt.Errorf("value %d overflows %s", val, field.Type())
			}
			field.SetInt(val)
		} else {
			return fmt.Errorf("cannot parse '%s' as int", v)
		}
	case bool:
		if v {
			field.SetInt(1)
		} else {
			field.SetInt(0)
		}
	default:
		return fmt.Errorf("cannot convert %T to %s", v, field.Type())
	}
	return nil
}

// setUintField converts various types to uint
func setUintField(field reflect.Value, value any) error {
	switch v := value.(type) {
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(v).Uint()
		if field.OverflowUint(val) {
			return fmt.Errorf("value %d overflows %s", val, field.Type())
		}
		field.SetUint(val)
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		if val < 0 {
			return fmt.Errorf("cannot assign negative value %d to %s", val, field.Type())
		}
		uintVal := uint64(val)
		if field.OverflowUint(uintVal) {
			return fmt.Errorf("value %d overflows %s", val, field.Type())
		}
		field.SetUint(uintVal)
	case float32, float64:
		val := reflect.ValueOf(v).Float()
		if val < 0 {
			return fmt.Errorf("cannot assign negative value %f to %s", val, field.Type())
		}
		uintVal := uint64(val)
		if field.OverflowUint(uintVal) {
			return fmt.Errorf("value %f overflows %s", val, field.Type())
		}
		field.SetUint(uintVal)
	case string:
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			if field.OverflowUint(val) {
				return fmt.Errorf("value %d overflows %s", val, field.Type())
			}
			field.SetUint(val)
		} else {
			return fmt.Errorf("cannot parse '%s' as uint", v)
		}
	case bool:
		if v {
			field.SetUint(1)
		} else {
			field.SetUint(0)
		}
	default:
		return fmt.Errorf("cannot convert %T to %s", v, field.Type())
	}
	return nil
}

// setFloatField converts various types to float
func setFloatField(field reflect.Value, value any) error {
	switch v := value.(type) {
	case float32, float64:
		val := reflect.ValueOf(v).Float()
		if field.OverflowFloat(val) {
			return fmt.Errorf("value %f overflows %s", val, field.Type())
		}
		field.SetFloat(val)
	case int, int8, int16, int32, int64:
		val := float64(reflect.ValueOf(v).Int())
		if field.OverflowFloat(val) {
			return fmt.Errorf("value %f overflows %s", val, field.Type())
		}
		field.SetFloat(val)
	case uint, uint8, uint16, uint32, uint64:
		val := float64(reflect.ValueOf(v).Uint())
		if field.OverflowFloat(val) {
			return fmt.Errorf("value %f overflows %s", val, field.Type())
		}
		field.SetFloat(val)
	case string:
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			if field.OverflowFloat(val) {
				return fmt.Errorf("value %f overflows %s", val, field.Type())
			}
			field.SetFloat(val)
		} else {
			return fmt.Errorf("cannot parse '%s' as float", v)
		}
	case bool:
		if v {
			field.SetFloat(1.0)
		} else {
			field.SetFloat(0.0)
		}
	default:
		return fmt.Errorf("cannot convert %T to %s", v, field.Type())
	}
	return nil
}

// setStringField converts any type to string
func setStringField(field reflect.Value, value any) error {
	field.SetString(fmt.Sprintf("%v", value))
	return nil
}

// setSliceField converts compatible slice types
func setSliceField(field reflect.Value, value any) error {
	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() != reflect.Slice && valueReflect.Kind() != reflect.Array {
		return fmt.Errorf("cannot assign %s to slice field", valueReflect.Kind())
	}

	// Create new slice with same length
	newSlice := reflect.MakeSlice(field.Type(), valueReflect.Len(), valueReflect.Len())
	
	// Convert each element
	for i := 0; i < valueReflect.Len(); i++ {
		srcElem := valueReflect.Index(i)
		dstElem := newSlice.Index(i)
		
		if err := setFieldValue(dstElem, srcElem.Interface()); err != nil {
			return fmt.Errorf("failed to convert slice element at index %d: %w", i, err)
		}
	}
	
	field.Set(newSlice)
	return nil
}

// setMapField converts compatible map types
func setMapField(field reflect.Value, value any) error {
	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() != reflect.Map {
		return fmt.Errorf("cannot assign %s to map field", valueReflect.Kind())
	}

	// Create new map
	newMap := reflect.MakeMap(field.Type())
	
	// Convert each key-value pair
	for _, key := range valueReflect.MapKeys() {
		srcValue := valueReflect.MapIndex(key)
		
		// Convert key
		newKey := reflect.New(field.Type().Key()).Elem()
		if err := setFieldValue(newKey, key.Interface()); err != nil {
			return fmt.Errorf("failed to convert map key %v: %w", key.Interface(), err)
		}
		
		// Convert value
		newValue := reflect.New(field.Type().Elem()).Elem()
		if err := setFieldValue(newValue, srcValue.Interface()); err != nil {
			return fmt.Errorf("failed to convert map value for key %v: %w", key.Interface(), err)
		}
		
		newMap.SetMapIndex(newKey, newValue)
	}
	
	field.Set(newMap)
	return nil
}