package flatten

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func Test_convertToRealValue(t *testing.T) {
	type args struct {
		fieldType    string
		valueStr     string
		dependsOnKey string
	}
	tests := []struct {
		name       string
		args       args
		wantValue  interface{}
		wantErrMsg string
	}{
		// Type String .....
		{"test type STRING",
			args{"STRING", "abc", "key_name"},
			"abc", ""},
		{"test type string",
			args{configloader.TypeString, "abcd", "key_name"},
			"abcd", ""},
		{"test type string empty",
			args{configloader.TypeString, "", "key_name"},
			"", ""},

		// Type integer ...
		{"test type inteGer, valid input",
			args{"inteGer", "1", "key_name"},
			int64(1), ""},
		{"test type INTEGER, valid input",
			args{"INTEGER", "2", "key_name"},
			int64(2), ""},
		{"test type integer, valid input",
			args{configloader.TypeInteger, "3", "key_name"},
			int64(3), ""},
		{"test type integer, valid input -10223",
			args{configloader.TypeInteger, "-10223", "key_name"},
			int64(-10223), ""},
		{"test type integer, valid input MAX_INT",
			args{configloader.TypeInteger, strconv.Itoa(math.MaxInt32), "key_name"},
			int64(math.MaxInt32), ""},
		{"test type integer, valid input MIN_INT",
			args{configloader.TypeInteger, strconv.Itoa(math.MinInt32), "key_name"},
			int64(math.MinInt32), ""},
		{"test type integer, invalid input",
			args{configloader.TypeInteger, "112sa", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, invalid input 1.0",
			args{configloader.TypeInteger, "1.0", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, invalid input 100000000.99999999",
			args{configloader.TypeInteger, "100000000.99999999", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, valid input MAX_LONG",
			args{configloader.TypeInteger, strconv.Itoa(math.MaxInt64), "key_name"},
			int64(math.MaxInt64), ""},
		{"test type integer, valid input MIN_LONG",
			args{configloader.TypeInteger, strconv.Itoa(math.MinInt64), "key_name"},
			int64(math.MinInt64), ""},
		{"test type integer empty",
			args{configloader.TypeInteger, "", "key_name"},
			int64(0), ""},

		// Type number .....
		{"test type numbEr, valid input",
			args{"numbEr", "0.3", "key_name"},
			0.3, ""},
		{"test type NUMBER, valid input",
			args{"NUMBER", "0.2", "key_name"},
			0.2, ""},
		{"test type number, valid input",
			args{configloader.TypeNumber, "0.1", "key_name"},
			0.1, ""},
		{"test type number, valid input 1.0",
			args{configloader.TypeNumber, "1.0", "key_name"},
			1.0, ""},
		{"test type number, valid input many 0000",
			args{configloader.TypeNumber, "10000.0000001", "key_name"},
			10000.0000001, ""},
		{"test type number, valid input MAX_DOUBLE",
			args{configloader.TypeNumber, fmt.Sprintf("%f", math.MaxFloat64), "key_name"},
			math.MaxFloat64, ""},
		{"test type number, invalid input",
			args{configloader.TypeNumber, "11.2sa", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type number empty",
			args{configloader.TypeNumber, "", "key_name"},
			float64(0), ""},

		// Type bool .....
		{"test type booleAN, valid input",
			args{"booleAN", "true", "key_name"},
			true, ""},
		{"test type boolean, valid input",
			args{configloader.TypeBoolean, "true", "key_name"},
			true, ""},
		{"test type boolean, valid input",
			args{configloader.TypeBoolean, "false", "key_name"},
			false, ""},
		{"test type boolean, invalid input",
			args{configloader.TypeBoolean, "falsee", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type BOOLEAN, valid input",
			args{"BOOLEAN", "falSE", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type BOOLEAN empty",
			args{"BOOLEAN", "", "key_name"},
			false, ""},

		// Type json .....
		{"test type json, valid input",
			args{configloader.TypeJson, "[123,456]", "key_name"},
			[]interface{}{float64(123), float64(456)}, ""},
		{"test type json, valid input",
			args{configloader.TypeJson, "[123.321,0.0001]", "key_name"},
			[]interface{}{123.321, 0.0001}, ""},
		{"test type json, valid input",
			args{configloader.TypeJson, "[\"abc\",\"cde\"]", "key_name"},
			[]interface{}{"abc", "cde"}, ""},
		{"test type json empty",
			args{configloader.TypeJson, "", "key_name"},
			nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := convertToRealValue(tt.args.fieldType, tt.args.valueStr, tt.args.dependsOnKey)
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("convertToRealValue() got = %v, want %v", got, tt.wantValue)
			}
			if got1 != tt.wantErrMsg {
				t.Errorf("convertToRealValue() got1 = %v, want %v", got1, tt.wantErrMsg)
			}
		})
	}
}
