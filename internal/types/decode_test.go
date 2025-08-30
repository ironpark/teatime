package types

import (
	"reflect"
	"testing"
)

// Test struct definitions
type UnmarshalTestStruct struct {
	Name        string    `prop:"Name"`
	Age         int       `prop:"Age"`
	Score       float64   `prop:"Score"`
	Active      bool      `prop:"Active"`
	Tags        []string  `prop:"Tags"`
	Config      map[string]string `prop:"Config"`
	Data        interface{} `prop:"Data"`
}

type SimpleStruct struct {
	Value string `prop:"Value"`
}

func TestUnmarshalProp(t *testing.T) {
	tests := []struct {
		name       string
		target     interface{}
		properties []Property
		wantError  bool
		validate   func(t *testing.T, target interface{})
	}{
		{
			name:   "basic types",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Name", Value: "John Doe"},
				{Key: "Age", Value: 30},
				{Key: "Score", Value: 95.5},
				{Key: "Active", Value: true},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				if ts.Name != "John Doe" {
					t.Errorf("Name = %v, want %v", ts.Name, "John Doe")
				}
				if ts.Age != 30 {
					t.Errorf("Age = %v, want %v", ts.Age, 30)
				}
				if ts.Score != 95.5 {
					t.Errorf("Score = %v, want %v", ts.Score, 95.5)
				}
				if !ts.Active {
					t.Errorf("Active = %v, want %v", ts.Active, true)
				}
			},
		},
		{
			name:   "type conversion - string to int",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Age", Value: "25"},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				if ts.Age != 25 {
					t.Errorf("Age = %v, want %v", ts.Age, 25)
				}
			},
		},
		{
			name:   "type conversion - string to bool",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Active", Value: "true"},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				if !ts.Active {
					t.Errorf("Active = %v, want %v", ts.Active, true)
				}
			},
		},
		{
			name:   "type conversion - int to bool",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Active", Value: 1},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				if !ts.Active {
					t.Errorf("Active = %v, want %v", ts.Active, true)
				}
			},
		},
		{
			name:   "slice conversion",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Tags", Value: []string{"go", "test", "unmarshal"}},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				expected := []string{"go", "test", "unmarshal"}
				if !reflect.DeepEqual(ts.Tags, expected) {
					t.Errorf("Tags = %v, want %v", ts.Tags, expected)
				}
			},
		},
		{
			name:   "slice type conversion",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Tags", Value: []interface{}{"go", "test", "unmarshal"}},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				expected := []string{"go", "test", "unmarshal"}
				if !reflect.DeepEqual(ts.Tags, expected) {
					t.Errorf("Tags = %v, want %v", ts.Tags, expected)
				}
			},
		},
		{
			name:   "map conversion",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Config", Value: map[string]string{"key1": "value1", "key2": "value2"}},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				expected := map[string]string{"key1": "value1", "key2": "value2"}
				if !reflect.DeepEqual(ts.Config, expected) {
					t.Errorf("Config = %v, want %v", ts.Config, expected)
				}
			},
		},
		{
			name:   "interface field",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Data", Value: map[string]interface{}{"nested": "value"}},
			},
			validate: func(t *testing.T, target interface{}) {
				ts := target.(*UnmarshalTestStruct)
				expected := map[string]interface{}{"nested": "value"}
				if !reflect.DeepEqual(ts.Data, expected) {
					t.Errorf("Data = %v, want %v", ts.Data, expected)
				}
			},
		},
		{
			name:   "nil value",
			target: &SimpleStruct{Value: "initial"},
			properties: []Property{
				{Key: "Value", Value: nil},
			},
			validate: func(t *testing.T, target interface{}) {
				ss := target.(*SimpleStruct)
				if ss.Value != "" {
					t.Errorf("Value = %v, want empty string", ss.Value)
				}
			},
		},
		{
			name:       "nil target",
			target:     nil,
			properties: []Property{},
			wantError:  true,
		},
		{
			name:       "non-pointer target",
			target:     UnmarshalTestStruct{},
			properties: []Property{},
			wantError:  true,
		},
		{
			name:       "pointer to non-struct",
			target:     new(string),
			properties: []Property{},
			wantError:  true,
		},
		{
			name:   "unknown field",
			target: &SimpleStruct{},
			properties: []Property{
				{Key: "UnknownField", Value: "test"},
			},
			wantError: true,
		},
		{
			name:   "invalid string to int conversion",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Age", Value: "not_a_number"},
			},
			wantError: true,
		},
		{
			name:   "invalid string to bool conversion",
			target: &UnmarshalTestStruct{},
			properties: []Property{
				{Key: "Active", Value: "maybe"},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalProp(tt.target, tt.properties)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("UnmarshalProp() error = nil, wantError = true")
				}
				return
			}
			
			if err != nil {
				t.Errorf("UnmarshalProp() error = %v, wantError = false", err)
				return
			}
			
			if tt.validate != nil {
				tt.validate(t, tt.target)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	tests := []struct {
		name       string
		fieldType  reflect.Type
		value      interface{}
		wantError  bool
		validate   func(t *testing.T, field reflect.Value)
	}{
		{
			name:      "bool from string true",
			fieldType: reflect.TypeOf(true),
			value:     "true",
			validate: func(t *testing.T, field reflect.Value) {
				if !field.Bool() {
					t.Errorf("Expected true, got %v", field.Bool())
				}
			},
		},
		{
			name:      "bool from int 0",
			fieldType: reflect.TypeOf(true),
			value:     0,
			validate: func(t *testing.T, field reflect.Value) {
				if field.Bool() {
					t.Errorf("Expected false, got %v", field.Bool())
				}
			},
		},
		{
			name:      "int from string",
			fieldType: reflect.TypeOf(int(0)),
			value:     "42",
			validate: func(t *testing.T, field reflect.Value) {
				if field.Int() != 42 {
					t.Errorf("Expected 42, got %v", field.Int())
				}
			},
		},
		{
			name:      "int from float",
			fieldType: reflect.TypeOf(int(0)),
			value:     42.7,
			validate: func(t *testing.T, field reflect.Value) {
				if field.Int() != 42 {
					t.Errorf("Expected 42, got %v", field.Int())
				}
			},
		},
		{
			name:      "uint from negative int",
			fieldType: reflect.TypeOf(uint(0)),
			value:     -1,
			wantError: true,
		},
		{
			name:      "float from string",
			fieldType: reflect.TypeOf(float64(0)),
			value:     "3.14",
			validate: func(t *testing.T, field reflect.Value) {
				if field.Float() != 3.14 {
					t.Errorf("Expected 3.14, got %v", field.Float())
				}
			},
		},
		{
			name:      "string from any type",
			fieldType: reflect.TypeOf(""),
			value:     123,
			validate: func(t *testing.T, field reflect.Value) {
				if field.String() != "123" {
					t.Errorf("Expected '123', got %v", field.String())
				}
			},
		},
		{
			name:      "invalid bool string",
			fieldType: reflect.TypeOf(true),
			value:     "invalid",
			wantError: true,
		},
		{
			name:      "invalid int string",
			fieldType: reflect.TypeOf(int(0)),
			value:     "invalid",
			wantError: true,
		},
		{
			name:      "invalid float string",
			fieldType: reflect.TypeOf(float64(0)),
			value:     "invalid",
			wantError: true,
		},
		{
			name:      "unsupported type",
			fieldType: reflect.TypeOf(make(chan int)),
			value:     "test",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.New(tt.fieldType).Elem()
			err := setFieldValue(field, tt.value)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("setFieldValue() error = nil, wantError = true")
				}
				return
			}
			
			if err != nil {
				t.Errorf("setFieldValue() error = %v, wantError = false", err)
				return
			}
			
			if tt.validate != nil {
				tt.validate(t, field)
			}
		})
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := &UnmarshalTestStruct{
		Name:   "Alice",
		Age:    25,
		Score:  88.5,
		Active: true,
		Tags:   []string{"developer", "golang"},
		Config: map[string]string{"theme": "dark", "lang": "en"},
		Data:   map[string]interface{}{"count": 42},
	}

	// Marshal to properties
	properties, err := MarshalProp(original)
	if err != nil {
		t.Fatalf("MarshalProp() error = %v", err)
	}

	// Unmarshal back to struct
	result := &UnmarshalTestStruct{}
	err = UnmarshalProp(result, properties)
	if err != nil {
		t.Fatalf("UnmarshalProp() error = %v", err)
	}

	// Compare values
	if result.Name != original.Name {
		t.Errorf("Name = %v, want %v", result.Name, original.Name)
	}
	if result.Age != original.Age {
		t.Errorf("Age = %v, want %v", result.Age, original.Age)
	}
	if result.Score != original.Score {
		t.Errorf("Score = %v, want %v", result.Score, original.Score)
	}
	if result.Active != original.Active {
		t.Errorf("Active = %v, want %v", result.Active, original.Active)
	}
	if !reflect.DeepEqual(result.Tags, original.Tags) {
		t.Errorf("Tags = %v, want %v", result.Tags, original.Tags)
	}
	if !reflect.DeepEqual(result.Config, original.Config) {
		t.Errorf("Config = %v, want %v", result.Config, original.Config)
	}
	if !reflect.DeepEqual(result.Data, original.Data) {
		t.Errorf("Data = %v, want %v", result.Data, original.Data)
	}
}

func TestOverflowHandling(t *testing.T) {
	type OverflowStruct struct {
		SmallInt  int8  `prop:"SmallInt"`
		SmallUint uint8 `prop:"SmallUint"`
	}

	tests := []struct {
		name       string
		properties []Property
		wantError  bool
	}{
		{
			name: "int8 overflow",
			properties: []Property{
				{Key: "SmallInt", Value: 1000}, // > 127
			},
			wantError: true,
		},
		{
			name: "uint8 overflow",
			properties: []Property{
				{Key: "SmallUint", Value: 1000}, // > 255
			},
			wantError: true,
		},
		{
			name: "valid small values",
			properties: []Property{
				{Key: "SmallInt", Value: 100},
				{Key: "SmallUint", Value: 200},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &OverflowStruct{}
			err := UnmarshalProp(target, tt.properties)
			
			if tt.wantError && err == nil {
				t.Errorf("UnmarshalProp() error = nil, wantError = true")
			}
			if !tt.wantError && err != nil {
				t.Errorf("UnmarshalProp() error = %v, wantError = false", err)
			}
		})
	}
}