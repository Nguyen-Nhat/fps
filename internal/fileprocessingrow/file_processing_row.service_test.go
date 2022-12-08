package fileprocessingrow

import "testing"

func Test_isSuccess(t *testing.T) {
	type args struct {
		totalSuccess int
		totalFailed  int
		total        int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "", args: args{0, 0, 2}, want: false},
		{name: "", args: args{1, 0, 2}, want: false},
		{name: "", args: args{1, 1, 2}, want: true},
		{name: "", args: args{0, 1, 2}, want: false},
		{name: "", args: args{0, 2, 2}, want: true},
		{name: "", args: args{2, 0, 2}, want: true},
		{name: "", args: args{2, 1, 2}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFinished(tt.args.totalSuccess, tt.args.totalFailed, tt.args.total); got != tt.want {
				t.Errorf("isSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
