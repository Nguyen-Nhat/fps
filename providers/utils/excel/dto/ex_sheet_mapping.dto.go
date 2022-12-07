package dto

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
	"strconv"
	"strings"
)

type SheetMappingMetadata struct {
	TaskId   CellData[int]
	Endpoint CellData[string]
	Header   CellData[string]
	Request  CellData[string]
	Response CellData[string]
}

var _ Converter[SheetMappingMetadata, MappingRow] = &SheetMappingMetadata{}

// GetHeaders implement from Converter interface
func (f *SheetMappingMetadata) GetHeaders() []string {
	var headers []string
	if len(f.TaskId.ColumnName) > 0 {
		headers = append(headers, f.TaskId.ColumnName)
	}
	if len(f.Endpoint.ColumnName) > 0 {
		headers = append(headers, f.Endpoint.ColumnName)
	}
	if len(f.Header.ColumnName) > 0 {
		headers = append(headers, f.Header.ColumnName)
	}
	if len(f.Request.ColumnName) > 0 {
		headers = append(headers, f.Request.ColumnName)
	}
	if len(f.Response.ColumnName) > 0 {
		headers = append(headers, f.Response.ColumnName)
	}
	return headers
}

// ToOutput implement from Converter interface
func (f *SheetMappingMetadata) ToOutput(rowId int) (MappingRow, error) {
	noData := MappingRow{}
	// 1. convert Header
	headerMap, err := converter.StringToMap(f.Header.ColumnName, f.Header.Value, !f.Header.Constrains.IsRequired)
	if err != nil {
		return noData, err
	}

	// 2. convert Request
	requestMap, err := converter.StringToMap(f.Request.ColumnName, f.Request.Value, !f.Request.Constrains.IsRequired)
	if err != nil {
		return noData, err
	}
	requestMapConverted, err := convertRequestMapping(f.TaskId.Value, requestMap)
	if err != nil {
		return noData, err
	}

	// 3. covert Response
	response, err := converter.StringJsonToStruct(f.Response.ColumnName, f.Response.Value, MappingResponse{})
	if err != nil {
		return noData, err
	}
	// todo validate response

	return MappingRow{
		RowId:    rowId,
		TaskId:   f.TaskId.Value,
		Endpoint: f.Endpoint.Value, // todo validate endpoint format
		Header:   headerMap,
		Request:  requestMapConverted,
		Response: *response,
	}, nil
}

// GetMetadata implement from Converter interface
func (f *SheetMappingMetadata) GetMetadata() *SheetMappingMetadata {
	return f
}

// ---------------------------------------------------------------------------------------------------------------------

type (
	MappingRow struct {
		RowId    int
		TaskId   int
		Endpoint string
		Header   map[string]string
		Request  map[string]MappingRequest // map[field_name]MappingRequest
		Response MappingResponse
	}

	MappingResponse struct {
		HttpStatusSuccess string         `json:"httpStatusSuccess"`
		Code              MappingResCode `json:"code"`
		Message           MappingResMsg  `json:"message"`
	}

	MappingResCode struct {
		Path          string `json:"path"`
		SuccessValues string `json:"successValues"`
	}

	MappingResMsg struct {
		Path string `json:"path"`
	}

	MappingRequest struct {
		FieldName                string
		Value                    string
		IsMappingExcel           bool
		IsMappingResponse        bool
		MappingKey               string // is in {A,B,C} when isMappingExcel=true, is jsonPath when isMappingResponse=true
		DependOnResponseOfTaskId int
	}
)

const (
	prefixMappingRequest         = "$"
	prefixMappingRequestResponse = "$response"
)

func convertRequestMapping(taskId int, requestMap map[string]string) (map[string]MappingRequest, error) {
	result := make(map[string]MappingRequest)
	for fieldName, valueMapping := range requestMap {
		mappingRequest := MappingRequest{FieldName: fieldName}
		if strings.HasPrefix(valueMapping, prefixMappingRequest) {
			if len(valueMapping) == 2 {
				columnIndex := string(valueMapping[1])
				logger.Infof("----- task %v, field %v is mapping with column %v", taskId, fieldName, columnIndex)
				mappingRequest.IsMappingExcel = true
				mappingRequest.MappingKey = columnIndex
			} else if len(valueMapping) > len(prefixMappingRequestResponse)+2 && strings.HasPrefix(valueMapping, prefixMappingRequestResponse) {
				template := strings.TrimPrefix(valueMapping, prefixMappingRequestResponse)
				dependOnTaskId, err := strconv.Atoi(string(template[0]))
				if err != nil || template[1] != '.' {
					logger.Infof("----- task %v, field %v has invalid value is %v", taskId, fieldName, valueMapping)
					return nil, fmt.Errorf("mapping request is invalid: %v", valueMapping)
				}
				responsePath := template[2:]
				mappingRequest.IsMappingResponse = true
				mappingRequest.MappingKey = responsePath
				mappingRequest.DependOnResponseOfTaskId = dependOnTaskId
			} else {
				logger.Errorf("----- task %v, field %v has invalid value is %v", taskId, fieldName, valueMapping)
				return nil, fmt.Errorf("mapping request is invalid: %v", valueMapping)
			}
		} else {
			mappingRequest.Value = valueMapping
		}
		result[fieldName] = mappingRequest
	}
	return result, nil
}
