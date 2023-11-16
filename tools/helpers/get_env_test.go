package helpers

import (
	"os"
	"testing"
)

func Test_Getenv(t *testing.T) {
	type args struct {
		key      string
		value    string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test Getenv when key not existed then return fallback",
			args{"", "", "this is fallback"}, "this is fallback"},
		{"test Getenv when key not existed then return correct value",
			args{"my_key", "this is value", "this is fallback"}, "this is value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv(tt.args.key, tt.args.value)
			if got := Getenv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("Getenv() = %v, want %v", got, tt.want)
			}
		})
	}
}
