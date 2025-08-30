package types

import (
	"testing"
)

type TestStruct struct {
	Name         string   `prop:"name" description:"The name field"`
	Age          int      `prop:"age,optional" description:"Age in years" input:"number,min=0,max=150"`
	Email        string   `prop:"email" description:"Email address" input:"text,placeholder=Enter your email"`
	IsActive     bool     `prop:"active" description:"Whether the user is active"`
	Score        float64  `prop:"score,hide" description:"User score" input:"range,min=0,max=100,step=0.1"`
	Category     string   `prop:"category,enum(basic,premium,enterprise)" description:"User category" input:"select"`
	Tags         []string `prop:"tags" description:"User tags"`
	unexported   string   // Should be ignored
	SkippedField string   `prop:"-" description:"This should be skipped"`
}

func TestMarshalProp_BasicTypes(t *testing.T) {
	testStruct := &TestStruct{
		Name:         "John Doe",
		Age:          30,
		Email:        "john@example.com",
		IsActive:     true,
		Score:        85.5,
		Category:     "premium",
		Tags:         []string{"developer", "go"},
		unexported:   "should be ignored",
		SkippedField: "should be skipped",
	}

	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedFields := map[string]struct {
		name          string
		propType      PropertyType
		optional      bool
		hideOnPreview bool
		hasEnums      bool
		enumCount     int
		hasInput      bool
		valType       PropertyType
	}{
		"Name":     {"name", String, false, false, false, 0, true, 0},
		"Age":      {"age", Int64, true, false, false, 0, true, 0},
		"Email":    {"email", String, false, false, false, 0, true, 0},
		"IsActive": {"active", Bool, false, false, false, 0, true, 0},
		"Score":    {"score", Float64, false, true, false, 0, true, 0},
		"Category": {"category", String, false, false, true, 3, true, 0},
		"Tags":     {"tags", Array, false, false, false, 0, true, String},
	}

	if len(props) != len(expectedFields) {
		t.Errorf("Expected %d properties, got %d", len(expectedFields), len(props))
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	for key, expected := range expectedFields {
		prop, exists := propMap[key]
		if !exists {
			t.Errorf("Expected property %s not found", key)
			continue
		}

		if prop.Name != expected.name {
			t.Errorf("Property %s: expected name %s, got %s", key, expected.name, prop.Name)
		}

		if prop.Type != expected.propType {
			t.Errorf("Property %s: expected type %s, got %s", key, expected.propType, prop.Type)
		}

		if prop.Optional != expected.optional {
			t.Errorf("Property %s: expected optional %v, got %v", key, expected.optional, prop.Optional)
		}

		if prop.HideOnPreview != expected.hideOnPreview {
			t.Errorf("Property %s: expected hideOnPreview %v, got %v", key, expected.hideOnPreview, prop.HideOnPreview)
		}

		if expected.hasEnums {
			if len(prop.Enums) != expected.enumCount {
				t.Errorf("Property %s: expected %d enums, got %d", key, expected.enumCount, len(prop.Enums))
			}
		}

		if expected.hasInput && prop.Input == nil {
			t.Errorf("Property %s: expected input config, got nil", key)
		} else if !expected.hasInput && prop.Input != nil {
			t.Errorf("Property %s: expected no input config, got %v", key, prop.Input)
		}

		if expected.valType != 0 && prop.ValType != expected.valType {
			t.Errorf("Property %s: expected valType %s, got %s", key, expected.valType, prop.ValType)
		}
	}
}

func TestMarshalProp_TagParsing(t *testing.T) {
	type TagTestStruct struct {
		OptionalField string `prop:"optional_field,optional" description:"Optional field"`
		HiddenField   string `prop:"hidden_field,hide" description:"Hidden field"`
		EnumField     string `prop:"enum_field,enum(option1,option2,option3)" description:"Enum field"`
		CombinedField string `prop:"combined_field,optional,hide" description:"Combined options"`
	}

	testStruct := &TagTestStruct{
		OptionalField: "test",
		HiddenField:   "hidden",
		EnumField:     "option1",
		CombinedField: "combined",
	}

	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	// Test optional field
	if prop, exists := propMap["OptionalField"]; exists {
		if !prop.Optional {
			t.Error("OptionalField should be optional")
		}
	}

	// Test hidden field
	if prop, exists := propMap["HiddenField"]; exists {
		if !prop.HideOnPreview {
			t.Error("HiddenField should be hidden on preview")
		}
	}

	// Test enum field
	if prop, exists := propMap["EnumField"]; exists {
		if len(prop.Enums) != 3 {
			t.Errorf("EnumField should have 3 enums, got %d, enums: %+v", len(prop.Enums), prop.Enums)
		}
		expectedEnums := []string{"option1", "option2", "option3"}
		for i, enum := range prop.Enums {
			if enum != expectedEnums[i] {
				t.Errorf("EnumField enum %d: expected %s, got %v", i, expectedEnums[i], enum)
			}
		}
	} else {
		t.Error("EnumField not found in properties")
	}

	// Test combined field
	if prop, exists := propMap["CombinedField"]; exists {
		if !prop.Optional {
			t.Error("CombinedField should be optional")
		}
		if !prop.HideOnPreview {
			t.Error("CombinedField should be hidden on preview")
		}
	}
}

func TestMarshalProp_InputConfig(t *testing.T) {
	type InputTestStruct struct {
		TextInput     string            `prop:"text_input" input:"text,placeholder=Enter text"`
		NumberInput   int               `prop:"number_input" input:"number,min=0,max=100,step=5"`
		RangeInput    float64           `prop:"range_input" input:"range,min=0,max=1,step=0.1"`
		TextareaInput string            `prop:"textarea_input" input:"textarea,rows=5,placeholder=Enter description"`
		SelectInput   string            `prop:"select_input" input:"select,multiple"`
		SwitchInput   bool              `prop:"switch_input" input:"switch"`
		DateInput     string            `prop:"date_input" input:"expression"`
		ListInput     []string          `prop:"list_input" input:"list"`
		JsonInput     any               `prop:"json_input" input:"json"`
		KeyValueInput map[string]string `prop:"kv_input" input:"kv"`
	}

	testStruct := &InputTestStruct{}
	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	// Test text input
	if prop, exists := propMap["TextInput"]; exists && prop.Input != nil {
		if prop.Input.Type != InputTypeText {
			t.Errorf("TextInput: expected type %v, got %v", InputTypeText, prop.Input.Type)
		}
		if prop.Input.Placeholder != "Enter text" {
			t.Errorf("TextInput: expected placeholder 'Enter text', got '%s'", prop.Input.Placeholder)
		}
	}

	// Test number input
	if prop, exists := propMap["NumberInput"]; exists && prop.Input != nil {
		if prop.Input.Type != InputTypeNumber {
			t.Errorf("NumberInput: expected type %v, got %v", InputTypeNumber, prop.Input.Type)
		}
		if prop.Input.Min == nil || *prop.Input.Min != 0 {
			t.Errorf("NumberInput: expected min 0, got %v", prop.Input.Min)
		}
		if prop.Input.Max == nil || *prop.Input.Max != 100 {
			t.Errorf("NumberInput: expected max 100, got %v", prop.Input.Max)
		}
		if prop.Input.Step == nil || *prop.Input.Step != 5 {
			t.Errorf("NumberInput: expected step 5, got %v", prop.Input.Step)
		}
	}

	// Test range input
	if prop, exists := propMap["RangeInput"]; exists && prop.Input != nil {
		if prop.Input.Type != InputTypeRange {
			t.Errorf("RangeInput: expected type %v, got %v", InputTypeRange, prop.Input.Type)
		}
		if prop.Input.Step == nil || *prop.Input.Step != 0.1 {
			t.Errorf("RangeInput: expected step 0.1, got %v", prop.Input.Step)
		}
	}

	// Test textarea input
	if prop, exists := propMap["TextareaInput"]; exists && prop.Input != nil {
		if prop.Input.Type != InputTypeTextarea {
			t.Errorf("TextareaInput: expected type %v, got %v", InputTypeTextarea, prop.Input.Type)
		}
	}

	// Test select input
	if prop, exists := propMap["SelectInput"]; exists && prop.Input != nil {
		if prop.Input.Type != InputTypeSelect {
			t.Errorf("SelectInput: expected type %v, got %v", InputTypeSelect, prop.Input.Type)
		}
		if !prop.Input.Multiple {
			t.Error("SelectInput: expected multiple to be true")
		}
	}

	// Test new input types
	newInputTypes := map[string]InputType{
		"SwitchInput":   InputTypeSwitch,
		"DateInput":     InputTypeExpression,
		"ListInput":     InputTypeList,
		"JsonInput":     InputTypeJson,
		"KeyValueInput": InputTypeKeyValue,
	}

	for key, expectedType := range newInputTypes {
		if prop, exists := propMap[key]; exists && prop.Input != nil {
			if prop.Input.Type != expectedType {
				t.Errorf("%s: expected input type %v, got %v", key, expectedType, prop.Input.Type)
			}
		} else {
			t.Errorf("%s: property not found or input config is nil", key)
		}
	}
}

func TestMarshalProp_ErrorConditions(t *testing.T) {
	// Test nil input
	_, err := MarshalProp(nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	// Test non-pointer input
	_, err = MarshalProp(TestStruct{})
	if err == nil {
		t.Error("Expected error for non-pointer input")
	}

	// Test non-struct input
	str := "test"
	_, err = MarshalProp(&str)
	if err == nil {
		t.Error("Expected error for non-struct input")
	}
}

func TestMarshalProp_SkippedFields(t *testing.T) {
	type SkipTestStruct struct {
		IncludedField string `prop:"included" description:"This should be included"`
		SkippedField1 string `prop:"-" description:"This should be skipped"`
		SkippedField2 string `description:"This should be skipped - no prop tag"`
		unexported    string `prop:"unexported" description:"This should be skipped - unexported"`
	}

	testStruct := &SkipTestStruct{
		IncludedField: "included",
		SkippedField1: "skipped1",
		SkippedField2: "skipped2",
		unexported:    "unexported",
	}

	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(props) != 1 {
		t.Errorf("Expected 1 property, got %d", len(props))
	}

	if props[0].Key != "IncludedField" {
		t.Errorf("Expected property key 'IncludedField', got '%s'", props[0].Key)
	}
}

func TestMarshalProp_AllPropertyTypes(t *testing.T) {
	type AllTypesStruct struct {
		BoolField      bool           `prop:"bool_field"`
		IntField       int            `prop:"int_field"`
		Int64Field     int64          `prop:"int64_field"`
		Float32Field   float32        `prop:"float32_field"`
		Float64Field   float64        `prop:"float64_field"`
		StringField    string         `prop:"string_field"`
		SliceField     []string       `prop:"slice_field"`
		MapField       map[string]any `prop:"map_field"`
		InterfaceField any            `prop:"interface_field"`
	}

	testStruct := &AllTypesStruct{}
	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedTypes := map[string]PropertyType{
		"BoolField":      Bool,
		"IntField":       Int64,
		"Int64Field":     Int64,
		"Float32Field":   Float64,
		"Float64Field":   Float64,
		"StringField":    String,
		"SliceField":     Array,
		"MapField":       Map,
		"InterfaceField": JSON,
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	for key, expectedType := range expectedTypes {
		if prop, exists := propMap[key]; exists {
			if prop.Type != expectedType {
				t.Errorf("Property %s: expected type %s, got %s", key, expectedType, prop.Type)
			}
		} else {
			t.Errorf("Property %s not found", key)
		}
	}
}

func TestMarshalProp_FallbackInputTypes(t *testing.T) {
	type FallbackTestStruct struct {
		StringField string            `prop:"string_field"` // Should fallback to InputTypeText
		BoolField   bool              `prop:"bool_field"`   // Should fallback to InputTypeSwitch
		IntField    int               `prop:"int_field"`    // Should fallback to InputTypeNumber
		FloatField  float64           `prop:"float_field"`  // Should fallback to InputTypeNumber
		ArrayField  []string          `prop:"array_field"`  // Should fallback to InputTypeList
		MapField    map[string]string `prop:"map_field"`    // Should fallback to InputTypeKeyValue
		JsonField   any               `prop:"json_field"`   // Should fallback to InputTypeJson
	}

	testStruct := &FallbackTestStruct{}
	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedFallbacks := map[string]InputType{
		"StringField": InputTypeText,
		"BoolField":   InputTypeSwitch,
		"IntField":    InputTypeNumber,
		"FloatField":  InputTypeNumber,
		"ArrayField":  InputTypeList,
		"MapField":    InputTypeKeyValue,
		"JsonField":   InputTypeJson,
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	for key, expectedType := range expectedFallbacks {
		if prop, exists := propMap[key]; exists {
			if prop.Input == nil {
				t.Errorf("Property %s: expected fallback input config, got nil", key)
				continue
			}
			if prop.Input.Type != expectedType {
				t.Errorf("Property %s: expected fallback input type %v, got %v", key, expectedType, prop.Input.Type)
			}
		} else {
			t.Errorf("Property %s not found", key)
		}
	}
}

func TestMarshalProp_ValType(t *testing.T) {
	type ValTypeTestStruct struct {
		StringArray []string          `prop:"string_array"`
		IntArray    []int             `prop:"int_array"`
		StringMap   map[string]string `prop:"string_map"`
		IntMap      map[string]int    `prop:"int_map"`
		AnyMap      map[string]any    `prop:"any_map"`
	}

	testStruct := &ValTypeTestStruct{}
	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedValTypes := map[string]PropertyType{
		"StringArray": String,
		"IntArray":    Int64,
		"StringMap":   String,
		"IntMap":      Int64,
		"AnyMap":      JSON,
	}

	propMap := make(map[string]Property)
	for _, prop := range props {
		propMap[prop.Key] = prop
	}

	for key, expectedValType := range expectedValTypes {
		if prop, exists := propMap[key]; exists {
			if prop.ValType != expectedValType {
				t.Errorf("Property %s: expected valType %s, got %s", key, expectedValType, prop.ValType)
			}
		} else {
			t.Errorf("Property %s not found", key)
		}
	}
}

func TestMarshalProp_ValidationIntegration(t *testing.T) {
	// Test struct with incompatible input types
	type InvalidInputStruct struct {
		SwitchOnString string `prop:"switch_field" input:"switch" description:"This should fail"`
		ListOnString   string `prop:"list_field" input:"list" description:"This should fail"`
		JsonOnString   string `prop:"json_field" input:"json" description:"This should fail"`
		SelectOnArray  []int  `prop:"select_field" input:"select" description:"This should fail"`
	}

	testCases := []struct {
		name      string
		input     any
		expectErr bool
		errMsg    string
	}{
		{
			name: "switch input on string field",
			input: &struct {
				Field string `prop:"field" input:"switch"`
			}{},
			expectErr: true,
			errMsg:    "input type incompatible with property type",
		},
		{
			name: "list input on string field",
			input: &struct {
				Field string `prop:"field" input:"list"`
			}{},
			expectErr: true,
			errMsg:    "input type incompatible with property type",
		},
		{
			name: "json input on string field", 
			input: &struct {
				Field string `prop:"field" input:"json"`
			}{},
			expectErr: true,
			errMsg:    "input type incompatible with property type",
		},
		{
			name: "select input on array field",
			input: &struct {
				Field []string `prop:"field" input:"select"`
			}{},
			expectErr: true,
			errMsg:    "input type incompatible with property type",
		},
		{
			name: "multi-select input on string field",
			input: &struct {
				Field string `prop:"field" input:"multi-select"`
			}{},
			expectErr: true,
			errMsg:    "input type incompatible with property type",
		},
		{
			name: "valid switch input on bool field",
			input: &struct {
				Field bool `prop:"field" input:"switch"`
			}{},
			expectErr: false,
		},
		{
			name: "valid list input on array field",
			input: &struct {
				Field []string `prop:"field" input:"list"`
			}{},
			expectErr: false,
		},
		{
			name: "valid expression input on any field",
			input: &struct {
				Field int `prop:"field" input:"expression"`
			}{},
			expectErr: false,
		},
		{
			name: "invalid range input with missing parameters",
			input: &struct {
				Field float64 `prop:"field" input:"range"`
			}{},
			expectErr: true,
			errMsg:    "invalid input config",
		},
		{
			name: "valid range input with parameters",
			input: &struct {
				Field float64 `prop:"field" input:"range,min=0,max=100,step=1"`
			}{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := MarshalProp(tc.input)
			
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tc.errMsg != "" {
					// Check if error message contains expected substring
					found := false
					if len(tc.errMsg) <= len(err.Error()) {
						for i := 0; i <= len(err.Error())-len(tc.errMsg); i++ {
							if err.Error()[i:i+len(tc.errMsg)] == tc.errMsg {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("Expected error message to contain %q, got %q", tc.errMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestMarshalProp_FallbackValidation(t *testing.T) {
	// Test that fallback input types are always valid
	type FallbackTestStruct struct {
		StringField string  `prop:"string_field"`
		BoolField   bool    `prop:"bool_field"`
		IntField    int     `prop:"int_field"`
		FloatField  float64 `prop:"float_field"`
		ArrayField  []int   `prop:"array_field"`
		MapField    map[string]string `prop:"map_field"`
		JsonField   any     `prop:"json_field"`
	}

	testStruct := &FallbackTestStruct{}
	props, err := MarshalProp(testStruct)
	if err != nil {
		t.Fatalf("Expected no error with fallback input types, got: %v", err)
	}

	// Verify all properties were created successfully
	expectedCount := 7
	if len(props) != expectedCount {
		t.Errorf("Expected %d properties with fallback inputs, got %d", expectedCount, len(props))
	}

	// Verify each property has a valid input configuration
	for _, prop := range props {
		if prop.Input == nil {
			t.Errorf("Property %s missing input configuration", prop.Key)
			continue
		}
		
		// All fallback inputs should validate successfully
		if err := prop.Input.Validate(); err != nil {
			t.Errorf("Fallback input validation failed for %s: %v", prop.Key, err)
		}
		
		if err := prop.Input.ValidateType(prop.Type); err != nil {
			t.Errorf("Fallback input type validation failed for %s: %v", prop.Key, err)
		}
	}
}