package flatten

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

const (
	errFileCannotLoad fileprocessing.ErrorDisplay = "không thể tải file import"
	errFileInvalid    fileprocessing.ErrorDisplay = "file tải lên không đúng định dạng"
	errFileNoData     fileprocessing.ErrorDisplay = "file tải lên không có dữ liệu"
	// error row
	errRowMissingDataColumn = "thiếu dữ liệu cột"
	// error config
	errConfigMapping      = "lỗi cấu hình hệ thống"
	errConfigMissingParam = "thiếu cấu hình hệ thống cho"
	// error data type
	errTypeWrong      = "sai kiểu dữ liệu"
	errTypeNotSupport = "không hỗ trợ kiểu dữ liệu"
)

// validateImportingData ...
// check: empty, invalid data type, constrains
func validateImportingData(sheetData [][]string, cfgMapping configloader.ConfigMappingMD) ([]configloader.ConfigMappingMD, []ErrorRow, error) {
	// 1. Empty or no data (at start row)
	dataStartAt := cfgMapping.DataStartAtRow
	if dataStartAt <= 1 { // must >= 2
		logger.ErrorT("%v: DataStartAtRow = %s", errConfigMapping, dataStartAt)
		return nil, nil, errors.New(errConfigMapping)
	}
	if len(sheetData) == 0 || len(sheetData) < dataStartAt {
		return nil, nil, errors.New(string(errFileNoData))
	}

	// 2. Validate
	var errorRows []ErrorRow
	var configMappings []configloader.ConfigMappingMD

	for rowID := dataStartAt - 1; rowID < len(sheetData); rowID++ {
		cfgMappingWithConvertedData, errorRowsInRow := validateImportingDataRowAndCloneConfigMapping(rowID, sheetData[rowID], cfgMapping)

		// check error rows
		if len(errorRowsInRow) > 0 {
			errorRows = append(errorRows, errorRowsInRow...)
		} else {
			configMappings = append(configMappings, cfgMappingWithConvertedData)
		}
	}

	return configMappings, errorRows, nil
}

// validateImportingDataRowAndCloneConfigMapping ...
func validateImportingDataRowAndCloneConfigMapping(rowID int, rowData []string, configMapping configloader.ConfigMappingMD) (configloader.ConfigMappingMD, []ErrorRow) {
	var errorRows []ErrorRow

	// 1. Get value for each RequestField in each Task
	var tasksUpdated []configloader.ConfigTaskMD
	for _, orgTask := range configMapping.Tasks {
		task := orgTask.Clone()
		// 1.1. RequestField in Request Params
		for fieldName, reqField := range task.RequestParamsMap {
			// 1.1.1. Get value in String type
			var valueStr string
			switch reqField.ValueDependsOn {
			case configloader.ValueDependsOnExcel:
				cellValue, errorRowsExel := validateAndGetValueForRequestFieldExcel(rowID, rowData, reqField)
				if len(errorRowsExel) == 0 {
					valueStr = cellValue
				} else {
					errorRows = append(errorRows, errorRowsExel...)
				}
			case configloader.ValueDependsOnParam:
				paramValue, errorRowsParams := validateAndGetValueForFieldParam(rowID, reqField, configMapping.FileParameters)
				if len(errorRowsParams) == 0 {
					valueStr = paramValue
				} else {
					errorRows = append(errorRows, errorRowsParams...)
				}
			case configloader.ValueDependsOnNone:
				valueStr = reqField.Value
			default:
				// No support ValueDependsOnTask
				// Because data of task only is gotten after call api to provider => cannot get at this time
				continue // go to next reqField
			}

			// 1.1.2. Get real value
			realValue, errMsg := convertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
			if len(errMsg) > 0 {
				errorRows = append(errorRows, ErrorRow{rowID, errMsg})
			} else {
				task.RequestParams[reqField.Field] = realValue
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestParamsMap, fieldName)
			}
		}

		// 1.2. RequestField in Request Body
		for fieldName, reqField := range task.RequestBodyMap {
			// 1.2.1. Get value in String type
			var valueStr string
			switch reqField.ValueDependsOn {
			case configloader.ValueDependsOnExcel:
				cellValue, errorRowsExel := validateAndGetValueForRequestFieldExcel(rowID, rowData, reqField)
				if len(errorRowsExel) == 0 {
					errorRows = append(errorRows, errorRowsExel...)
					valueStr = cellValue
				}
			case configloader.ValueDependsOnParam:
				paramValue, errorRowsParams := validateAndGetValueForFieldParam(rowID, reqField, configMapping.FileParameters)
				if len(errorRowsParams) == 0 {
					errorRows = append(errorRows, errorRowsParams...)
					valueStr = paramValue
				}
			case configloader.ValueDependsOnNone:
				valueStr = reqField.Value
			default:
				// No support ValueDependsOnTask
				// Because data of task only is gotten after call api to provider => cannot get at this time
				continue // go to next reqField
			}

			// 1.2.2. Get real value
			realValue, errMsg := convertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
			if len(errMsg) > 0 {
				errorRows = append(errorRows, ErrorRow{rowID, errMsg})
			} else {
				task.RequestBody[reqField.Field] = realValue
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestBodyMap, fieldName)
			}
		}

		// 1.3. Set value for remaining data
		task.ImportRowData = rowData
		task.ImportRowIndex = rowID
		tasksUpdated = append(tasksUpdated, task)
	}

	// 2. Update Tasks with updated value
	configMapping.Tasks = tasksUpdated

	// 3. Return
	return configMapping, errorRows
}

func validateAndGetValueForFieldParam(rowID int, reqField *configloader.RequestFieldMD, fileParameters map[string]string) (string, []ErrorRow) {
	var errorRows []ErrorRow
	paramKey := reqField.ValueDependsOnKey

	// Validate Require
	paramValue, existed := fileParameters[paramKey]
	if !existed {
		reason := fmt.Sprintf("%s %s", errConfigMissingParam, paramKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
	}

	// Validate Data type
	paramValue = strings.TrimSpace(paramValue)
	// todo re-check, because Validating Data also works in convertToRealValue() function
	//if isWrongDataType(reqField.Type, paramValue) {
	//	reason := fmt.Sprintf("%s %s", errConfigWrongType, paramKey)
	//	errorRow := ErrorRow{RowId: rowID, Reason: reason}
	//	errorRows = append(errorRows, errorRow)
	//}

	reqField.Value = paramValue

	return paramValue, errorRows
}

func validateAndGetValueForRequestFieldExcel(rowID int, rowData []string, reqField *configloader.RequestFieldMD) (string, []ErrorRow) {
	var errorRows []ErrorRow
	columnKey := reqField.ValueDependsOnKey
	columnIndex := int(strings.ToUpper(columnKey)[0]) - int('A') // get first character then

	// Validate Require
	if columnIndex >= len(rowData) || // column request out of range
		(reqField.Required && len(strings.TrimSpace(rowData[columnIndex])) == 0) { // column is required by value is empty
		reason := fmt.Sprintf("%s %s", errRowMissingDataColumn, columnKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
	}

	// Validate Data type
	cellValue := strings.TrimSpace(rowData[columnIndex])
	// todo re-check
	//if isWrongDataType(reqField.Type, cellValue) {
	//	reason := fmt.Sprintf("%s %s", errRowWrongType, columnKey)
	//	errorRow := ErrorRow{RowId: rowID, Reason: reason}
	//	errorRows = append(errorRows, errorRow)
	//}

	reqField.Value = cellValue

	return cellValue, errorRows
}

// isWrongDataType ...
// supporting type: Int, ...
func isWrongDataType(fieldType string, valueStr string) bool {
	switch fieldType {
	case configloader.TypeInt:
		_, err := strconv.Atoi(valueStr)
		if err != nil {
			return true
		}
	}
	return false
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
