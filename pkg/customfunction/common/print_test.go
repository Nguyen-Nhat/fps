package customFunc

import "testing"

func TestPrintFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		params   []string
		expected string
	}{
		{
			name:     "Simple string",
			format:   "Hello, %s!",
			params:   []string{"Alice"},
			expected: "Hello, Alice!",
		},
		{
			name:     "Multiple placeholders",
			format:   "%s has %s apples.",
			params:   []string{"Bob", "5"},
			expected: "Bob has 5 apples.",
		},
		{
			name:     "No placeholders",
			format:   "Just a plain string",
			params:   []string{},
			expected: "Just a plain string",
		},
		{
			name:     "Empty params with placeholders",
			format:   "%s, %s!",
			params:   []string{"", ""},
			expected: ", !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrintFormat(tt.format, tt.params)
			if result.Result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Result)
			}
		})
	}
}
