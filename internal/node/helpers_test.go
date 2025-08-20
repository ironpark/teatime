package node

import "testing"

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"float64", float64(3.14), true},
		{"float32", float32(3.14), true},
		{"int", 42, true},
		{"int64", int64(42), true},
		{"int32", int32(42), true},
		{"uint", uint(42), true},
		{"uint64", uint64(42), true},
		{"string", "42", false},
		{"bool", true, false},
		{"nil", nil, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumber(tt.v); got != tt.want {
				t.Errorf("isNumber(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsString(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"string", "hello", true},
		{"empty string", "", true},
		{"int", 42, false},
		{"bool", true, false},
		{"nil", nil, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isString(tt.v); got != tt.want {
				t.Errorf("isString(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsBool(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"string", "true", false},
		{"int", 1, false},
		{"nil", nil, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBool(tt.v); got != tt.want {
				t.Errorf("isBool(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsWholeNumber(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"whole positive", 42.0, true},
		{"whole negative", -42.0, true},
		{"whole zero", 0.0, true},
		{"decimal", 42.5, false},
		{"small decimal", 42.001, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWholeNumber(tt.f); got != tt.want {
				t.Errorf("isWholeNumber(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsPositiveWholeNumber(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"positive whole", 42.0, true},
		{"zero", 0.0, true},
		{"negative whole", -42.0, false},
		{"positive decimal", 42.5, false},
		{"negative decimal", -42.5, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPositiveWholeNumber(tt.f); got != tt.want {
				t.Errorf("isPositiveWholeNumber(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}