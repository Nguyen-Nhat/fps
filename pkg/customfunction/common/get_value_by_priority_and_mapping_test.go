package customFunc

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
)

func TestGetValueByPriority(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want FuncResult
	}{
		{
			name: "Test case 1 - type string",
			args: args{
				values: []string{"", "string", "a", "b", "c"},
			},
			want: FuncResult{Result: "a"},
		},
		{
			name: "Test case 2 - type integer",
			args: args{
				values: []string{"", "integer", "1", "b", "c"},
			},
			want: FuncResult{Result: int64(1)},
		},
		{
			name: "Test case 3 - type string has empty",
			args: args{
				values: []string{"", "string", "", "b", "c"},
			},
			want: FuncResult{Result: "b"},
		},
		{
			name: "Test case 4 - type boolean",
			args: args{
				values: []string{"", "boolean", "yEs", "N", "Y"},
			},
			want: FuncResult{Result: true},
		},
		{
			name: "Test case 5 - cant parse value",
			args: args{
				values: []string{"", "integer", "yEs", "N", "Y"},
			},
			want: FuncResult{ErrorMessage: errorz.ErrCantParseValue("yEs", "integer")},
		},
		{
			name: "Test case 6 - has dictionary",
			args: args{
				values: []string{"{\"a\": \"res_a\",\"b\": \"res_b\"}", "string", "a", "b", "c"},
			},
			want: FuncResult{Result: "res_a"},
		},
		{
			name: "Test case 7 - has dictionary but not in",
			args: args{
				values: []string{"{\"b\": \"res_b\"}", "string", "a", "b", "c"},
			},
			want: FuncResult{ErrorMessage: errorz.ErrNotExistValueInList("a", []string{"b"}).Error()},
		},
		{
			name: "Test case 8 - dictionary is invalid",
			args: args{
				values: []string{"{\"b\": \"res_b\"", "string", "a", "b", "c"},
			},
			want: FuncResult{ErrorMessage: "unexpected end of JSON input"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValueByPriorityAndMapping(tt.args.values[0], tt.args.values[1], tt.args.values[2:]); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValueByPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		valueStr string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test true", args{"true"}, true, false},
		{"test TRUE", args{"TRUE"}, true, false},
		{"test T", args{"T"}, true, false},
		{"test Y", args{"Y"}, true, false},
		{"test YeS", args{"YeS"}, true, false},
		{"test 1", args{"1"}, true, false},
		{"test false", args{"false"}, false, false},
		{"test FALSE", args{"FALSE"}, false, false},
		{"test F", args{"F"}, false, false},
		{"test N", args{"N"}, false, false},
		{"test NO", args{"NO"}, false, false},
		{"test 0", args{"0"}, false, false},
		{"test invalid", args{"invalid"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBool(tt.args.valueStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}
