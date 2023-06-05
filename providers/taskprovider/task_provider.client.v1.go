package taskprovider

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

type (
	// IClientV1 ...
	IClientV1 interface {
		Execute(taskIndex int, taskMappingStr string, previousResponse map[int32]string) (map[string]interface{}, string, string, bool, string)
	}

	// ClientV1 ....
	clientV1 struct {
		client *http.Client
	}
)

var _ IClientV1 = &clientV1{}

// NewClientV1 ...
func NewClientV1() IClientV1 {
	return &clientV1{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *clientV1) Execute(taskIndex int, taskMappingStr string, previousResponse map[int32]string) (map[string]interface{}, string, string, bool, string) {
	// 1. Load Data and Mapping
	configMapping, err := converter.StringJsonToStruct("data_raw", taskMappingStr, configloader.ConfigMappingMD{})
	if err != nil {
		return nil, "", "", false, "failed to load config map"
	}

	// 2. Build request
	task, err := mapDataByPreviousResponse(taskIndex, *configMapping, previousResponse)
	if err != nil {
		return nil, "", "", false, err.Error()
	}

	// 3. Request
	reqHeader := task.Header
	reqHeader["Content-Type"] = "application/json"
	logger.Infof("Prepare call %v with header=%v, requestParams=%+v, requestBody=%+v", task.Endpoint, reqHeader, task.RequestParams, task.RequestBody)
	httpStatus, responseBody, curl, err := utils.SendHTTPRequestRaw(c.client, task.Method, task.Endpoint, reqHeader, task.RequestParams, task.RequestBody)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v", task.Endpoint, err)
		logger.Errorf("-----> curl:\n%s", curl)
		responseBody = fmt.Sprintf("httpStatus=%v, responseBody=%v, error=%v", httpStatus, responseBody, err)
		return task.RequestBody, curl, responseBody, false, err.Error()
	}

	// 4. Handle response
	logger.Infof("response is status=%v, body=%v", httpStatus, responseBody)
	isSuccess := isTaskSuccess(responseBody, task, httpStatus)
	messageRes := gjson.Get(responseBody, task.Response.Message.Path).String()
	return task.RequestBody, curl, responseBody, isSuccess, messageRes
}

// private method ------------------------------------------------------------------------------------------------------

// mapDataByPreviousResponse ...
func mapDataByPreviousResponse(taskIndex int, configMapping configloader.ConfigMappingMD, previousResponses map[int32]string) (
	configloader.ConfigTaskMD, error) {
	// 1. Get Task config
	task, isTaskExisted := getTaskConfig(taskIndex, configMapping)
	if !isTaskExisted {
		return configloader.ConfigTaskMD{}, fmt.Errorf("wrong config")
	}

	// 2. If no remaining request field that need to convert data -> do nothing
	if len(task.RequestParamsMap) == 0 && len(task.RequestBodyMap) == 0 {
		return task, nil
	}

	// 3. Convert request params
	for _, reqField := range task.RequestParamsMap {
		// 1.2.1. Get value in String type
		var valueStr string
		switch reqField.ValueDependsOn {
		case configloader.ValueDependsOnTask:
			valueInTask, err := getValueByPreviousTaskResponse(reqField, previousResponses)
			if err != nil {
				return configloader.ConfigTaskMD{}, err
			}
			valueStr = valueInTask
		default:
			return configloader.ConfigTaskMD{}, fmt.Errorf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
		}

		// 1.2.2. Get real value
		realValue, errMsg := convertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
		if len(errMsg) > 0 {
			return configloader.ConfigTaskMD{}, errors.New(errMsg)
		} else {
			task.RequestParams[reqField.Field] = realValue
		}
	}

	// 4. Convert request body
	// todo: haven't supported ValueDependsOnTask for nested object (defined in RequestFieldMD.ArrayItem) -> will update later
	for _, reqField := range task.RequestBodyMap {
		if reqField.Type == configloader.TypeArray {
			continue // ignore this type
		}
		// 1.2.1. Get value in String type
		var valueStr string
		switch reqField.ValueDependsOn {
		case configloader.ValueDependsOnTask:
			valueInTask, err := getValueByPreviousTaskResponse(reqField, previousResponses)
			if err != nil {
				return configloader.ConfigTaskMD{}, err
			}
			valueStr = valueInTask
		default:
			return configloader.ConfigTaskMD{}, fmt.Errorf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
		}

		// 1.2.2. Get real value
		realValue, errMsg := convertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
		if len(errMsg) > 0 {
			logger.Errorf("convertToRealValue ... error = %s", errMsg)
			return configloader.ConfigTaskMD{}, errors.New(errMsg)
		} else {
			task.RequestBody[reqField.Field] = realValue
		}
	}

	// 5. return
	return task, nil
}

func getTaskConfig(taskIndex int, configMapping configloader.ConfigMappingMD) (configloader.ConfigTaskMD, bool) {
	if len(configMapping.Tasks) == 0 {
		return configloader.ConfigTaskMD{}, false
	}

	for _, t := range configMapping.Tasks {
		if t.TaskIndex == taskIndex {
			return t, true
		}
	}
	return configloader.ConfigTaskMD{}, false
}

func getValueByPreviousTaskResponse(reqField *configloader.RequestFieldMD, previousResponses map[int32]string) (string, error) {
	dependOn := int32(reqField.ValueDependsOnTaskID)

	// get from previous task
	previousResponse, existed := previousResponses[dependOn]
	if !existed {
		logger.Errorf("task %v not existed", dependOn)
		return "", fmt.Errorf("no task contain %v in response", reqField.ValueDependsOnKey)
	}

	// get data by path
	codeRes := gjson.Get(previousResponse, reqField.ValueDependsOnKey)
	if !codeRes.Exists() {
		logger.Errorf("---- get data by path %v, but not found in previous response %v", reqField.ValueDependsOnKey, previousResponses)
		return "", fmt.Errorf("path %v not existed in response of previous task", reqField.ValueDependsOnKey)
	}

	return codeRes.String(), nil
}

func isTaskSuccess(responseBody string, task configloader.ConfigTaskMD, httpStatus int) bool {
	// 1. case check by httpStatus
	responseMD := task.Response
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

func convertToRealValue(fieldType string, valueStr string, dependsOnKey string) (interface{}, string) {
	var realValue interface{}
	// todo support ARRAY
	switch strings.ToUpper(fieldType) {
	case configloader.TypeString:
		realValue = valueStr
	//case configloader.TypeInt: // todo re-check
	//	if valueInt, err := strconv.Atoi(valueStr); err == nil {
	//		realValue = valueInt
	//	} else {
	//		return nil, fmt.Sprintf("%s (%s)", errTypeWrong, dependsOnKey)
	//	}
	case configloader.TypeInt, configloader.TypeLong:
		if valueInt64, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			realValue = valueInt64
		} else {
			return nil, fmt.Sprintf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	default:
		return nil, fmt.Sprintf("%s %s", errTypeNotSupport, fieldType)
	}
	return realValue, ""
}
