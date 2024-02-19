package taskprovider

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	ct "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customtype"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
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
	}
	tests := []struct {
		name            string
		args            args
		wantMsgContains string
	}{
		{"http 200, responseBody=empty --> default error",
			args{configloader.ResponseMsg{Path: "message"}, 200, ""},
			errMsgByHttp},
		{"http 400, responseBody={message: 'invalid data'} --> invalid data",
			args{configloader.ResponseMsg{Path: "message"}, 400, `{"message":"invalid data"}`},
			"invalid data"},
		{"http 500, responseBody={message: 'system error'} --> system error",
			args{configloader.ResponseMsg{Path: "message"}, 500, `{"message":"system error"}`},
			"system error"},
		{"http 502, responseBody=empty --> default error",
			args{configloader.ResponseMsg{Path: "message"}, 502, ""},
			errMsgByHttp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getResponseMessage(tt.args.responseMsg, tt.args.httpStatus, tt.args.responseBody)
			if (len(got) == 0 && got != tt.wantMsgContains) || // no message
				(len(got) != 0 && !strings.Contains(got, tt.wantMsgContains)) { // have message contains
				t.Errorf("getResponseMessage() = %v, want contains %v", got, tt.wantMsgContains)
			}
		})
	}
}

func Test_checkMustHaveValueInPath(t *testing.T) {
	taskName := "this is task name"
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
		responseBody string
		task         configloader.ConfigTaskMD
	}
	tests := []struct {
		name                 string
		args                 args
		wantHasRequiredField bool
		wantMessage          string
	}{
		// task no require field
		{"responseBody=empty, taskNoRequireField --> no missing",
			args{"", taskNoRequireField},
			false, ""},

		// require field but it is not existed
		{"responseBody not contain required fields, taskWithDataIsArray --> missing",
			args{`{"data":{"txns":[{"iddd":"123"}]}`, taskWithDataIsArray},
			true, messageContainTaskName},
		{"responseBody not contain required fields, taskWithDataIsArrayAndPathHasFilter --> missing",
			args{`{"data":{"txns":[{"nameeee":"quy"}]}`, taskWithDataIsArrayAndPathHasFilter},
			true, messageContainTaskName},

		// require field and it is existed
		{"responseBody not empty, taskNoRequireField --> no missing",
			args{`{"data": "abc"}`, taskNoRequireField},
			false, ""},
		{"responseBody contains required field, taskWithDataIsObject --> no missing",
			args{`{"data":{"name":"quy"}`, taskWithDataIsObject},
			false, ""},
		{"responseBody contains required field, taskWithDataIsArray --> no missing",
			args{`{"data":{"txns":[{"id":"123"}]}`, taskWithDataIsArray},
			false, ""},
		{"responseBody contains required field, taskWithDataIsArrayAndPathHasFilter --> no missing",
			args{`{"data":{"txns":[{"name":"123"},{"name":"quy"},{"name":"abc"}]}`, taskWithDataIsArrayAndPathHasFilter},
			false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHasRequiredField, gotMessage := checkMissingRequiredField(tt.args.responseBody, tt.args.task)
			if gotMessage != tt.wantMessage || gotHasRequiredField != tt.wantHasRequiredField {
				t.Errorf("checkMustHaveValueInPath() = %v,%v while want %v, %v",
					gotMessage, gotHasRequiredField, tt.wantMessage, tt.wantHasRequiredField)
			}
		})
	}
}

func Test_transformMessage(t *testing.T) {
	msgRes := "this is message response"
	fileHeader := []string{"Header A", "Header B", "Header C"}
	rowData := []string{"value A", "value B", "value C"}

	emptyMsgTransforms := map[int]configloader.MessageTransformation{}

	msgTransforms := map[int]configloader.MessageTransformation{
		0:   {Message: "default message {{ $A }} abc"},
		400: {HttpStatus: 400, Message: "message for http 400 {{ $header.B }} abc"},
		500: {HttpStatus: 500, Message: "message for http 500 {{ $header.B }} {{ $A }} {{$response.message}} xyz"},
	}

	type args struct {
		httpStatus      int
		messageResponse string
		msgTransforms   map[int]configloader.MessageTransformation
		fileHeader      []string
		rowData         []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test transformMessage: no config Message Transformations -> empty message", args{400, msgRes, emptyMsgTransforms, fileHeader, rowData},
			msgRes},
		{"test transformMessage: httpStatus=400 -> correct message", args{400, msgRes, msgTransforms, fileHeader, rowData},
			"message for http 400 Header B abc"},
		{"test transformMessage: httpStatus=500 -> correct message", args{500, msgRes, msgTransforms, fileHeader, rowData},
			fmt.Sprintf("message for http 500 %s %s %s xyz", "Header B", "value A", msgRes)},
		{"test transformMessage: httpStatus = 401 (not in config) -> default message", args{401, msgRes, msgTransforms, fileHeader, rowData},
			"default message value A abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformMessage(tt.args.httpStatus, tt.args.messageResponse, tt.args.msgTransforms, tt.args.fileHeader, tt.args.rowData); got != tt.want {
				t.Errorf("transformMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchMessagePattern(t *testing.T) {
	fileHeader := []string{"Header A", "Header B", "Header C"}
	rowData := []string{"value A", "value B", "value C"}
	messageResponse := "this is message response"

	type args struct {
		messagePattern  string
		fileHeader      []string
		rowData         []string
		messageResponse string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// messagePattern empty or no replacement
		{"test matchMessagePattern: messagePattern empty -> empty",
			args{"", fileHeader, rowData, messageResponse},
			""},
		{"test matchMessagePattern: messagePattern no replacement -> empty",
			args{"this is message", fileHeader, rowData, messageResponse},
			"this is message"},

		// messagePattern has 1 replacement, and it existed
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is column -> colum value",
			args{"abc {{$A}} def", fileHeader, rowData, messageResponse},
			"abc value A def"},
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is header -> header value",
			args{"abc {{ $header.A }} def", fileHeader, rowData, messageResponse},
			"abc Header A def"},
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is current response msg -> response message value",
			args{"abc {{$response.message}} def", fileHeader, rowData, messageResponse},
			fmt.Sprintf("abc %s def", messageResponse)},

		// messagePattern has 1 replacement, and it is not existed
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is not existed column -> no replace",
			args{"abc {{$Y}} def", fileHeader, rowData, messageResponse},
			"abc {{$Y}} def"},
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is header -> no replace",
			args{"abc {{ $header.M }} def", fileHeader, rowData, messageResponse},
			"abc {{ $header.M }} def"},

		// messagePattern has 1 replacement, and it wrong
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is not existed column -> no replace",
			args{"abc {{$AA}} def", fileHeader, rowData, messageResponse},
			"abc {{$AA}} def"},
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is header -> no replace",
			args{"abc {{ $header.BB }} def", fileHeader, rowData, messageResponse},
			"abc {{ $header.BB }} def"},
		{"test matchMessagePattern: messagePattern has 1 replacement, and it is current response msg -> no replace",
			args{"abc {{$response1.message}} def", fileHeader, rowData, messageResponse},
			"abc {{$response1.message}} def"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchMessagePattern(tt.args.messagePattern, tt.args.fileHeader, tt.args.rowData, tt.args.messageResponse); got != tt.want {
				t.Errorf("matchMessagePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_abc(t *testing.T) {
	logger.Infof("\nadjustmentTypeId = %s\n",
		gjson.Get(responseTest, `data.adjustmentTypes.#(reasons.#(reasonDescription=="quytm_3")).adjustmentTypeId`).String())

	logger.Infof("\nreasonId = %s\n",
		gjson.Get(responseTest, `data.adjustmentTypes.#(reasons.#(reasonDescription=="quytm_3")).reasons.#(reasonDescription=="quytm_3").reasonId`).String())

	//logger.Infof("\nreasonId = %s\n",
	//	gjson.Get(responseTest, `data.adjustmentTypes.#.reasons.#(reasonDescription=="quytm_3").reasonId`).String())
}

var responseTest = `
{
  "code": 0,
  "message": "string",
  "data": {
    "adjustmentTypes": [
      {
        "adjustmentTypeId": 1111,
        "adjustmentTypeName": "string",
        "reasons": [
          {
            "reasonId": 12,
            "reasonDescription": "quytm_1"
          }
        ]
      },
      {
        "adjustmentTypeId": 22222,
        "adjustmentTypeName": "string",
        "reasons": [
          {
            "reasonId": 34,
            "reasonDescription": "quytm_2"
          },
          {
            "reasonId": 56,
            "reasonDescription": "quytm_3"
          }
        ]
      }
    ]
  }
}
`

func Test_convertRequestParams(t *testing.T) {
	tests := []struct {
		name      string
		reqParams map[string]interface{}
		want      []ct.Pair[string, string]
	}{
		{"test case normal with all primitive type and case empty value",
			map[string]interface{}{
				"param_int":    123,
				"param_float":  123.567,
				"param_string": "abc",
				"param_bool":   true,
				"param_empty":  "", // expect ignore
			},
			[]ct.Pair[string, string]{
				{Key: "param_bool", Value: "true"},
				{Key: "param_float", Value: "123.567"},
				{Key: "param_int", Value: "123"},
				{Key: "param_string", Value: "abc"},
			}},
		{"test case has param that is array",
			map[string]interface{}{
				"param_array":   []interface{}{1, 2, 3},
				"param_array_2": []interface{}{"1", "2", "3"},
			},
			[]ct.Pair[string, string]{
				{Key: "param_array", Value: "1"},
				{Key: "param_array", Value: "2"},
				{Key: "param_array", Value: "3"},
				{Key: "param_array_2", Value: "1"},
				{Key: "param_array_2", Value: "2"},
				{Key: "param_array_2", Value: "3"},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertRequestParams(tt.reqParams)
			sort.Slice(got, func(i, j int) bool { return got[i].Key < got[j].Key })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertRequestParams() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func Test_replacePathParams(t *testing.T) {
	type args struct {
		endpoint     string
		reqFieldName string
		realValue    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test case normal",
			args{"http://localhost:8080/api/v1/abc/{$param1}", "param1", "123"},
			"http://localhost:8080/api/v1/abc/123"},
		{"test case normal not full param",
			args{"http://localhost:8080/api/v1/abc/{$param1}/def/{$param2}", "param2", "456"},
			"http://localhost:8080/api/v1/abc/{$param1}/def/456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplacePathParams(tt.args.endpoint, tt.args.reqFieldName, tt.args.realValue); got != tt.want {
				t.Errorf("ReplacePathParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
