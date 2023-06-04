package configloader

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

// databaseConfigLoaderV1 ...
type databaseConfigLoaderV1 struct {
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

func (cl *databaseConfigLoaderV1) Load(file fileprocessing.ProcessingFile) (ConfigMappingMD, error) {
	// 1. Get mapping in DB
	cfgMapping, cfgTasks, err := cl.loadConfigFromDB(file)
	if err != nil {
		logger.Errorf("find config mapping failed, error = %+v", err)
		return ConfigMappingMD{}, err
	}

	// 2. Convert mapping DB to ConfigTaskMD
	var tasks []ConfigTaskMD
	for _, task := range cfgTasks {
		cfg, err := toConfigTaskMD(*task)
		if err != nil {
			return ConfigMappingMD{}, err
		}
		tasks = append(tasks, cfg)
	}

	// 3. Return
	return ConfigMappingMD{
		// metadata get from file processing
		DataStartAtRow:     int(cfgMapping.DataStartAtRow),
		RequireColumnIndex: strings.Split(cfgMapping.RequireColumnIndex, ","),
		ErrorColumnIndex:   cfgMapping.ErrorColumnIndex,
		// parameter in file
		FileParameters: toFileParameters(file.FileParameters),
		// List ConfigTaskMD
		Tasks: tasks,
	}, nil
}

// Private methods -----------------------------------------------------------------------------------------------------

func (cl *databaseConfigLoaderV1) loadConfigFromDB(file fileprocessing.ProcessingFile) (*configmapping.ConfigMapping, []*configtask.ConfigTask, error) {
	ctx := context.Background()

	// 1. Config mapping
	cfgMapping, err := cl.cfgMappingService.FindByClientID(ctx, file.ClientID)
	if err != nil {
		return nil, nil, err
	}

	// 2. Config task
	cfgTasks, err := cl.cfgTaskService.FindByConfigMappingID(ctx, int32(cfgMapping.ID))
	if err != nil {
		return nil, nil, err
	}

	return cfgMapping, cfgTasks, nil
}

// toFileParameters ... return map[key]value
func toFileParameters(rawJson string) map[string]string {
	headerMap, _ := converter.StringToMap("fileParameters", rawJson, true)
	return headerMap
}

// toConfigTaskMD ...
func toConfigTaskMD(task configtask.ConfigTask) (ConfigTaskMD, error) {
	// 1. convert Header
	headerMap, err := converter.StringToMap("header", task.Header, true)
	if err != nil {
		return ConfigTaskMD{}, err
	}

	// 2. convert Request Params
	requestParamsMap, err := toMapRequestFieldMD(task.TaskIndex, "configRequestParam", task.RequestParams)
	if err != nil {
		return ConfigTaskMD{}, err
	}

	// 2. convert Request Body
	requestBodyMap, err := toMapRequestFieldMD(task.TaskIndex, "configRequestBody", task.RequestBody)
	if err != nil {
		return ConfigTaskMD{}, err
	}

	// 3. convert ResponseMD
	responseMD, err := toResponseMD(task)
	if err != nil {
		return ConfigTaskMD{}, err
	}

	// 4. Return result
	return ConfigTaskMD{
		TaskIndex: int(task.TaskIndex),
		// Request
		Endpoint: task.EndPoint,
		Method:   task.Method,
		Header:   headerMap,
		//RequestParams: requestParamsMap,
		RequestParamsMap: requestParamsMap,
		RequestBodyMap:   requestBodyMap,
		RequestParams:    map[string]interface{}{},
		RequestBody:      map[string]interface{}{},
		// Response
		Response: responseMD,
	}, nil
}

// toMapRequestFieldMD ... return map[fieldName]RequestFieldMD, error
func toMapRequestFieldMD(taskID int32, dataName string, dataRaw string) (map[string]*RequestFieldMD, error) {
	// 1. Check empty
	if len(dataRaw) == 0 {
		return map[string]*RequestFieldMD{}, nil
	}

	// 2. Unmarshal
	var list []*RequestFieldMD
	if err := json.Unmarshal([]byte(dataRaw), &list); err != nil {
		logger.Errorf("error when convert %v: value=%v, err=%v", dataName, dataRaw, err)
		return nil, fmt.Errorf("cannot convert %v", dataName)
	}

	// 3. Convert to Map
	fieldMap := map[string]*RequestFieldMD{}
	for _, reqField := range list {
		requestFieldEnriched, err := enrichRequestFieldMD(taskID, *reqField)
		if err != nil {
			return nil, fmt.Errorf("cannot get config %v", dataName)
		}

		fieldMap[reqField.Field] = &requestFieldEnriched
	}

	// 4. return
	return fieldMap, nil
}

// enrichRequestFieldMD ... enrich more data
func enrichRequestFieldMD(taskId int32, reqField RequestFieldMD) (RequestFieldMD, error) {
	fieldName := reqField.Field
	valuePattern := reqField.ValuePattern

	if strings.HasPrefix(valuePattern, prefixMappingRequest) {
		if len(valuePattern) == 2 {
			columnIndex := string(valuePattern[1]) // if `$A` -> columnIndex = `A`
			logger.Infof("----- task %v, field %v is mapping with column %v", taskId, fieldName, columnIndex)
			reqField.ValueDependsOn = ValueDependsOnExcel
			reqField.ValueDependsOnKey = columnIndex
		} else if len(valuePattern) > len(prefixMappingRequestResponse)+2 && strings.HasPrefix(valuePattern, prefixMappingRequestResponse) {
			template := strings.TrimPrefix(valuePattern, prefixMappingRequestResponse) // $response1.field_abc -> template = 1.field_abc
			dependOnTaskId, err := strconv.Atoi(string(template[0]))                   // 1.field_abc -> 1
			if err != nil || template[1] != '.' {
				logger.Infof("----- task %v, field %v has invalid value is %v", taskId, fieldName, valuePattern)
				return RequestFieldMD{}, fmt.Errorf("mapping request is invalid: %v", valuePattern)
			}

			responsePath := template[2:] // 1.field_abc -> field_abc
			reqField.ValueDependsOn = ValueDependsOnTask
			reqField.ValueDependsOnKey = responsePath
			reqField.ValueDependsOnTaskID = dependOnTaskId
		} else if len(valuePattern) > len(prefixMappingRequestParameter)+2 && strings.HasPrefix(valuePattern, prefixMappingRequestParameter) {
			template := strings.TrimPrefix(valuePattern, prefixMappingRequestParameter) // $param.field_abc -> paramKey = .field_abc
			if len(template) <= 1 || template[0] != '.' {
				logger.Infof("----- task %v, field %v has invalid value is %v", taskId, fieldName, valuePattern)
				return RequestFieldMD{}, fmt.Errorf("mapping request is invalid: %v", valuePattern)
			}

			reqField.ValueDependsOn = ValueDependsOnParam
			reqField.ValueDependsOnKey = template[1:]
		} else {
			logger.Errorf("----- task %v, field %v has invalid value is %v", taskId, fieldName, valuePattern)
			return RequestFieldMD{}, fmt.Errorf("mapping request is invalid: %v", valuePattern)
		}
	} else {
		reqField.Value = valuePattern // raw data
		reqField.ValueDependsOn = ValueDependsOnNone
	}

	return reqField, nil
}

// toResponseMD ...
func toResponseMD(task configtask.ConfigTask) (ResponseMD, error) {
	// 1. ResponseCode
	var responseCode ResponseCode
	if err := json.Unmarshal([]byte(task.ResponseSuccessCodeSchema), &responseCode); err != nil {
		logger.Errorf("error when convert ResponseSuccessCodeSchema: value=%v, err=%v", task.ResponseSuccessCodeSchema, err)
		return ResponseMD{}, fmt.Errorf("cannot convert ResponseSuccessCodeSchema")
	}

	// 2. ResponseMsg
	var responseMsg ResponseMsg
	if err := json.Unmarshal([]byte(task.ResponseMessageSchema), &responseMsg); err != nil {
		logger.Errorf("error when convert ResponseMessageSchema: value=%v, err=%v", task.ResponseMessageSchema, err)
		return ResponseMD{}, fmt.Errorf("cannot convert ResponseMessageSchema")
	}

	// 3. Response
	responseMD := ResponseMD{
		HttpStatusSuccess: task.ResponseSuccessHTTPStatus,
		Code:              responseCode,
		Message:           responseMsg,
	}
	return responseMD, nil
}
