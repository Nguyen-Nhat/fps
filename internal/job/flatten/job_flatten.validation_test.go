package flatten

import (
	"reflect"
	"testing"
)

func Test_validateAndMatchJsonPath(t *testing.T) {
	type args struct {
		rowID    int
		rowData  []string
		jsonPath string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantError *ErrorRow
	}{
		// case not contains variable
		{"test case path is object and not contains variable", args{1, []string{"a", "b", "c"}, `data.transaction.id`},
			`data.transaction.id`, nil},
		{"test case path is array and not contains variable", args{1, []string{"a", "b", "c"}, `data.transactions.#(name=="quy")`},
			`data.transactions.#(name=="quy")`, nil},

		// case contains variable
		{"test case path is object and contains variable", args{1, []string{"a", "id", "c"}, `data.transaction.{{ $B }}`},
			`data.transaction.id`, nil},
		{"test case path is array and contains variable", args{1, []string{"a", "b", "quy"}, `data.transactions.#(name=="{{ $C }}").id`},
			`data.transactions.#(name=="quy").id`, nil},

		// case contains multiple variables
		{"test case path is object and contains multiple variables -> no support all", args{1, []string{"a", "id", "member"}, `data.transaction.{{ $C }}.{{ $B }}`},
			`data.transaction.member.{{ $B }}`, nil},
		{"test case path is array and contains multiple variables -> no support all", args{1, []string{"id", "b", "quy"}, `data.transactions.#(name=="{{ $C }}").{{ $A }}`},
			`data.transactions.#(name=="quy").{{ $A }}`, nil},

		// case contains space in variable
		{"test case path is array and contains variable", args{1, []string{"a", "b", "quy"}, `data.transactions.#(name=="{{               $C                 }}").id`},
			`data.transactions.#(name=="quy").id`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateAndMatchJsonPath(tt.args.rowID, tt.args.rowData, tt.args.jsonPath)
			if got != tt.want {
				t.Errorf("validateAndMatchJsonPath() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantError) {
				t.Errorf("validateAndMatchJsonPath() got1 = %v, want %v", got1, tt.wantError)
			}
		})
	}
}

func Test_mapValueForCustomFunctionParams(t *testing.T) {
	rowData := []string{"value A", "value B", "value C"}
	fileParameters := map[string]interface{}{
		"field_num":   12,
		"field_num_2": 23,
		"field_str":   "abc_def",
		"field_str_2": "abc_3209",
	}

	type args struct {
		paramsRaw      []string
		rowData        []string
		fileParameters map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test mapValueForCustomFunctionParams: case raw params is empty",
			args{[]string{}, rowData, fileParameters},
			[]string{}},
		{"test mapValueForCustomFunctionParams: case raw params contains number, string",
			args{[]string{"1", "abc"}, rowData, fileParameters},
			[]string{"1", "abc"}},
		{"test mapValueForCustomFunctionParams: case raw params contains $A, $param",
			args{[]string{"1", "abc", "$A", "$C", "$param.field_num", "$param.field_str"}, rowData, fileParameters},
			[]string{"1", "abc", "value A", "value C", "12", "abc_def"}},
		{"test mapValueForCustomFunctionParams: case raw params contains pattern not support are $response, $header",
			args{[]string{"1", "abc", "$response1.data.id", "$header.A"}, rowData, fileParameters},
			[]string{"1", "abc", "$response1.data.id", "$header.A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapValueForCustomFunctionParams(tt.args.paramsRaw, tt.args.rowData, tt.args.fileParameters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapValueForCustomFunctionParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
