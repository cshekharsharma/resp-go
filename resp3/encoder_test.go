package resp3

import (
	"errors"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Basic Types
		{
			name:     "SmallString as Simple String",
			input:    "hello",
			expected: "+hello\r\n",
		},
		{
			name:     "LongString as Bulk String",
			input:    "This is a long string of length > 16",
			expected: "$36\r\nThis is a long string of length > 16\r\n",
		},
		{
			name:     "Integer",
			input:    123,
			expected: ":123\r\n",
		},
		{
			name:     "Float",
			input:    3.14,
			expected: ",3.140000\r\n",
		},
		{
			name:     "Boolean True",
			input:    true,
			expected: "#t\r\n",
		},
		{
			name:     "Boolean False",
			input:    false,
			expected: "#f\r\n",
		},
		{
			name:     "Nil",
			input:    nil,
			expected: "_\r\n",
		},
		{
			name:     "Error",
			input:    errors.New("an error"),
			expected: "-an error\r\n",
		},

		// Arrays
		{
			name:     "Array of Integers",
			input:    []int{1, 2, 3},
			expected: "*3\r\n:1\r\n:2\r\n:3\r\n",
		},
		{
			name:     "Array of Strings",
			input:    []string{"a", "b", "c"},
			expected: "*3\r\n+a\r\n+b\r\n+c\r\n",
		},
		{
			name:     "Array of Mixed Types",
			input:    []interface{}{"a", 123, true},
			expected: "*3\r\n+a\r\n:123\r\n#t\r\n",
		},
		{
			name:     "Array of Booleans",
			input:    []bool{true, false, true},
			expected: "*3\r\n#t\r\n#f\r\n#t\r\n",
		},
		{
			name:     "Array of Floats",
			input:    []float64{1.23, 4.56, 7.89},
			expected: "*3\r\n,1.230000\r\n,4.560000\r\n,7.890000\r\n",
		},

		// Maps
		{
			name:     "Map with String Keys",
			input:    map[string]interface{}{"a": 1, "b": 2},
			expected: "%4\r\n+a\r\n:1\r\n+b\r\n:2\r\n",
		},
		{
			name:     "Map with Interface Keys",
			input:    map[interface{}]interface{}{"a": 1, 2: "b"},
			expected: "%4\r\n+a\r\n:1\r\n:2\r\n+b\r\n",
		},
		{
			name:     "Map with Mixed Keys",
			input:    map[interface{}]interface{}{"a": 1, 2: true, 3.14: "pi"},
			expected: "%6\r\n+a\r\n:1\r\n:2\r\n#t\r\n,3.140000\r\n+pi\r\n",
		},

		// Structs
		{
			name: "Struct",
			input: struct {
				Name string
				Post string
				Age  int
			}{"Alice", "Senior Software Engineer", 25},
			expected: "%6\r\n+Name\r\n+Alice\r\n+Post\r\n$24\r\nSenior Software Engineer\r\n+Age\r\n:25\r\n",
		},

		// Time
		{
			name:     "Time",
			input:    time.UnixMilli(1620832335000),
			expected: ":1620832335000\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("Encode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEncodeErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input chan int
	}{
		{
			name:  "Unsupported Type",
			input: make(chan int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encode(tt.input)
			if err == nil {
				t.Fatalf("expected error, got none")
			}
		})
	}
}
