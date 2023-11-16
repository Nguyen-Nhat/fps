package configloader

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_toConfigTaskMD(t *testing.T) {
	http200 := int32(200)
	defaultConfigTask := configtask.ConfigTask{ConfigTask: ent.ConfigTask{
		ID:              1,
		ConfigMappingID: 1,
		TaskIndex:       2,
		Name:            "abc",
		EndPoint:        "https://abc.com/call",
		Method:          "POST",
		Header:          `{"authorization":"abcd", "sellerID": 1, "isEnable": "true"}`,
		RequestParams: `[
							{"field":"siteIds","type":"integer","valuePattern":"$param.siteId","required":true},
							{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true}
						]`,
		RequestBody:               loadJson("test_config_json_all_type_level_1.json"),
		ResponseSuccessHTTPStatus: int32(200),
		ResponseSuccessCodeSchema: `{
										"path":"code",
										"successValues":"0",
										"mustHaveValueInPath":"data.#(sellerSiteCode==\"{{ $A }}\").id"
									}`,
		ResponseMessageSchema: `{"path": "message"}`,
		MessageTransformations: `[
									{
										"httpStatus": 400,
										"message": "Error {{$header.A}} {{$A}}, {{$header.B}} {{$B}}: {{$response.message}}"
									},
									{ "httpStatus": 0, "message": "{{$response.message}}" }
								]`,
		GroupByColumns:   "A,B",
		GroupBySizeLimit: 300,
		CreatedAt:        time.Now(),
		CreatedBy:        "quy.tm@teko.vn",
		UpdatedAt:        time.Now(),
	}}

	wantConfigTask := ConfigTaskMD{
		TaskIndex: 2,
		TaskName:  "abc",

		Endpoint: "https://abc.com/call",
		Method:   "POST",

		RequestHeaderMap: make(map[string]*RequestFieldMD),
		RequestParamsMap: map[string]*RequestFieldMD{
			"siteIds": {
				Field:             "siteIds",
				Type:              "integer",
				ValuePattern:      "$param.siteId",
				Required:          true,
				ValueDependsOn:    "PARAM",
				ValueDependsOnKey: "siteId"},
			"sellerId": {
				Field:             "sellerId",
				Type:              "integer",
				ValuePattern:      "$param.sellerId",
				Required:          true,
				ValueDependsOn:    "PARAM",
				ValueDependsOnKey: "sellerId"},
		},
		RequestBodyMap: wantRequestBody_1(),
		RequestHeader: map[string]interface{}{
			"authorization": []interface{}{"abcd"}[0],
			"sellerID":      []interface{}{float64(1)}[0],
			"isEnable":      []interface{}{"true"}[0],
		},
		RequestParams: nil,
		RequestBody:   nil,

		Response: ResponseMD{
			&http200,
			ResponseCode{"code", "0", `data.#(sellerSiteCode=="{{ $A }}").id`},
			ResponseMsg{"message"},
			map[int]MessageTransformation{
				0:   {HttpStatus: 0, Message: "{{$response.message}}"},
				400: {HttpStatus: 400, Message: "Error {{$header.A}} {{$A}}, {{$header.B}} {{$B}}: {{$response.message}}"},
			},
		},

		RowGroup: RowGroupMD{"A,B", []int{0, 1}, 300},

		ImportRowHeader: nil,
		ImportRowData:   nil,
		ImportRowIndex:  0,
	}

	type args struct {
		task configtask.ConfigTask
	}
	tests := []struct {
		name    string
		args    args
		want    ConfigTaskMD
		wantErr bool
	}{
		{"case success", args{defaultConfigTask}, wantConfigTask, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toConfigTaskMD(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("toConfigTaskMD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, wantConfigTask.TaskIndex, got.TaskIndex)
			assert.Equal(t, wantConfigTask.TaskName, got.TaskName)
			assert.Equal(t, wantConfigTask.Endpoint, got.Endpoint)
			assert.Equal(t, wantConfigTask.Method, got.Method)
			assert.Equal(t, wantConfigTask.RequestHeaderMap, got.RequestHeaderMap)
			// Request Param
			assert.Equal(t, len(wantConfigTask.RequestParamsMap), len(got.RequestParamsMap))
			for key, wantValue := range wantConfigTask.RequestParamsMap {
				assert.NotNil(t, got.RequestParamsMap[key], fmt.Sprintf("RequestParamsMap.%s: \nwant: %+v \ngot : nil", key, *wantValue))
				assert.True(t, reflect.DeepEqual(*wantValue, *got.RequestParamsMap[key]), fmt.Sprintf("RequestParamsMap.%s: \nwant: %+v \ngot : %+v", key, *wantValue, *got.RequestParamsMap[key]))
			}
			// Request Body
			assert.Equal(t, len(wantConfigTask.RequestBodyMap), len(got.RequestBodyMap))
			for key, wantValue := range wantConfigTask.RequestBodyMap {
				assert.NotNil(t, got.RequestBodyMap[key], fmt.Sprintf("RequestBodyMap.%s: \nwant: %+v \ngot : nil", key, *wantValue))
				assert.True(t, reflect.DeepEqual(*wantValue, *got.RequestBodyMap[key]), fmt.Sprintf("RequestBodyMap.%s: \nwant: %+v \ngot : %+v", key, *wantValue, *got.RequestBodyMap[key]))
			}
			//assert.Equal(t, nil, got.RequestBodyMap)
			assert.True(t, reflect.DeepEqual(wantConfigTask.RequestHeader, got.RequestHeader), fmt.Sprintf("RequestHeader: \nwant: %+v \ngot : %+v", wantConfigTask.RequestHeader, got.RequestHeader))
			assert.Equal(t, map[string]interface{}{}, got.RequestParams)
			assert.Equal(t, map[string]interface{}{}, got.RequestBody)
			assert.True(t, reflect.DeepEqual(wantConfigTask.Response, got.Response), fmt.Sprintf("Response: \nwant: %+v \ngot : %+v", wantConfigTask.Response, got.Response))
			assert.True(t, reflect.DeepEqual(wantConfigTask.RowGroup, got.RowGroup), fmt.Sprintf("RowGroup: \nwant: %+v \ngot : %+v", wantConfigTask.RowGroup, got.RowGroup))
			assert.True(t, got.ImportRowHeader == nil)
			assert.True(t, got.ImportRowData == nil)
			assert.Equal(t, 0, got.ImportRowIndex)
		})
	}
}

func Test_response(t *testing.T) {
	var responseCode ResponseCode
	jsonStr := `{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\"{{ $A }}\").sellerSiteCode"}`
	if err := json.Unmarshal([]byte(jsonStr), &responseCode); err != nil {
		logger.Errorf("error when convert ResponseSuccessCodeSchema: value=%v, err=%v", jsonStr, err)
	} else {
		logger.Infof("Ok %+v", responseCode)
	}

}

func loadJson(path string) string {
	if data, err := os.ReadFile(path); err != nil {
		panic(err)
	} else {
		return string(data)
	}
}
