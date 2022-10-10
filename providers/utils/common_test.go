package utils

import "testing"

func TestHiddenString(t *testing.T) {
	type args struct {
		input            string
		numberOfTailChar int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "input is empty",
			args: args{"", 3},
			want: "",
		},
		{
			name: "numberOfTailChar is zero",
			args: args{"abc", 0},
			want: "***",
		},
		{
			name: "input shorter than tail",
			args: args{"ab", 3},
			want: "ab",
		},
		{
			name: "input longer than tail",
			args: args{"abcd", 3},
			want: "***bcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HiddenString(tt.args.input, tt.args.numberOfTailChar); got != tt.want {
				t.Errorf("HiddenString() = %v, want %v", got, tt.want)
			}
		})
	}
}
