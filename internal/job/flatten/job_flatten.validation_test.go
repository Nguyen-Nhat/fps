package flatten

import (
	"context"
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

func Test_validateAndMatchJsonPath(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rowID      int
		rowData    []string
		fileHeader []string
		jsonPath   string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantError *ErrorRow
	}{
		// case not contains variable
		{"test case path is object and not contains variable", args{1, []string{"a", "b", "c"}, []string{"Header A", "Header B", "Header C"}, `data.transaction.id`},
			`data.transaction.id`, nil},
		{"test case path is array and not contains variable", args{1, []string{"a", "b", "c"}, []string{"Header A", "Header B", "Header C"}, `data.transactions.#(name=="quy")`},
			`data.transactions.#(name=="quy")`, nil},

		// case contains variable
		{"test case path is object and contains variable", args{1, []string{"a", "id", "c"}, []string{"Header A", "Header B", "Header C"}, `data.transaction.{{ $B }}`},
			`data.transaction.id`, nil},
		{"test case path is array and contains variable", args{1, []string{"a", "b", "quy"}, []string{"Header A", "Header B", "Header C"}, `data.transactions.#(name=="{{ $C }}").id`},
			`data.transactions.#(name=="quy").id`, nil},

		// case contains multiple variables
		{"test case path is object and contains multiple variables -> no support all", args{1, []string{"a", "id", "member"}, []string{"Header A", "Header B", "Header C"}, `data.transaction.{{ $C }}.{{ $B }}`},
			`data.transaction.member.{{ $B }}`, nil},
		{"test case path is array and contains multiple variables -> no support all", args{1, []string{"id", "b", "quy"}, []string{"Header A", "Header B", "Header C"}, `data.transactions.#(name=="{{ $C }}").{{ $A }}`},
			`data.transactions.#(name=="quy").{{ $A }}`, nil},

		// case contains space in variable
		{"test case path is array and contains variable", args{1, []string{"a", "b", "quy"}, []string{"Header A", "Header B", "Header C"}, `data.transactions.#(name=="{{               $C                 }}").id`},
			`data.transactions.#(name=="quy").id`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateAndMatchJsonPath(ctx, tt.args.rowID, tt.args.rowData, tt.args.jsonPath, tt.args.fileHeader)
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

func Test_validateAndGetValueForRequestFieldExcel(t *testing.T) {
	ctx := context.Background()
	_, _ = i18n.LoadI18n("../../../resources/messages")
	reqFieldKeyARequired := configloader.RequestFieldMD{ValueDependsOnKey: "A", Required: true}
	reqFieldKeyANotRequired := configloader.RequestFieldMD{ValueDependsOnKey: "A", Required: false}
	reqFieldKeyABRequired := configloader.RequestFieldMD{ValueDependsOnKey: "AB", Required: true}
	reqFieldKeyABNotRequired := configloader.RequestFieldMD{ValueDependsOnKey: "AB", Required: true}
	reqFieldKeyWrong := configloader.RequestFieldMD{ValueDependsOnKey: "A1", Required: true}
	reqFieldKeyARequiredWithDefault := configloader.RequestFieldMD{ValueDependsOnKey: "A", Required: true, DefaultValuePattern: "valDefault"}
	reqFieldKeyANotRequiredWithDefault := configloader.RequestFieldMD{ValueDependsOnKey: "A", Required: false, DefaultValuePattern: "valDefault"}

	type args struct {
		rowID      int
		rowData    []string
		fileHeader []string
		reqField   *configloader.RequestFieldMD
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test case wrong column key -> error", args{1, []string{}, []string{"Header A", "Header B", "Header C"}, &reqFieldKeyWrong}, "", true},
		{"test case field require but key out of column range -> error", args{1, []string{"val1", "val2"}, []string{"Header A", "Header B", "Header C"}, &reqFieldKeyABRequired}, "", true},
		{"test case field require but value is empty -> error", args{1, []string{"", "val2"}, []string{"Header A", "Header B", "Header C"}, &reqFieldKeyARequired}, "", true},

		{"test case field A require and key in column range -> correct value", args{1, []string{"val1", "val2"}, []string{"Header A", "Header B", "Header C"}, &reqFieldKeyARequired}, "val1", false},
		{"test case field A NOT require, key in column range -> correct value", args{1, []string{"", "val2"}, []string{"Header A", "Header B", "Header C"}, &reqFieldKeyANotRequired}, "", false},

		{"test case field AB require, key in column range -> correct value", args{1, []string{
			"val1", "val2", "val3", "val4", "val5", "val6", "val7", "val8", "val9", "val10",
			"val1", "val2", "val3", "val4", "val5", "val6", "val7", "val8", "val9", "val10",
			"val1", "val2", "val3", "val4", "val5", "val_Z", "val_AA", "val_AB", "val9", "val10",
		}, []string{}, &reqFieldKeyABRequired}, "val_AB", false},
		{"test case field AB NOT require, key OUT of column range -> empty value", args{1, []string{"val1", "val2"}, []string{"Header A", "Header B"}, &reqFieldKeyABNotRequired}, "", false},

		{"test case field A require and key in column range but empty -> get default value", args{1, []string{"", "val2"}, []string{"Header A", "Header B"}, &reqFieldKeyARequiredWithDefault}, "valDefault", false},
		{"test case field A NOT require, key in column range but empty -> get default value", args{1, []string{"", "val2"}, []string{"Header A"}, &reqFieldKeyANotRequiredWithDefault}, "valDefault", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validateAndGetValueForRequestFieldExcel(ctx, tt.args.rowID, tt.args.rowData, tt.args.reqField, tt.args.fileHeader)
			if got != tt.want {
				t.Errorf("validateAndGetValueForRequestFieldExcel() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr && len(gotErr) == 0 {
				t.Errorf("validateAndGetValueForRequestFieldExcel() got1 = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_getValueDependOnKey(t *testing.T) {
	type args struct {
		valueDependsOn    configloader.ValueDependsOn
		valueDependsOnKey string
		fileHeader        []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test case valueDependsOn is EXCEL -> return header name", args{configloader.ValueDependsOnExcel, "A", []string{"Header A", "Header B", "Header C"}}, "Header A"},
		{"test case valueDependsOn is EXCEL but header empty -> return valueDependsOnKey", args{configloader.ValueDependsOnExcel, "D", []string{"Header A", "Header B", "Header C"}}, "D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueDependOnKeyExcel(tt.args.valueDependsOn, tt.args.valueDependsOnKey, tt.args.fileHeader); got != tt.want {
				t.Errorf("getValueDependOnKeyExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}
