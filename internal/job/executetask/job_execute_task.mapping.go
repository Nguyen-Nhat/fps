package executetask

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

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
	return basejobmanager.ConvertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
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