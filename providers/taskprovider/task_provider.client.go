package taskprovider

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/tidwall/gjson"
)

type (
	// IClient ...
	IClient interface {
		Execute(dataRaw string, taskMappingStr string, previousResponse map[int32]string) (map[string]interface{}, string, bool, string)
	}

	// Client ....
	Client struct {
		client *http.Client
	}
)

var _ IClient = &Client{}

// NewClient ...
func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *Client) Execute(dataRaw string, taskMappingStr string, previousResponse map[int32]string) (map[string]interface{}, string, bool, string) {
	// 1. Load Data and Mapping
	taskMapping, dataMap, err := convertDataAndMappingFromRawData(dataRaw, taskMappingStr)
	if err != nil {
		return nil, "", false, err.Error()
	}

	// 2. Build request
	reqHeader, reqEndPoint, requestBody, err := buildRequestForTask(taskMapping, dataMap, previousResponse)
	if err != nil {
		return nil, "", false, err.Error()
	}

	// 3. Request
	reqHeader["Content-Type"] = "application/json"
	logger.Infof("Prepare call %v with header=%v, requestBody=%v", reqEndPoint, reqHeader, requestBody)
	httpStatus, responseBody, _, err := utils.SendHTTPRequestRaw(c.client, http.MethodPost, reqEndPoint, reqHeader, nil, requestBody)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v", reqEndPoint, err)
		responseBody = fmt.Sprintf("httpStatus=%v, responseBody=%v, error=%v", httpStatus, responseBody, err)
		return requestBody, responseBody, false, err.Error()
	}

	// 4. Handle response
	logger.Infof("response is status=%v, body=%v", httpStatus, responseBody)
	isSuccess := isSuccessBaseOnMapping(responseBody, taskMapping, httpStatus)
	messageRes := gjson.Get(responseBody, taskMapping.Response.Message.Path).String()
	return requestBody, responseBody, isSuccess, messageRes
}

// private method ------------------------------------------------------------------------------------------------------

func convertDataAndMappingFromRawData(dataRaw string, taskMappingStr string) (*dto.MappingRow, map[string]string, error) {
	// 1. Convert Data
	dataMap, _ := converter.StringToMap("data_raw", dataRaw, true)
	if len(dataMap) == 0 {
		logger.ErrorT("data map is empty")
		return nil, nil, fmt.Errorf("failed to load data map")
	}

	// 2. Convert Mapping
	taskMapping, err := converter.StringJsonToStruct("mapping", taskMappingStr, dto.MappingRow{})
	if err != nil {
		logger.ErrorT("cannot convert mapping")
		return nil, nil, fmt.Errorf("failed to load mapping")
	}

	return taskMapping, dataMap, nil
}

// buildRequestForTask ... return (reqHeader, reqEndPoint, requestBody, error)
func buildRequestForTask(taskMapping *dto.MappingRow, dataMap map[string]string, previousResponses map[int32]string) (
	map[string]interface{}, string, map[string]interface{}, error) {
	requestBody := make(map[string]interface{})
	for _, reqMapping := range taskMapping.Request {
		if reqMapping.IsMappingExcel {
			reqMapping.Value = dataMap[reqMapping.MappingKey]
		}
		if reqMapping.IsMappingResponse {
			value, err := getValueBaseOnPreviousTaskResponse(reqMapping, previousResponses)
			if err != nil {
				return nil, "", nil, err
			}
			reqMapping.Value = value
		}
		if len(reqMapping.Value) > 0 {
			requestBody[reqMapping.FieldName] = reqMapping.Value
		}
	}
	logger.Infof("----- Request object = %v", requestBody)
	reqEndPoint := taskMapping.Endpoint
	reqHeader := taskMapping.Header
	return reqHeader, reqEndPoint, requestBody, nil
}

func getValueBaseOnPreviousTaskResponse(reqMapping dto.MappingRequest, previousResponses map[int32]string) (string, error) {
	dependOn := int32(reqMapping.DependOnResponseOfTaskId)

	// get from previous task
	previousResponse, existed := previousResponses[dependOn]
	if !existed {
		logger.Errorf("task %v not existed", dependOn)
		return "", fmt.Errorf("no task contain %v in response", reqMapping.MappingKey)
	}

	// get data by path
	codeRes := gjson.Get(previousResponse, reqMapping.MappingKey)
	if !codeRes.Exists() {
		logger.Errorf("---- get data by path %v, but not found in previous response %v", reqMapping.MappingKey, previousResponses)
		return "", fmt.Errorf("path %v not existed in response of previous task", reqMapping.MappingKey)
	}

	return codeRes.String(), nil
}

func isSuccessBaseOnMapping(responseBody string, taskMapping *dto.MappingRow, httpStatus int) bool {
	// 1. case check by httpStatus
	if len(taskMapping.Response.HttpStatusSuccess) != 0 {
		return strconv.Itoa(httpStatus) == taskMapping.Response.HttpStatusSuccess
	}

	// 2. case check by Code in response body
	codePath := taskMapping.Response.Code.Path
	codeSuccessValues := taskMapping.Response.Code.SuccessValues
	codeRes := gjson.Get(responseBody, codePath)
	if !codeRes.Exists() {
		logger.Errorf("---- get code by path %v, but not found in response %v", codePath, responseBody)
		return false
	}
	// code in response belongs to mapping
	return strings.Contains(codeSuccessValues, codeRes.String())
}
