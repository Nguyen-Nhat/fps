package executerowgroup

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

const json1 = `
	{
	  "items": [
		{
		  "quantity": 7,
		  "sku": "13017"
		}
	  ],
	  "site": "7",
	  "type": "S42"
	}
	`

const json2 = `
	{
	  "items": [
		{
		  "quantity": 8,
		  "sku": "13018"
		}
	  ],
	  "site": "7",
	  "type": "S42"
	}
	`

const json1And2 = `
	{
	  "items": [
		{
		  "quantity": 7,
		  "sku": "13017"
		},
		{
		  "quantity": 8,
		  "sku": "13018"
		}
	  ],
	  "site": "7",
	  "type": "S42"
	}
	`

func Test_mergeMapInterface(t *testing.T) {
	var input1, input2, input1And2 map[string]interface{}
	_ = json.Unmarshal([]byte(json1), &input1)
	_ = json.Unmarshal([]byte(json2), &input2)
	_ = json.Unmarshal([]byte(json1And2), &input1And2)

	type args struct {
		first  map[string]interface{}
		second map[string]interface{}
	}
	tests := []struct {
		name            string
		args            args
		want            map[string]interface{}
		wantErr         bool
		wantErrContains string
	}{
		// case CAN merge -----------
		{"test mergeMapInterface case 2 maps are empty -> return TRUE and empty map",
			args{map[string]interface{}{}, map[string]interface{}{}},
			map[string]interface{}{}, false, ""},

		{"test mergeMapInterface case 1 maps are empty -> return TRUE and correct map",
			args{map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2}}, map[string]interface{}{}},
			map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2}}, false, ""},

		{"test mergeMapInterface case 1 maps are empty -> return TRUE and correct map",
			args{map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2}},
				map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{2, 3}}},
			map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2, 2, 3}}, false, ""},

		{"test mergeMapInterface case HAPPY with array object -> return TRUE and correct map",
			args{input1, input2},
			input1And2, false, ""},

		// case CAN NOT merge -----------
		{"test mergeMapInterface case field int not match -> return FALSE",
			args{map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2}},
				map[string]interface{}{"type_int": 2, "type_string": "abc", "type_arr": []int{2, 3}}},
			nil, true, "type_int"},

		{"test mergeMapInterface case field int not match -> return FALSE",
			args{map[string]interface{}{"type_int": 1, "type_string": "abc", "type_arr": []int{1, 2}},
				map[string]interface{}{"type_int": 1, "type_string": "abcd", "type_arr": []int{2, 3}}},
			nil, true, "type_string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotError := mergeMapInterface(tt.args.first, tt.args.second)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeMapInterface() got = %v, want %v", got, tt.want)
			}
			if (gotError != nil && !tt.wantErr) || (gotError == nil && tt.wantErr) {
				t.Errorf("mergeMapInterface() gotError = %v, wantErr %v", gotError, tt.wantErr)
			}
			if gotError != nil && !strings.Contains(gotError.Error(), tt.wantErrContains) {
				t.Errorf("mergeMapInterface() gotError = %v, wantErrContains %v", gotError, tt.wantErrContains)
			}
		})
	}
}
