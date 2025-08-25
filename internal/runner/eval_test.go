package runner

import (
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		states     map[string]any
		ref        string
		want       any
		wantErr    bool
	}{
		{
			name:       "simple arithmetic",
			expression: "1 + 2",
			states:     map[string]any{},
			ref:        "test",
			want:       3,
			wantErr:    false,
		},
		{
			name:       "string concatenation",
			expression: `"hello" + " " + "world"`,
			states:     map[string]any{},
			ref:        "test",
			want:       "hello world",
			wantErr:    false,
		},
		{
			name:       "access node output",
			expression: "node1.output.value",
			states: map[string]any{
				"node1.output.value": 42,
			},
			ref:     "test",
			want:    42,
			wantErr: false,
		},
		{
			name:       "access node input",
			expression: "node1.input.text",
			states: map[string]any{
				"node1.input.text": "test string",
			},
			ref:     "test",
			want:    "test string",
			wantErr: false,
		},
		{
			name:       "nested property access",
			expression: "node1.output.data",
			states: map[string]any{
				"node1.output.data": map[string]any{
					"key": "value",
				},
			},
			ref:     "test",
			want:    map[string]any{"key": "value"},
			wantErr: false,
		},
		{
			name:       "len function with string",
			expression: `len("hello")`,
			states:     map[string]any{},
			ref:        "test",
			want:       5,
			wantErr:    false,
		},
		{
			name:       "len function with array from state",
			expression: "len(node1.output.items)",
			states: map[string]any{
				"node1.output.items": []any{1, 2, 3, 4},
			},
			ref:     "test",
			want:    4,
			wantErr: false,
		},
		{
			name:       "strContains function",
			expression: `strContains("hello world", "world")`,
			states:     map[string]any{},
			ref:        "test",
			want:       true,
			wantErr:    false,
		},
		{
			name:       "toLowerCase function",
			expression: `toLowerCase("HELLO")`,
			states:     map[string]any{},
			ref:        "test",
			want:       "hello",
			wantErr:    false,
		},
		{
			name:       "toUpperCase function",
			expression: `toUpperCase("hello")`,
			states:     map[string]any{},
			ref:        "test",
			want:       "HELLO",
			wantErr:    false,
		},
		{
			name:       "toString function",
			expression: "toString(42)",
			states:     map[string]any{},
			ref:        "test",
			want:       "42",
			wantErr:    false,
		},
		{
			name:       "comparison operators",
			expression: "node1.output.count > 10",
			states: map[string]any{
				"node1.output.count": 15,
			},
			ref:     "test",
			want:    true,
			wantErr: false,
		},
		{
			name:       "logical operators",
			expression: "node1.output.active && node2.output.ready",
			states: map[string]any{
				"node1.output.active": true,
				"node2.output.ready":  true,
			},
			ref:     "test",
			want:    true,
			wantErr: false,
		},
		{
			name:       "complex expression",
			expression: `len(node1.output.items) > 0 && strContains(node2.output.text, "success")`,
			states: map[string]any{
				"node1.output.items": []any{1, 2, 3},
				"node2.output.text":  "operation success",
			},
			ref:     "test",
			want:    true,
			wantErr: false,
		},
		{
			name:       "array indexing",
			expression: "node1.output.items[0]",
			states: map[string]any{
				"node1.output.items": []any{"first", "second", "third"},
			},
			ref:     "test",
			want:    "first",
			wantErr: false,
		},
		{
			name:       "map property access",
			expression: `node1.output.config["key"]`,
			states: map[string]any{
				"node1.output.config": map[string]any{
					"key": "value",
				},
			},
			ref:     "test",
			want:    "value",
			wantErr: false,
		},
		{
			name:       "invalid expression",
			expression: "invalid syntax [[",
			states:     map[string]any{},
			ref:        "test",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "undefined variable",
			expression: "undefinedVar",
			states:     map[string]any{},
			ref:        "test",
			want:       nil,
			wantErr:    false, // expr-lang returns nil for undefined variables
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create WorkflowState and populate it with test data
			workflowState := NewWorkflowState()
			for key, value := range tt.states {
				workflowState.Set(key, value)
			}
			got, err := Eval(tt.expression, workflowState)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() = %v (type %T), want %v (type %T)", got, got, tt.want, tt.want)
			}
		})
	}
}
