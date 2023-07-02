package taskprovider

import (
	"encoding/json"
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
	isSuccess := isTaskSuccess(task.Response, httpStatus, responseBody)
	messageRes := getResponseMessage(task.Response.Message, httpStatus, responseBody, isSuccess)
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
	for reqFieldName, reqField := range task.RequestParamsMap {
		realValue, err := getValueStringFromConfig(reqField, previousResponses)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		} else {
			task.RequestParams[reqFieldName] = realValue
		}
	}

	// 4. Convert request body
	for reqFieldName, reqField := range task.RequestBodyMap {
		// 4.1. Convert ArrayItem
		if len(reqField.ArrayItemMap) > 0 {
			// 4.1.1. For each child fields
			for reqFieldChildName, reqFieldChild := range reqField.ArrayItemMap {
				// get value
				realChildValue, err := getValueStringFromConfig(reqFieldChild, previousResponses)
				if err != nil {
					return configloader.ConfigTaskMD{}, err
				}

				// set value to RequestBody, that will be used for requesting to api
				// task.RequestBody[reqFieldName] is following format `array[map[string]interface{}]`
				items, err := setValueForChild(realChildValue, task.RequestBody[reqFieldName], reqFieldName, reqFieldChildName)
				if err != nil {
					return configloader.ConfigTaskMD{}, err
				} else {
					task.RequestBody[reqFieldName] = items
				}
			}

			// 4.1.2. Continue -> no convert realValue for Array Item
			continue
		}

		// 4.2. Get value
		realValue, err := getValueStringFromConfig(reqField, previousResponses)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		} else {
			task.RequestBody[reqFieldName] = realValue
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

func getValueStringFromConfig(reqField *configloader.RequestFieldMD, previousResponses map[int32]string) (interface{}, error) {
	// 1. Get value in String type
	var valueStr string
	switch reqField.ValueDependsOn {
	case configloader.ValueDependsOnTask:
		valueInTask, err := getValueByPreviousTaskResponse(reqField, previousResponses)
		if err != nil {
			return "", err
		}
		valueStr = valueInTask
	default:
		return "", fmt.Errorf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
	}

	// 2. Get real value then return
	return convertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
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
	if !codeRes.Exists() || // case not existed
		(len(codeRes.String()) == 0 && reqField.Required) { // case no value
		logger.Errorf("---- get data by path %v, but not found in previous response %v", reqField.ValueDependsOnKey, previousResponses)
		return "", fmt.Errorf("path `%v` not existed in response of task %v", reqField.ValueDependsOnKey, dependOn)
	}

	return codeRes.String(), nil
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

func getResponseMessage(responseMsg configloader.ResponseMsg, httpStatus int, responseBody string, isSuccess bool) string {
	messageRes := gjson.Get(responseBody, responseMsg.Path).String()

	if !isSuccess && len(messageRes) == 0 { // if failed and no message -> return error with http status
		messageRes = fmt.Sprintf("%v %v", errMsgByHttp, httpStatus)
	}

	return messageRes
}

func convertToRealValue(fieldType string, valueStr string, dependsOnKey string) (interface{}, error) {
	var realValue interface{}
	switch strings.ToLower(fieldType) {
	case configloader.TypeString:
		realValue = valueStr
	case configloader.TypeInteger:
		if valueInt64, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			realValue = valueInt64
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeNumber:
		if valueFloat64, err := strconv.ParseFloat(valueStr, 64); err == nil {
			realValue = valueFloat64
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeBoolean:
		if valueBool, err := strconv.ParseBool(valueStr); err == nil {
			realValue = valueBool
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeJson:
		result := gjson.Parse(valueStr)
		realValue = result.Value()
	default:
		return nil, fmt.Errorf("%s %s", errTypeNotSupport, fieldType)
	}
	return realValue, nil
}

func setValueForChild(realChildValue interface{}, items interface{}, reqFieldName string, reqFieldChildName string) (interface{}, error) {
	// 1. Build first item value
	item := map[string]interface{}{reqFieldChildName: realChildValue}

	// 2. In case haven't valued -> return
	if items == nil { //
		return []interface{}{item}, nil
	}

	// 3. In case already have value
	switch itemsType := items.(type) {
	case []interface{}: // only check case items is Array
		// 3.1. in case items empty -> return
		if len(itemsType) <= 0 { //
			return []interface{}{item}, nil
		}

		// 3.2. array already has value, get first item -> convert to map -> set child value
		firstItem, err := json.Marshal(itemsType[0])
		if err != nil {
			logger.Errorf("failed convert first item in `%v` to string, raw=%+v", reqFieldName, itemsType[0])
			return nil, errors.New("system error")
		} else {
			firstItemMap := make(map[string]interface{})
			if err := json.Unmarshal(firstItem, &firstItemMap); err != nil {
				logger.Errorf("failed convert first item in `%v` to string, raw=%+v", reqFieldName, itemsType[0])
				return nil, errors.New("system error")
			}
			firstItemMap[reqFieldChildName] = realChildValue
			return []interface{}{firstItemMap}, nil
		}
	default:
		logger.Errorf("`%v` is not array, raw=%+v", reqFieldName, itemsType)
		return nil, errors.New("system error")
	}
}
