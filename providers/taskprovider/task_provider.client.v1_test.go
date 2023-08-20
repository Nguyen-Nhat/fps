package taskprovider

import (
	"fmt"
	"strings"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
)

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

	pathHasFilterByName := `data.txns.#(name=="quy").name`
	taskWithDataIsArrayAndPathHasFilter := configloader.ConfigTaskMD{
		TaskName: taskName,
		Response: configloader.ResponseMD{Code: configloader.ResponseCode{MustHaveValueInPath: pathHasFilterByName}},
	}

	type args struct {
		responseBody   string
		task           configloader.ConfigTaskMD
		isSuccess      bool
		defaultMessage string
	}
	tests := []struct {
		name          string
		args          args
		wantMessage   string
		wantIsSuccess bool
	}{
		// task is failed
		{"responseBody=empty, taskNoRequireField, failed --> default message",
			args{"", taskNoRequireField, false, defaultMessage},
			defaultMessage, false},
		{"responseBody contains required fields, taskWithDataIsObject, failed --> default message",
			args{`{"data":{"name":"quy"}`, taskWithDataIsObject, false, defaultMessage},
			defaultMessage, false},
		{"responseBody not contain required fields, taskWithDataIsArray, failed --> default message",
			args{`{"data":{"txns":[{"iddd":"123"}]}`, taskWithDataIsArray, false, defaultMessage},
			defaultMessage, false},
		{"responseBody not contain required fields, taskWithDataIsArrayAndPathHasFilter, failed --> default message",
			args{`{"data":{"txns":[{"nameeee":"quy"}]}`, taskWithDataIsArrayAndPathHasFilter, false, defaultMessage},
			defaultMessage, false},

		// no required field
		{"responseBody=empty, taskNoRequireField, success --> default message",
			args{"", taskNoRequireField, true, defaultMessage},
			defaultMessage, true},
		{"responseBody contains required fields, taskNoRequireField, success --> default message",
			args{`{"data":{"name":"quy"}`, taskNoRequireField, true, defaultMessage},
			defaultMessage, true},
		{"responseBody not contain required fields, taskNoRequireField, success --> default message",
			args{`{"data":{"txns":[{"iddd":"123"}]}`, taskNoRequireField, true, defaultMessage},
			defaultMessage, true},

		// success, response body empty
		{"responseBody=empty, taskNoRequireField, success --> default message",
			args{"", taskNoRequireField, true, defaultMessage},
			defaultMessage, true},
		{"responseBody=empty, taskWithDataIsObject, success --> message contains taskName",
			args{"", taskWithDataIsObject, true, defaultMessage},
			messageContainTaskName, false},
		{"responseBody=empty, taskWithDataIsArray, success --> message contains taskName",
			args{"", taskWithDataIsArray, true, defaultMessage},
			messageContainTaskName, false},
		{"responseBody=empty, taskWithDataIsArrayAndPathHasFilter, success --> message contains taskName",
			args{"", taskWithDataIsArrayAndPathHasFilter, true, defaultMessage},
			messageContainTaskName, false},

		// success, response body contains required field
		{"responseBody not empty, taskNoRequireField, success --> default message",
			args{`{"data": "abc"}`, taskNoRequireField, true, defaultMessage},
			defaultMessage, true},
		{"responseBody contains required field, taskWithDataIsObject, success --> message contains taskName",
			args{`{"data":{"name":"quy"}`, taskWithDataIsObject, true, defaultMessage},
			defaultMessage, true},
		{"responseBody contains required field, taskWithDataIsArray, success --> message contains taskName",
			args{`{"data":{"txns":[{"id":"123"}]}`, taskWithDataIsArray, true, defaultMessage},
			defaultMessage, true},
		{"responseBody contains required field, taskWithDataIsArrayAndPathHasFilter, success --> message contains taskName",
			args{`{"data":{"txns":[{"name":"123"},{"name":"quy"},{"name":"abc"}]}`, taskWithDataIsArrayAndPathHasFilter, true, defaultMessage},
			defaultMessage, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, gotIsSuccess := checkRequiredFieldWhenTaskSuccess(tt.args.responseBody, tt.args.task, tt.args.isSuccess, tt.args.defaultMessage)
			if gotMessage != tt.wantMessage || gotIsSuccess != tt.wantIsSuccess {
				t.Errorf("checkMustHaveValueInPath() = %v,%v while want %v, %v",
					gotMessage, gotIsSuccess, tt.wantMessage, tt.wantIsSuccess)
			}
		})
	}
}
