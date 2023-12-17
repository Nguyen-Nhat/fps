package taskprovider

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	t "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customtype"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

// var regexDoubleBrace = regexp.MustCompile(`{{([^{}]+)}}`)
var regexDoubleBrace = regexp.MustCompile(`\{\{(.*?)}}`)

type (
	// IClientV1 ...
	IClientV1 interface {
		Execute(configloader.ConfigTaskMD) (string, string, bool, string)
	}

	// ClientV1 ....
	clientV1 struct {
		client *http.Client
	}
)

const errMsgByHttp = "Xảy ra lỗi http"

var _ IClientV1 = &clientV1{}

// NewClientV1 ...
func NewClientV1() IClientV1 {
	return &clientV1{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *clientV1) Execute(task configloader.ConfigTaskMD) (string, string, bool, string) {
	// 0. Init and convert input
	reqHeader := task.RequestHeader
	reqHeader["Content-Type"] = "application/json" // default header
	reqParams := task.RequestParams
	listReqParams := convertRequestParams(reqParams)
	logger.Infof("Prepare call %v with header=%v, requestParams=%+v, requestBody=%+v", task.Endpoint, reqHeader, listReqParams, task.RequestBody)

	// 1. Request
	httpStatus, responseBody, curl, err := utils.SendHTTPRequestRaw(c.client, task.Method, task.Endpoint, reqHeader, listReqParams, task.RequestBody)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v", task.Endpoint, err)
		logger.Errorf("-----> curl:\n%s", curl)
		responseBody = fmt.Sprintf("httpStatus=%v, responseBody=%v, error=%v", httpStatus, responseBody, err)
		return curl, responseBody, false, err.Error()
	}

	// 2. Handle response
	logger.Infof("response is status=%v, body=%v", httpStatus, responseBody)
	isSuccess := isTaskSuccess(task.Response, httpStatus, responseBody)
	messageInResponse := getResponseMessage(task.Response.Message, httpStatus, responseBody)
	var messageDisplay string
	if isSuccess {
		if isMissing, defaultErrMsgWhenFieldMissing := checkMissingRequiredField(responseBody, task); isMissing {
			messageDisplay = transformMessage(httpStatusWhenMissingField, defaultErrMsgWhenFieldMissing, task.Response.MessageTransforms, task.ImportRowHeader, task.ImportRowData)
			isSuccess = false
		} else {
			messageDisplay = ""
		}
	} else {
		messageDisplay = transformMessage(httpStatus, messageInResponse, task.Response.MessageTransforms, task.ImportRowHeader, task.ImportRowData)
	}

	return curl, responseBody, isSuccess, messageDisplay
}

// private method ------------------------------------------------------------------------------------------------------

// convertRequestParams ...
func convertRequestParams(reqParams map[string]interface{}) []t.Pair[string, string] {
	// 1. Convert data to string
	var listParams []t.Pair[string, string]
	for k, v := range reqParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice, reflect.Array:
			if slice, ok := reflect.ValueOf(v).Interface().([]interface{}); ok {
				for _, param := range slice {
					listParams = append(listParams, t.Pair[string, string]{Key: k, Value: fmt.Sprintf("%v", param)})
				}
			}
		default:
			listParams = append(listParams, t.Pair[string, string]{Key: k, Value: fmt.Sprintf("%v", v)})
		}
	}

	// 2. Ignore empty param
	var listParamsIgnoreEmpty []t.Pair[string, string]
	for _, pair := range listParams {
		if len(pair.Key) > 0 && len(pair.Value) > 0 {
			listParamsIgnoreEmpty = append(listParamsIgnoreEmpty, pair)
		}
	}

	// 3. Return
	return listParamsIgnoreEmpty
}

func isTaskSuccess(responseMD configloader.ResponseMD, httpStatus int, responseBody string) bool {
	// 1. case check by httpStatus
	if responseMD.HttpStatusSuccess != nil {
		return *responseMD.HttpStatusSuccess == int32(httpStatus)
	}

	// 2. case check by Code in response body
	codePath := responseMD.Code.Path
	codeSuccessValues := responseMD.Code.SuccessValues
	codeRes := gjson.Get(responseBody, codePath)
	if !codeRes.Exists() {
		logger.Errorf("---- get code by path %v, but not found in response %v", codePath, responseBody)
		return false
	}
	// code in response belongs to mapping
	return strings.Contains(codeSuccessValues, codeRes.String())
}

func getResponseMessage(responseMsg configloader.ResponseMsg, httpStatus int, responseBody string) string {
	messageRes := gjson.Get(responseBody, responseMsg.Path).String()

	if len(messageRes) == 0 { // if failed and no message -> return error with http status
		messageRes = fmt.Sprintf("%v %v", errMsgByHttp, httpStatus)
	}

	return messageRes
}

// checkMissingRequiredField ... using this function in case task is success, but response body not contains Required Field
// if existed, return: true, ""
// else return false, <error message>
func checkMissingRequiredField(responseBody string, task configloader.ConfigTaskMD) (bool, string) {
	mustHaveValueInPath := task.Response.Code.MustHaveValueInPath
	if len(mustHaveValueInPath) <= 0 {
		return false, ""
	}

	if value := gjson.Get(responseBody, mustHaveValueInPath).String(); len(value) > 0 {
		// case have value in path -> return default message
		return false, ""
	} else {
		// case NOT have value by path -> return task failed
		return true, fmt.Sprintf("Không có dữ liệu khi %s", task.TaskName)
	}
}

// transformMessage ...
func transformMessage(httpStatus int, messageResponse string,
	msgTransforms map[int]configloader.MessageTransformation, fileHeader []string, rowData []string) string {
	// 1. If no config -> return default is messageResponse
	if len(msgTransforms) == 0 {
		return messageResponse
	}

	// 2. Find messageTransformation that match to httpStatus
	var matchedMsgTrans *configloader.MessageTransformation
	if msgTrans, ok := msgTransforms[httpStatus]; ok { // get by httpStatus
		matchedMsgTrans = &msgTrans
	} else if msgTransDefault, okDefault := msgTransforms[0]; okDefault { // get default when no config for httpStatus
		matchedMsgTrans = &msgTransDefault
	} else { // if not found default MessageTransformation -> return default is messageResponse
		return messageResponse
	}

	// 3. Get message pattern then map value
	messagePattern := matchedMsgTrans.Message
	return matchMessagePattern(messagePattern, fileHeader, rowData, messageResponse)
}

// matchMessagePattern ... replace value for messagePattern
// for example:
//   - MessagePattern: abc {{ $A }} def {{ $param.sellerId }}
//   - Excel data has $A = quy, parameters have sellerId=1
//     -> output: abc quy def 1
//
// Support pattern: $A, $header, $response.message
func matchMessagePattern(messagePattern string, fileHeader []string, rowData []string, messageResponse string) string {
	// 1. Extract data with format like `{{ ... }}`
	allMatchers := regexDoubleBrace.FindAllStringSubmatch(messagePattern, -1)

	// 2. Return if not match
	if len(allMatchers) == 0 {
		return messagePattern
	}

	// 3. Validate and Replace value
	for _, matchers := range allMatchers {
		// 3.0. If not have 2 matchers => ignore this case
		if len(matchers) != 2 {
			continue
		}

		// 3.1. Check value pattern
		valuePatternWithDoubleBrace := matchers[0]
		valuePattern := strings.TrimSpace(matchers[1])
		if !strings.HasPrefix(valuePattern, configloader.PrefixMappingRequest) {
			continue // ignore this case
		}

		// 3.2. Get replacement from value pattern
		replacement := ""
		if excel.IsColumnIndex(valuePattern) {
			columnKey := valuePattern[1:] // if `$A` -> columnIndex = `A`
			replacement = excel.GetValueFromColumnKey(columnKey, rowData)
		} else if strings.HasPrefix(valuePattern, configloader.PrefixMappingRequestHeader) {
			headerName := strings.TrimPrefix(valuePattern, configloader.PrefixMappingRequestHeader+".")
			replacement = excel.GetValueFromColumnKey(headerName, fileHeader)
		} else if configloader.PrefixMappingRequestCurrentResponseMessage == valuePattern {
			replacement = messageResponse
		}

		// 3.3. Replace value
		if len(replacement) > 0 {
			messagePattern = strings.ReplaceAll(messagePattern, valuePatternWithDoubleBrace, replacement)
		}
	}

	// 4. return
	return messagePattern
}
