package fileprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTotalErrorPlus(t *testing.T) {
	tests := []struct {
		name           string
		fileParameters string
		want           int32
	}{
		// return zero
		{"test getTotalErrorPlus case fileParameters is empty -> return 0", "", 0},
		{"test getTotalErrorPlus case fileParameters is json empty -> return 0", "{}", 0},
		{"test getTotalErrorPlus case fileParameters is not json -> return 0", "abc", 0},
		{"test getTotalErrorPlus case fileParameters is json, but not include totalErrorPlus -> return 0", `{"abc": "value1", "def": 123}`, 0},

		// return number
		{"test getTotalErrorPlus case fileParameters is json, and include totalErrorPlus 0 -> return 0",
			`{"abc": "value1", "def": 123, "totalErrorPlus": 0}`, 0},
		{"test getTotalErrorPlus case fileParameters is json, and include totalErrorPlus negative -> return 0",
			`{"abc": "value1", "def": 123, "totalErrorPlus": -12}`, 0},
		{"test getTotalErrorPlus case fileParameters is json, and include totalErrorPlus positive -> return correct",
			`{"abc": "value1", "def": 123, "totalErrorPlus": 32}`, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getTotalErrorPlus(tt.fileParameters), "getTotalErrorPlus(%v)", tt.fileParameters)
		})
	}
}
