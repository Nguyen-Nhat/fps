package taskprovider

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_convertToRealValue(t *testing.T) {
	type args struct {
		fieldType    string
		valueStr     string
		dependsOnKey string
	}
	tests := []struct {
		name      string
		args      args
		wantValue interface{}
		wantError error
	}{
		// Type String .....
		{"test type STRING",
			args{configloader.TypeString, "abc", "key_name"},
			"abc", nil},
		{"test type string",
			args{"string", "abcd", "key_name"},
			"abcd", nil},

		// Type integer ...
		{"test type inteGer, valid input",
			args{"inteGer", "1", "key_name"},
			int64(1), nil},
		{"test type INTEGER, valid input",
			args{"INTEGER", "2", "key_name"},
			int64(2), nil},
		{"test type integer, valid input",
			args{configloader.TypeInteger, "3", "key_name"},
			int64(3), nil},
		{"test type integer, valid input -10223",
			args{configloader.TypeInteger, "-10223", "key_name"},
			int64(-10223), nil},
		{"test type integer, valid input MAX_INT",
			args{configloader.TypeInteger, strconv.Itoa(math.MaxInt32), "key_name"},
			int64(math.MaxInt32), nil},
		{"test type integer, valid input MIN_INT",
			args{configloader.TypeInteger, strconv.Itoa(math.MinInt32), "key_name"},
			int64(math.MinInt32), nil},
		{"test type integer, invalid input",
			args{configloader.TypeInteger, "112sa", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, invalid input 1.0",
			args{configloader.TypeInteger, "1.0", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, invalid input 100000000.99999999",
			args{configloader.TypeInteger, "100000000.99999999", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},
		{"test type integer, valid input MAX_LONG",
			args{configloader.TypeInteger, strconv.Itoa(math.MaxInt64), "key_name"},
			int64(math.MaxInt64), nil},
		{"test type integer, valid input MIN_LONG",
			args{configloader.TypeInteger, strconv.Itoa(math.MinInt64), "key_name"},
			int64(math.MinInt64), nil},

		// Type number .....
		{"test type numbEr, valid input",
			args{"numbEr", "0.3", "key_name"},
			0.3, nil},
		{"test type NUMBER, valid input",
			args{"NUMBER", "0.2", "key_name"},
			0.2, nil},
		{"test type number, valid input",
			args{configloader.TypeNumber, "0.1", "key_name"},
			0.1, nil},
		{"test type number, valid input 1.0",
			args{configloader.TypeNumber, "1.0", "key_name"},
			1.0, nil},
		{"test type number, valid input many 0000",
			args{configloader.TypeNumber, "10000.0000001", "key_name"},
			10000.0000001, nil},
		{"test type number, valid input MAX_DOUBLE",
			args{configloader.TypeNumber, fmt.Sprintf("%f", math.MaxFloat64), "key_name"},
			math.MaxFloat64, nil},
		{"test type number, invalid input",
			args{configloader.TypeNumber, "11.2sa", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},

		// Type json .....
		{"test type booleAN, valid input",
			args{"booleAN", "true", "key_name"},
			true, nil},
		{"test type boolean, valid input",
			args{configloader.TypeBoolean, "true", "key_name"},
			true, nil},
		{"test type boolean, valid input",
			args{configloader.TypeBoolean, "false", "key_name"},
			false, nil},
		{"test type boolean, invalid input",
			args{configloader.TypeBoolean, "falsee", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},
		{"test type BOOLEAN, valid input",
			args{"BOOLEAN", "falSE", "key_name"},
			nil, fmt.Errorf("%s (%s)", errTypeWrong, "key_name")},

		// Type json .....
		{"test type json, valid input",
			args{configloader.TypeJson, "[123,456]", "key_name"},
			[]interface{}{float64(123), float64(456)}, nil},
		{"test type json, valid input",
			args{configloader.TypeJson, "[123.321,0.0001]", "key_name"},
			[]interface{}{123.321, 0.0001}, nil},
		{"test type json, valid input",
			args{configloader.TypeJson, "[\"abc\",\"cde\"]", "key_name"},
			[]interface{}{"abc", "cde"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRealValue, gotError := convertToRealValue(tt.args.fieldType, tt.args.valueStr, tt.args.dependsOnKey)
			if !reflect.DeepEqual(gotRealValue, tt.wantValue) {
				t.Errorf("convertToRealValue() gotRealValue = %v, want %v", gotRealValue, tt.wantValue)
			}
			if (gotError == nil && tt.wantError != nil) ||
				(gotError != nil && gotError.Error() != tt.wantError.Error()) {
				t.Errorf("convertToRealValue() gotError = %v, want %v", gotError, tt.wantError)
			}
		})
	}
}

func Test_isTaskSuccess(t *testing.T) {
	http200 := int32(200)

	type args struct {
		responseMD   configloader.ResponseMD
		httpStatus   int
		responseBody string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// check by http 200
		{"responseMD.httpsStatus=200, httpStatus=200, responseBody=empty --> true",
			args{configloader.ResponseMD{HttpStatusSuccess: &http200}, 200, ""},
			true},
		{"responseMD.httpsStatus=200, httpStatus=200, responseBody have value --> true",
			args{configloader.ResponseMD{HttpStatusSuccess: &http200}, 200, `{"code":"12"}`},
			true},
		{"responseMD.httpsStatus=200, httpStatus=0, responseBody have value --> false",
			args{configloader.ResponseMD{HttpStatusSuccess: &http200}, 0, `{"code":"12"}`},
			false},
		{"responseMD.httpsStatus=200, httpStatus=400, responseBody have value",
			args{configloader.ResponseMD{HttpStatusSuccess: &http200}, 400, `{"code":"12"}`},
			false},
		// check by code 00
		{"responseMD.Code is {code, 00}, httpStatus=200, responseBody=empty",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "00"}}, 200, ""},
			false},
		{"responseMD.Code is {code, 00}, httpStatus=200, responseBody={code:00} --> true",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "00"}}, 200, `{"code":"00"}`},
			true},
		// check by code 200
		{"responseMD.Code is {code, 200}, httpStatus=200, responseBody={code:200} --> true",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "200"}}, 200, `{"code":200}`},
			true},
		{"responseMD.Code is {code, 200}, httpStatus=200, responseBody={code:'200'} --> true",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "200"}}, 200, `{"code":"200"}`},
			true},
		// check by code 0
		{"responseMD.Code is {code, 0}, httpStatus=200, responseBody={code:0} --> true",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "0"}}, 200, `{"code":0}`},
			true},
		{"responseMD.Code is {code, 0}, httpStatus=200, responseBody={code:'0'} --> true",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "0"}}, 200, `{"code":"0"}`},
			true},
		// check by code error 400
		{"responseMD.Code is {code, 200}, httpStatus=400, responseBody={code:400} --> false",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "200"}}, 400, `{"code":400}`},
			false},
		{"responseMD.Code is {code, 200}, httpStatus=400, responseBody={code:'400'} --> false",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "200"}}, 400, `{"code":"400"}`},
			false},
		// check by code error 502
		{"responseMD.Code is {code, 200}, httpStatus=502, responseBody=empty --> false",
			args{configloader.ResponseMD{Code: configloader.ResponseCode{Path: "code", SuccessValues: "00"}}, 400, `{"code":400}`},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTaskSuccess(tt.args.responseMD, tt.args.httpStatus, tt.args.responseBody); got != tt.want {
				t.Errorf("isTaskSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getResponseMessage(t *testing.T) {
	type args struct {
		responseMsg  configloader.ResponseMsg
		httpStatus   int
		responseBody string
		isSuccess    bool
	}
	tests := []struct {
		name            string
		args            args
		wantMsgContains string
	}{
		{"http 200, responseBody=empty, success=true --> empty",
			args{configloader.ResponseMsg{Path: "message"}, 200, "", true},
			""},
		{"http 200, responseBody=empty, success=false --> default error",
			args{configloader.ResponseMsg{Path: "message"}, 200, "", false},
			errMsgByHttp},
		{"http 400, responseBody={message: 'invalid data'}, success=false --> invalid data",
			args{configloader.ResponseMsg{Path: "message"}, 400, `{"message":"invalid data"}`, false},
			"invalid data"},
		{"http 400, responseBody={message: 'invalid data'}, success=true --> invalid data",
			args{configloader.ResponseMsg{Path: "message"}, 400, `{"message":"invalid data"}`, true},
			"invalid data"},
		{"http 500, responseBody={message: 'system error'}, success=false --> system error",
			args{configloader.ResponseMsg{Path: "message"}, 500, `{"message":"system error"}`, false},
			"system error"},
		{"http 502, responseBody=empty, success=false --> default error",
			args{configloader.ResponseMsg{Path: "message"}, 502, "", false},
			errMsgByHttp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getResponseMessage(tt.args.responseMsg, tt.args.httpStatus, tt.args.responseBody, tt.args.isSuccess)
			if (len(got) == 0 && got != tt.wantMsgContains) || // no message
				(len(got) != 0 && !strings.Contains(got, tt.wantMsgContains)) { // have message contains
				t.Errorf("getResponseMessage() = %v, want contains %v", got, tt.wantMsgContains)
			}
		})
	}
}

func Test_checkMustHaveValueInPath(t *testing.T) {
	taskName := "this is task name"
	defaultMessage := "this is default message"
	messageContainTaskName := fmt.Sprintf("Không có dữ liệu khi %s", taskName)

	taskNoRequireField := configloader.ConfigTaskMD{
		TaskName: taskName,
		Response: configloader.ResponseMD{Code: configloader.ResponseCode{MustHaveValueInPath: ""}},
	}

	taskWithDataIsObject := configloader.ConfigTaskMD{
		TaskName: taskName,
		Response: configloader.ResponseMD{Code: configloader.ResponseCode{MustHaveValueInPath: "data.name"}},
	}

	taskWithDataIsArray := configloader.ConfigTaskMD{
		TaskName: taskName,
		Response: configloader.ResponseMD{Code: configloader.ResponseCode{MustHaveValueInPath: "data.txns.0.id"}},
	}

	type args struct {
		responseBody   string
		task           configloader.ConfigTaskMD
		isSuccess      bool
		defaultMessage string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// task is failed
		{"responseBody=empty, taskNoRequireField, failed --> default message",
			args{"", taskNoRequireField, false, defaultMessage},
			defaultMessage},
		{"responseBody contains required fields, taskWithDataIsObject, failed --> default message",
			args{`{"data":{"name":"quy"}`, taskWithDataIsObject, false, defaultMessage},
			defaultMessage},
		{"responseBody not contain required fields, taskWithDataIsArray, failed --> default message",
			args{`{"data":{"txns":[{"iddd":"123"}]}`, taskWithDataIsArray, false, defaultMessage},
			defaultMessage},
		// no required field
		{"responseBody=empty, taskNoRequireField, success --> default message",
			args{"", taskNoRequireField, true, defaultMessage},
			defaultMessage},
		{"responseBody contains required fields, taskNoRequireField, success --> default message",
			args{`{"data":{"name":"quy"}`, taskNoRequireField, true, defaultMessage},
			defaultMessage},
		{"responseBody not contain required fields, taskNoRequireField, success --> default message",
			args{`{"data":{"txns":[{"iddd":"123"}]}`, taskNoRequireField, true, defaultMessage},
			defaultMessage},

		// success, response body empty
		{"responseBody=empty, taskNoRequireField, success --> default message",
			args{"", taskNoRequireField, true, defaultMessage},
			defaultMessage},
		{"responseBody=empty, taskWithDataIsObject, success --> message contains taskName",
			args{"", taskWithDataIsObject, true, defaultMessage},
			messageContainTaskName},
		{"responseBody=empty, taskWithDataIsArray, success --> message contains taskName",
			args{"", taskWithDataIsArray, true, defaultMessage},
			messageContainTaskName},
		// success, response body contains required field
		{"responseBody not empty, taskNoRequireField, success --> default message",
			args{`{"data": "abc"}`, taskNoRequireField, true, defaultMessage},
			defaultMessage},
		{"responseBody contains required field, taskWithDataIsObject, success --> message contains taskName",
			args{`{"data":{"name":"quy"}`, taskWithDataIsObject, true, defaultMessage},
			defaultMessage},
		{"responseBody contains required field, taskWithDataIsArray, success --> message contains taskName",
			args{`{"data":{"txns":[{"id":"123"}]}`, taskWithDataIsArray, true, defaultMessage},
			defaultMessage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkRequiredFieldWhenTaskSuccess(tt.args.responseBody, tt.args.task, tt.args.isSuccess, tt.args.defaultMessage); got != tt.want {
				t.Errorf("checkMustHaveValueInPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
