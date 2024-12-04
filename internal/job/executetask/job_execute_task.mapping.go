package executetask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

var regexDoubleBrace = regexp.MustCompile(`\{\{(.*?)}}`)

// mapDataByPreviousResponseAndCustomFunction ...
func mapDataByPreviousResponseAndCustomFunction(ctx context.Context, processingFileRow *fileprocessingrow.ProcessingFileRow, configMapping configloader.ConfigMappingMD, previousResponses map[int32]string) (
	configloader.ConfigTaskMD, error) {
	// 1. Get Task config
	task, isTaskExisted := getTaskConfig(int(processingFileRow.TaskIndex), configMapping)
	if !isTaskExisted {
		return configloader.ConfigTaskMD{}, fmt.Errorf("wrong config")
	}

	// 2. If no remaining request field that need to convert data -> do nothing
	if len(task.RequestParamsMap) == 0 && len(task.RequestBodyMap) == 0 && len(task.PathParamsMap) == 0 {
		return task, nil
	}

	// 3. Convert path params
	for reqFieldName, reqField := range task.PathParamsMap {
		realValue, err := getValueStringFromConfig(ctx, processingFileRow, reqField, previousResponses, task.ImportRowData)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		}
		task.PathParams[reqFieldName] = realValue
		delete(task.PathParamsMap, reqFieldName)
	}

	// 3. Convert request params
	for reqFieldName, reqField := range task.RequestParamsMap {
		realValue, err := getValueStringFromConfig(ctx, processingFileRow, reqField, previousResponses, task.ImportRowData)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		} else {
			task.RequestParams[reqFieldName] = realValue
			delete(task.RequestParamsMap, reqFieldName)
		}
	}

	// 4. Convert request body
	for reqFieldName, reqField := range task.RequestBodyMap {
		// 4.1. Convert ArrayItem
		if len(reqField.ArrayItemMap) > 0 {
			// 4.1.1. For each child fields
			childMap, err := getValueFromConfig(ctx, processingFileRow, reqFieldName, reqField.ArrayItemMap, previousResponses, task.ImportRowData)
			if err != nil {
				return configloader.ConfigTaskMD{}, err
			}

			if task.RequestBody[reqFieldName] != nil {
				converter.Override([]map[string]interface{}{childMap}, task.RequestBody[reqFieldName])
			} else {
				task.RequestBody[reqFieldName] = []map[string]interface{}{childMap}
			}
			if len(childMap) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
				if val, ok := childMap[""]; ok {
					task.RequestBody[reqFieldName] = []interface{}{val}
				}
			} else if len(childMap) == 0 { // childMap is empty -> remove field
				delete(task.RequestBody, reqFieldName)
			}
			task.RequestBody[reqFieldName] = []map[string]interface{}{childMap}

			// 4.1.2. Continue -> no convert realValue for Array Item
			continue
		}
		if len(reqField.ItemsMap) > 0 {
			childMap, err := getValueFromConfig(ctx, processingFileRow, reqFieldName, reqField.ItemsMap, previousResponses, task.ImportRowData)
			if err != nil {
				return configloader.ConfigTaskMD{}, err
			}

			task.RequestBody[reqFieldName] = childMap
			continue
		}

		// 4.2. Get value
		realValue, err := getValueStringFromConfig(ctx, processingFileRow, reqField, previousResponses, task.ImportRowData)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		} else {
			task.RequestBody[reqFieldName] = realValue
			delete(task.RequestBodyMap, reqFieldName)
		}
	}

	// 5. return
	return task, nil
}

func getValueFromConfig(ctx context.Context, processingFileRow *fileprocessingrow.ProcessingFileRow, parentFieldName string, requestFieldList map[string]*configloader.RequestFieldMD, previousResponses map[int32]string, rowData []string) (
	map[string]interface{}, error) {

	childMap := make(map[string]interface{})

	for fieldNameChild, reqFieldChild := range requestFieldList {
		if len(reqFieldChild.ArrayItemMap) > 0 {
			childMapInArr, err := getValueFromConfig(ctx, processingFileRow, fieldNameChild, reqFieldChild.ArrayItemMap, previousResponses, rowData)
			if err != nil {
				return nil, err
			}

			childMap[fieldNameChild] = []map[string]interface{}{childMapInArr}
			if len(childMapInArr) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
				if val, ok := childMapInArr[""]; ok {
					childMap[fieldNameChild] = []interface{}{val}
				}
			} else if len(childMapInArr) == 0 { // childMapInArr is empty -> remove field
				delete(childMap, fieldNameChild)
			}

			// Continue -> no convert realValue for Array Item
			continue
		}

		if len(reqFieldChild.ItemsMap) > 0 {
			childMapInObj, err := getValueFromConfig(ctx, processingFileRow, fieldNameChild, reqFieldChild.ItemsMap, previousResponses, rowData)
			if err != nil {
				return nil, err
			}

			childMap[fieldNameChild] = childMapInObj

			// Continue -> no convert realValue for Array Item
			continue
		}

		// get value
		realChildValue, err := getValueStringFromConfig(ctx, processingFileRow, reqFieldChild, previousResponses, rowData)
		if err != nil {
			return nil, err
		}

		// set value to RequestBody, that will be used for requesting to api
		// task.RequestBody[reqFieldName] is following format `array[map[string]interface{}]`
		if reqFieldChild.Type == configloader.TypeArray {
			if children, err := setValueForChild(realChildValue, childMap[fieldNameChild], parentFieldName, fieldNameChild); err != nil {
				return nil, err
			} else {
				childMap[fieldNameChild] = children
			}
		} else if realChildValue != nil { // ignore case value is nil
			childMap[fieldNameChild] = realChildValue
		}
	}

	return childMap, nil
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

func getValueStringFromConfig(ctx context.Context, processingFileRow *fileprocessingrow.ProcessingFileRow, reqField *configloader.RequestFieldMD, previousResponses map[int32]string, rowData []string) (interface{}, error) {
	// 1. Get value in String type
	var valueStr string
	switch reqField.ValueDependsOn {
	case configloader.ValueDependsOnTask:
		valueInTask, err := getValueByPreviousTaskResponse(reqField, previousResponses)
		if err != nil {
			return nil, err
		}
		valueStr = valueInTask
	case configloader.ValueDependsOnFunc:
		// case custom function has args depend on previous response
		for idx, paramRaw := range reqField.ValueDependsOnFunc.ParamsRaw {
			if strings.HasPrefix(paramRaw, configloader.PrefixMappingRequestResponse) {
				realValue, err := getValueByPreviousTaskResponseForCustomFunc(paramRaw, previousResponses, rowData)
				if err != nil {
					return nil, err
				}
				reqField.ValueDependsOnFunc.ParamsMapped[idx] = realValue
			}
		}
		result, err := customFunc.ExecuteFunction(reqField.ValueDependsOnFunc)
		if err != nil {
			return nil, err
		} else if len(result.ErrorMessage) > 0 {
			return nil, fmt.Errorf(result.ErrorMessage)
		} else {
			return result.Result, nil
		}
	case configloader.ValueDependsOnDb:
		valueInDb, err := getValueFromFieldInDb(processingFileRow, reqField)
		if err != nil {
			return nil, err
		}
		return valueInDb, nil
	default:
		return nil, fmt.Errorf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
	}

	// 2. Get real value then return
	return basejobmanager.ConvertToRealValue(ctx, reqField.Type, valueStr, reqField.ValueDependsOnKey)
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

// Ex: paramRaw: $response1.data.products.0.productLifeCycle
// Ex: previousResponses: 1 -> {"code":0,"message":"Thao tác thành công","data":{"products":[{"sellerId":40,"sku":"10125","productLifeCycle":null,"isMarketShortage":null,"expectedEndOfShortageDate":null,"buyable":null,"autoReplenishment":null,"taxId":1,"skus":[],"isCoreProductLine":false}]}}
func getValueByPreviousTaskResponseForCustomFunc(valuePattern string, previousResponses map[int32]string, rowData []string) (string, error) {
	template := strings.TrimPrefix(valuePattern, configloader.PrefixMappingRequestResponse) // $response1.field_abc -> template = 1.field_abc
	if len(template) <= 2 || template[1] != constant.SplitByDotChar {
		return constant.EmptyString, fmt.Errorf("mapping request is invalid: %v", valuePattern)
	}
	dependOnTaskId, err := strconv.Atoi(string(template[0])) // 1.field_abc -> 1
	if err != nil {
		logger.Errorf("mapping request is invalid: %v, err: %v", valuePattern, err)
		return constant.EmptyString, fmt.Errorf("mapping request is invalid: %v", valuePattern)
	}

	valueDependsOnKey := matchJsonPath(rowData, template[2:]) // 1.field_abc -> field_abc

	dependOn := int32(dependOnTaskId)

	// get from previous task
	previousResponse, existed := previousResponses[dependOn]
	if !existed {
		logger.Errorf("task %v not existed", dependOn)
		return "", fmt.Errorf("no task contain %v in response", valueDependsOnKey)
	}

	// get data by path
	codeRes := gjson.Get(previousResponse, valueDependsOnKey)
	if !codeRes.Exists() { // case not existed
		logger.Errorf("---- get data by path %v, but not found in previous response %v", valueDependsOnKey, previousResponses)
		return "", fmt.Errorf("path `%v` not existed in response of task %v", valueDependsOnKey, dependOn)
	}

	return codeRes.String(), nil
}

// matchJsonPath ... support validate json path and update json path if it contains variable
// for example:
//   - Json path: data.transactions.#(name=="{{ $A }}").id
//   - Excel data has $A = quy
//     -> output: data.transactions.#(name=="quy").id
func matchJsonPath(rowData []string, jsonPath string) string {
	// 1. Extract data with format like `{{ $A }}`
	matchers := regexDoubleBrace.FindStringSubmatch(jsonPath)

	// 2. Return if not match
	if len(matchers) != 2 {
		return jsonPath
	}

	// 3. Validate and Replace value
	valuePatternWithDoubleBrace := matchers[0]
	valuePattern := strings.TrimSpace(matchers[1])
	if len(valuePattern) == 0 {
		return jsonPath
	}

	// 3.1. Get column key: $A -> A
	if !excel.IsColumnIndex(valuePattern) {
		logger.Errorf("validateResponseCode ... error %+v", valuePattern)
		return jsonPath
	}
	columnKey := valuePattern[1:] // if `$A` -> columnIndex = `A`
	columnIndex, err := excelize.ColumnNameToNumber(columnKey)
	if err != nil {
		logger.Errorf("validateResponseCode ... error %+v", err)
		return jsonPath
	}

	// 3.2. Validate value
	if columnIndex > len(rowData) || // column request out of range
		len(strings.TrimSpace(rowData[columnIndex-1])) == 0 { // column is required by value is empty
		return jsonPath
	}

	// 3.3. Replace value
	cellValue := strings.TrimSpace(rowData[columnIndex-1])
	jsonPath = strings.ReplaceAll(jsonPath, valuePatternWithDoubleBrace, cellValue)

	// 4. return
	return jsonPath
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

func getValueFromFieldInDb(processingFileRow *fileprocessingrow.ProcessingFileRow, reqField *configloader.RequestFieldMD) (interface{}, error) {
	dependOnField := reqField.ValueDependsOnKey
	switch dependOnField {
	case configloader.ValueDependsOnDbFieldTaskId:
		return processingFileRow.ID, nil
	case configloader.ValueDependsOnDbFieldFileId:
		return processingFileRow.FileID, nil
	default:
		return nil, fmt.Errorf("cannot get value from field %v", dependOnField)
	}
}
