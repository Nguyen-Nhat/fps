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
			args{configloader.TypeString, "abc", "key_name"},
			"abc", ""},
		{"test type string",
			args{"string", "abcd", "key_name"},
			"abcd", ""},

		// Type INT ...
		{"test type inT, valid input",
			args{"inT", "1", "key_name"},
			int64(1), ""},
		{"test type int, valid input",
			args{"int", "2", "key_name"},
			int64(2), ""},
		{"test type INT, valid input",
			args{configloader.TypeInt, "3", "key_name"},
			int64(3), ""},
		{"test type INT, valid input -10223",
			args{configloader.TypeInt, "-10223", "key_name"},
			int64(-10223), ""},
		{"test type INT, valid input MAX_INT",
			args{configloader.TypeInt, strconv.Itoa(math.MaxInt32), "key_name"},
			int64(math.MaxInt32), ""},
		{"test type INT, valid input MIN_INT",
			args{configloader.TypeInt, strconv.Itoa(math.MinInt32), "key_name"},
			int64(math.MinInt32), ""},
		{"test type INT, invalid input",
			args{configloader.TypeInt, "112sa", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type INT, invalid input 1.0",
			args{configloader.TypeInt, "1.0", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
		{"test type INT, invalid input 100000000.99999999",
			args{configloader.TypeInt, "100000000.99999999", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},

		// LONG ....
		{"test type LonG, valid input",
			args{"LonG", "1", "key_name"},
			int64(1), ""},
		{"test type long, valid input",
			args{"long", "2", "key_name"},
			int64(2), ""},
		{"test type LONG, valid input",
			args{configloader.TypeLong, "3242", "key_name"},
			int64(3242), ""},
		{"test type LONG, valid input MAX_LONG",
			args{configloader.TypeLong, strconv.Itoa(math.MaxInt64), "key_name"},
			int64(math.MaxInt64), ""},
		{"test type LONG, valid input MIN_LONG",
			args{configloader.TypeLong, strconv.Itoa(math.MinInt64), "key_name"},
			int64(math.MinInt64), ""},

		// Type Double .....
		{"test type douBle, valid input",
			args{"douBle", "0.3", "key_name"},
			0.3, ""},
		{"test type double, valid input",
			args{"douBle", "0.2", "key_name"},
			0.2, ""},
		{"test type DOUBLE, valid input",
			args{configloader.TypeDouble, "0.1", "key_name"},
			0.1, ""},
		{"test type double, valid input 1.0",
			args{configloader.TypeDouble, "1.0", "key_name"},
			1.0, ""},
		{"test type double, valid input many 0000",
			args{configloader.TypeDouble, "10000.0000001", "key_name"},
			10000.0000001, ""},
		{"test type DOUBLE, valid input MAX_DOUBLE",
			args{configloader.TypeDouble, fmt.Sprintf("%f", math.MaxFloat64), "key_name"},
			math.MaxFloat64, ""},
		{"test type double, invalid input",
			args{configloader.TypeDouble, "11.2sa", "key_name"},
			nil, fmt.Sprintf("%s (%s)", errTypeWrong, "key_name")},
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
