package customFunc

import "testing"

func TestIsContain(t *testing.T) {
	tests := []struct {
		name              string
		str               string
		subStr            string
		isCaseInsensitive string
		mustEqual         string
		expectedResult    bool
	}{
		// mustEqual = true
		{"Equal strings, case sensitive", "Hello", "Hello", "false", "true", true},
		{"Equal strings, case insensitive", "Hello", "hello", "true", "true", true},
		{"Not equal strings", "Hello", "World", "false", "true", false},

		// isCaseInsensitive = "true"
		{"SubStr contained, case insensitive", "Hello World", "world", "true", "false", true},
		{"SubStr not contained, case insensitive", "Hello World", "planet", "true", "false", false},

		// isCaseInsensitive = "false"
		{"SubStr contained, case sensitive", "Hello World", "World", "false", "false", true},
		{"SubStr not contained, case sensitive", "Hello World", "world", "false", "false", false},

		// Edge cases
		{"Empty str and subStr, mustEqual", "", "", "false", "true", true},
		{"Empty str, non-empty subStr", "", "Hello", "false", "false", false},
		{"Non-empty str, empty subStr", "Hello World", "", "false", "false", true}, // contains always returns "true" for empty subStr
		{"Empty str and subStr, case insensitive", "", "", "true", "false", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsContain(tt.str, tt.subStr, tt.isCaseInsensitive, tt.mustEqual)
			if result.Result != tt.expectedResult {
				t.Errorf("IsContain(%q, %q, %v, %v) = %v; want %v",
					tt.str, tt.subStr, tt.isCaseInsensitive, tt.mustEqual, result.Result, tt.expectedResult)
			}
		})
	}
}
