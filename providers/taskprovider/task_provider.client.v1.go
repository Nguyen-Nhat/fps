package taskprovider

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

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
	// 1. Request
	reqHeader := task.RequestHeader
	reqHeader["Content-Type"] = "application/json" // default header
	logger.Infof("Prepare call %v with header=%v, requestParams=%+v, requestBody=%+v", task.Endpoint, reqHeader, task.RequestParams, task.RequestBody)
	httpStatus, responseBody, curl, err := utils.SendHTTPRequestRaw(c.client, task.Method, task.Endpoint, reqHeader, task.RequestParams, task.RequestBody)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v", task.Endpoint, err)
		logger.Errorf("-----> curl:\n%s", curl)
		responseBody = fmt.Sprintf("httpStatus=%v, responseBody=%v, error=%v", httpStatus, responseBody, err)
		return curl, responseBody, false, err.Error()
	}

	// 2. Handle response
	logger.Infof("response is status=%v, body=%v", httpStatus, responseBody)
	isSuccess := isTaskSuccess(task.Response, httpStatus, responseBody)
	messageRes := getResponseMessage(task.Response.Message, httpStatus, responseBody, isSuccess)
	messageRes, isSuccess = checkRequiredFieldWhenTaskSuccess(responseBody, task, isSuccess, messageRes)
	return curl, responseBody, isSuccess, messageRes
}

// checkRequiredFieldWhenTaskSuccess ... using this function in case task is success, but response body not contains Required Field
func checkRequiredFieldWhenTaskSuccess(responseBody string, task configloader.ConfigTaskMD, isSuccess bool, defaultMessage string) (string, bool) {
	mustHaveValueInPath := task.Response.Code.MustHaveValueInPath
	if !isSuccess || len(mustHaveValueInPath) <= 0 {
		return defaultMessage, isSuccess
	}

	if value := gjson.Get(responseBody, mustHaveValueInPath).String(); len(value) > 0 {
		// case have value in path -> return default message
		return defaultMessage, isSuccess
	} else {
		// case NOT have value by path -> return task failed
		return fmt.Sprintf("Không có dữ liệu khi %s", task.TaskName), false
	}
}

// private method ------------------------------------------------------------------------------------------------------

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

func getResponseMessage(responseMsg configloader.ResponseMsg, httpStatus int, responseBody string, isSuccess bool) string {
	messageRes := gjson.Get(responseBody, responseMsg.Path).String()

	if !isSuccess && len(messageRes) == 0 { // if failed and no message -> return error with http status
		messageRes = fmt.Sprintf("%v %v", errMsgByHttp, httpStatus)
	}

	return messageRes
}
