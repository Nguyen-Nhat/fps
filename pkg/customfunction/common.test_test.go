package customFunc

import (
	"reflect"
	"testing"
)

func Test_testFunc(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want FuncResult
	}{
		{"test testFunc with 1 and 2", args{1, 2}, FuncResult{3, ""}},
		{"test testFunc with -1 and 2", args{-1, 2}, FuncResult{1, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testFunc(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testFuncError(t *testing.T) {
	tests := []struct {
		name string
		want FuncResult
	}{
		{"test testFuncError return default error message", FuncResult{nil, "this is testing error function"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testFuncError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testFuncError() = %v, want %v", got, tt.want)
			}
		})
	}
}
