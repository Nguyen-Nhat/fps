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
