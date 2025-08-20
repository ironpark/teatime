package node

import (
	"testing"
	"time"
)

func TestNodeProperty_ValidateValue(t *testing.T) {
	tests := []struct {
		name     string
		prop     NodeProperty
		value    any
		wantErr  bool
	}{
		// Bool type tests
		{
			name:     "bool valid true",
			prop:     NodeProperty{Type: Bool},
			value:    true,
			wantErr:  false,
		},
		{
			name:     "bool valid false",
			prop:     NodeProperty{Type: Bool},
			value:    false,
			wantErr:  false,
		},
		{
			name:     "bool invalid string",
			prop:     NodeProperty{Type: Bool},
			value:    "true",
			wantErr:  true,
		},
		{
			name:     "bool nil required",
			prop:     NodeProperty{Type: Bool, Optional: false},
			value:    nil,
			wantErr:  true,
		},
		{
			name:     "bool nil optional",
			prop:     NodeProperty{Type: Bool, Optional: true},
			value:    nil,
			wantErr:  false,
		},
		
		// Int64 type tests
		{
			name:     "int64 valid int64",
			prop:     NodeProperty{Type: Int64},
			value:    int64(42),
			wantErr:  false,
		},
		{
			name:     "int64 valid int",
			prop:     NodeProperty{Type: Int64},
			value:    42,
			wantErr:  false,
		},
		{
			name:     "int64 valid float64 whole number",
			prop:     NodeProperty{Type: Int64},
			value:    float64(42),
			wantErr:  false,
		},
		{
			name:     "int64 invalid float64 decimal",
			prop:     NodeProperty{Type: Int64},
			value:    42.5,
			wantErr:  true,
		},
		{
			name:     "int64 invalid string",
			prop:     NodeProperty{Type: Int64},
			value:    "42",
			wantErr:  true,
		},
		
		// Uint64 type tests
		{
			name:     "uint64 valid uint64",
			prop:     NodeProperty{Type: Uint64},
			value:    uint64(42),
			wantErr:  false,
		},
		{
			name:     "uint64 valid float64 positive whole",
			prop:     NodeProperty{Type: Uint64},
			value:    float64(42),
			wantErr:  false,
		},
		{
			name:     "uint64 invalid negative float64",
			prop:     NodeProperty{Type: Uint64},
			value:    float64(-42),
			wantErr:  true,
		},
		{
			name:     "uint64 invalid float64 decimal",
			prop:     NodeProperty{Type: Uint64},
			value:    42.5,
			wantErr:  true,
		},
		
		// Float64 type tests
		{
			name:     "float64 valid float64",
			prop:     NodeProperty{Type: Float64},
			value:    42.5,
			wantErr:  false,
		},
		{
			name:     "float64 valid int",
			prop:     NodeProperty{Type: Float64},
			value:    42,
			wantErr:  false,
		},
		{
			name:     "float64 valid int64",
			prop:     NodeProperty{Type: Float64},
			value:    int64(42),
			wantErr:  false,
		},
		{
			name:     "float64 invalid string",
			prop:     NodeProperty{Type: Float64},
			value:    "42.5",
			wantErr:  true,
		},
		
		// String type tests
		{
			name:     "string valid",
			prop:     NodeProperty{Type: String},
			value:    "hello",
			wantErr:  false,
		},
		{
			name:     "string invalid int",
			prop:     NodeProperty{Type: String},
			value:    42,
			wantErr:  true,
		},
		
		// JSON type tests
		{
			name:     "json valid map",
			prop:     NodeProperty{Type: JSON},
			value:    map[string]any{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "json valid array",
			prop:     NodeProperty{Type: JSON},
			value:    []any{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "json invalid string",
			prop:     NodeProperty{Type: JSON},
			value:    "not json",
			wantErr:  true,
		},
		
		// Date type tests
		{
			name:     "date valid time.Time",
			prop:     NodeProperty{Type: Date},
			value:    time.Now(),
			wantErr:  false,
		},
		{
			name:     "date valid string",
			prop:     NodeProperty{Type: Date},
			value:    "2024-01-01T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "date invalid int",
			prop:     NodeProperty{Type: Date},
			value:    42,
			wantErr:  true,
		},
		
		// StringArray type tests
		{
			name:     "string array valid []string",
			prop:     NodeProperty{Type: StringArray},
			value:    []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "string array valid []any with strings",
			prop:     NodeProperty{Type: StringArray},
			value:    []any{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "string array invalid []any with mixed types",
			prop:     NodeProperty{Type: StringArray},
			value:    []any{"a", 1, "c"},
			wantErr:  true,
		},
		{
			name:     "string array invalid single string",
			prop:     NodeProperty{Type: StringArray},
			value:    "not an array",
			wantErr:  true,
		},
		
		// NumberArray type tests
		{
			name:     "number array valid []float64",
			prop:     NodeProperty{Type: NumberArray},
			value:    []float64{1.0, 2.0, 3.0},
			wantErr:  false,
		},
		{
			name:     "number array valid []any with numbers",
			prop:     NodeProperty{Type: NumberArray},
			value:    []any{float64(1), int(2), int64(3)},
			wantErr:  false,
		},
		{
			name:     "number array invalid []any with string",
			prop:     NodeProperty{Type: NumberArray},
			value:    []any{1, "two", 3},
			wantErr:  true,
		},
		
		// BooleanArray type tests
		{
			name:     "boolean array valid []bool",
			prop:     NodeProperty{Type: BooleanArray},
			value:    []bool{true, false, true},
			wantErr:  false,
		},
		{
			name:     "boolean array valid []any with bools",
			prop:     NodeProperty{Type: BooleanArray},
			value:    []any{true, false, true},
			wantErr:  false,
		},
		{
			name:     "boolean array invalid []any with mixed",
			prop:     NodeProperty{Type: BooleanArray},
			value:    []any{true, "false", true},
			wantErr:  true,
		},
		
		// Invalid property type
		{
			name:     "invalid property type",
			prop:     NodeProperty{Type: PropertyType(999)},
			value:    "anything",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prop.ValidateValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}