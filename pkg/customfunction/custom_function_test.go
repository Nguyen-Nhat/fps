package customFunc

import (
	"reflect"
	"testing"
)

func TestIsCustomFunction(t *testing.T) {
	tests := []struct {
		name            string
		functionPattern string
		want            bool
	}{
		// case empty
		{"functionPattern is empty -> false", "", false},
		// case invalid
		{"functionPattern is not function case 1 -> false", "$", false},
		{"functionPattern is not function case 2 -> false", "$func", false},
		{"functionPattern is not function case 3 -> false", "$funcc", false},
		{"functionPattern is not function case 3 -> false", "$funcc.abc", false},
		// case valid
		{"functionPattern is function 1 -> false", "$func.abc", true},
		{"functionPattern is function 2 -> false", "$func.randomInt", true},
		{"functionPattern is function 3 -> false", "$func.randomInt()", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCustomFunction(tt.functionPattern); got != tt.want {
				t.Errorf("IsCustomFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCustomFunction(t *testing.T) {
	tests := []struct {
		name            string
		functionPattern string
		wantErr         bool
		want            CustomFunction
	}{
		// not function
		{"test ToCustomFunction: functionPattern is empty", "",
			true, CustomFunction{}},
		{"test ToCustomFunction: functionPattern is number", "1234",
			true, CustomFunction{}},
		{"test ToCustomFunction: functionPattern is string", "abcd",
			true, CustomFunction{}},
		{"test ToCustomFunction: functionPattern is column value", "$A",
			true, CustomFunction{}},

		// function wrong pattern
		{"test ToCustomFunction: functionPattern is wrong pattern", "$funcc.randomInt",
			true, CustomFunction{}},

		// function is correct
		{"test ToCustomFunction: functionPattern is correct, with no param", "$func.randomInt",
			false, CustomFunction{FunctionPattern: "$func.randomInt", Name: "randomInt", ParamsRaw: []string{}}},
		{"test ToCustomFunction: functionPattern is correct, with 1 param number", "$func.randomInt;123",
			false, CustomFunction{FunctionPattern: "$func.randomInt;123", Name: "randomInt", ParamsRaw: []string{"123"}}},
		{"test ToCustomFunction: functionPattern is correct, with 2 params number and string", "$func.randomInt;123;abce",
			false, CustomFunction{FunctionPattern: "$func.randomInt;123;abce", Name: "randomInt", ParamsRaw: []string{"123", "abce"}}},
		{"test ToCustomFunction: functionPattern is correct, with param is column value", "$func.randomInt;{{$A}}",
			false, CustomFunction{FunctionPattern: "$func.randomInt;{{$A}}", Name: "randomInt", ParamsRaw: []string{"{{$A}}"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToCustomFunction(tt.functionPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToCustomFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToCustomFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecuteFunction(t *testing.T) {
	tests := []struct {
		name    string
		cusFunc CustomFunction
		wantErr bool
		want    FuncResult
	}{
		// not match function name
		{"test ExecuteFunction: case not match function", CustomFunction{Name: "this is function that not existed"}, true, FuncResult{}},

		// funcTestError ...
		{"test ExecuteFunction: case funcTestError", CustomFunction{Name: funcTestError},
			false, FuncResult{nil, "this is testing error function"}},

		// funcTest ...
		{"test ExecuteFunction: case funcTest no params -> error",
			CustomFunction{Name: funcTest, ParamsMapped: []string{}},
			true, FuncResult{}},
		{"test ExecuteFunction: case funcTest has 1 param",
			CustomFunction{Name: funcTest, ParamsMapped: []string{"1"}},
			true, FuncResult{}},
		{"test ExecuteFunction: case funcTest has 2 params but wrong type",
			CustomFunction{Name: funcTest, ParamsMapped: []string{"1", "a"}},
			true, FuncResult{}},
		{"test ExecuteFunction: case funcTest has 2 params and all is valid",
			CustomFunction{Name: funcTest, ParamsMapped: []string{"1", "20"}},
			false, FuncResult{21, ""}},
		{"test ExecuteFunction: case funcTest has 3 params and all is valid",
			CustomFunction{Name: funcTest, ParamsMapped: []string{"1", "2", "a"}},
			false, FuncResult{3, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteFunction(tt.cusFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
