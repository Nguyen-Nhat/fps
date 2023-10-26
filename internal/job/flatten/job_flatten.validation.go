package flatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
)

const (
	errFileCannotLoad fileprocessing.ErrorDisplay = "không thể tải file import"
	errFileInvalid    fileprocessing.ErrorDisplay = "file tải lên không đúng định dạng"
	errFileNoData     fileprocessing.ErrorDisplay = "file tải lên không có dữ liệu"
	// error row
	errRowMissingDataColumn = "thiếu dữ liệu cột"
	errRowMissingData       = "không có dữ liệu"
	// error config
	errConfigMapping      = "lỗi cấu hình hệ thống"
	errConfigMissingParam = "thiếu cấu hình hệ thống cho"
)

var regexDoubleBrace = regexp.MustCompile(`\{\{(.*?)}}`)

// validateImportingData ...
// check: empty, invalid data type, constrains
func validateImportingData(sheetData [][]string, cfgMapping configloader.ConfigMappingMD) ([]*configloader.ConfigMappingMD, []ErrorRow, error) {
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
	var configMappings []*configloader.ConfigMappingMD

	for id := dataStartAt - 1; id < len(sheetData); id++ {
		rowID := id - dataStartAt + 1 // rowID is index of data (not include header), start from 1
		cfgMappingWithConvertedData, errorRowsInRow := validateImportingDataRowAndCloneConfigMapping(rowID, sheetData[id], cfgMapping)

		// check error rows
		if len(errorRowsInRow) > 0 {
			errorRows = append(errorRows, errorRowsInRow...)
		} else {
			configMappings = append(configMappings, &cfgMappingWithConvertedData)
		}
	}

	return configMappings, errorRows, nil
}

// validateImportingDataRowAndCloneConfigMapping ...
func validateImportingDataRowAndCloneConfigMapping(rowID int, rowData []string, configMapping configloader.ConfigMappingMD) (configloader.ConfigMappingMD, []ErrorRow) {
	var errorRows []ErrorRow

	// 0. Check row empty
	if len(rowData) == 0 {
		return configloader.ConfigMappingMD{}, []ErrorRow{{RowId: rowID, Reason: errRowMissingData}}
	}

	// 1. Get value for each RequestField in each Task
	var tasksUpdated []configloader.ConfigTaskMD
	fileParameters := configMapping.FileParameters
	for _, orgTask := range configMapping.Tasks {
		task := orgTask.Clone()
		// 1.1. RequestField in Request Params
		for fieldName, reqField := range task.RequestParamsMap {
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(rowID, rowData, reqField, fileParameters)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.1.2. Get real value
			realValue, err := basejobmanager.ConvertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
			if err != nil {
				errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
			} else {
				if realValue != nil {
					task.RequestParams[reqField.Field] = realValue
				}
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestParamsMap, fieldName)
			}
		}

		// 1.2. RequestField in Request Body (support ArrayItem)
		for fieldName, reqField := range task.RequestBodyMap {
			// 1.2.1. Validate ArrayItemMap
			if len(reqField.ArrayItemMap) > 0 {
				arrayItemMapUpdated, childMap, errorRowsForArrayItem := validateArrayItemMap(rowID, rowData, reqField.ArrayItemMap, fileParameters)
				if len(errorRowsForArrayItem) > 0 {
					errorRows = append(errorRows, errorRowsForArrayItem...)
					continue
				}

				reqField.ArrayItemMap = arrayItemMapUpdated
				task.RequestBody[fieldName] = []map[string]interface{}{childMap}
				if len(childMap) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
					val, ok := childMap[""]
					if ok && len(fmt.Sprintf("%+v", val)) > 0 { // val is not empty value -> set field is array[val]
						task.RequestBody[fieldName] = []interface{}{val}
					} else if ok && len(fmt.Sprintf("%+v", val)) == 0 { // val is empty -> remove field
						delete(task.RequestBody, fieldName)
					}
				}
				if len(arrayItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove task.RequestBodyMap field value
					delete(task.RequestBodyMap, fieldName)
				}
			}

			// 1.2.2. Validate field
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(rowID, rowData, reqField, fileParameters)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.2.2. Get real value
			realValue, err := basejobmanager.ConvertToRealValue(reqField.Type, valueStr, reqField.ValueDependsOnKey)
			if err != nil {
				errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
			} else {
				if realValue != nil {
					task.RequestBody[reqField.Field] = realValue
				}
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestBodyMap, fieldName)
			}
		}

		// 1.3. Validate ResponseCode config
		resultAfterMatch, errorRow := validateAndMatchJsonPath(rowID, rowData, task.Response.Code.MustHaveValueInPath)
		if errorRow != nil {
			errorRows = append(errorRows, *errorRow)
		} else {
			task.Response.Code.MustHaveValueInPath = resultAfterMatch
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

func validateArrayItemMap(rowID int, rowData []string, arrayItemMap map[string]*configloader.RequestFieldMD, fileParameters map[string]interface{}) (
	map[string]*configloader.RequestFieldMD, map[string]interface{}, []ErrorRow) {
	// 1. Init value
	var errorRows []ErrorRow
	childMap := make(map[string]interface{})

	// 2. Explore each item in array
	for fieldNameChild, reqFieldChild := range arrayItemMap {
		// 2.1. Case field type is Array
		if len(reqFieldChild.ArrayItemMap) > 0 {
			// 2.1.1. Validate
			arrayItemMapUpdated, childMapInArr, errorRowsForArrayItem := validateArrayItemMap(rowID, rowData, reqFieldChild.ArrayItemMap, fileParameters)
			if len(errorRowsForArrayItem) > 0 {
				errorRows = append(errorRows, errorRowsForArrayItem...)
				continue
			}

			// 2.1.2. Update childMap and remove field if it was mapped data
			reqFieldChild.ArrayItemMap = arrayItemMapUpdated
			childMap[fieldNameChild] = []map[string]interface{}{childMapInArr}
			if len(childMapInArr) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
				val, ok := childMapInArr[""]
				if ok && len(fmt.Sprintf("%+v", val)) > 0 { // val is not empty value -> set field is array[val]
					childMap[fieldNameChild] = []interface{}{val}
				} else if ok && len(fmt.Sprintf("%+v", val)) == 0 { // val is empty -> remove field
					delete(childMap, fieldNameChild)
				}
			}
			if len(arrayItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove that field in arrayItemMap
				delete(arrayItemMap, fieldNameChild)
			}

			// 2.1.3. Continue
			continue
		}

		// 2.2. Case field type is Object
		if len(reqFieldChild.ItemsMap) > 0 {
			// 2.2.1. Validate
			objectItemMapUpdated, childMapInObj, errorRowsForArrayItem := validateArrayItemMap(rowID, rowData, reqFieldChild.ItemsMap, fileParameters)
			if len(errorRowsForArrayItem) > 0 {
				errorRows = append(errorRows, errorRowsForArrayItem...)
				continue
			}

			// 2.2.2. Update childMap and remove field if it was mapped data
			reqFieldChild.ArrayItemMap = objectItemMapUpdated
			childMap[fieldNameChild] = childMapInObj
			if len(objectItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove that field in arrayItemMap
				delete(arrayItemMap, fieldNameChild)
			}

			// 2.2.3. Continue
			continue
		}

		// 2.3. In Normal case, field maybe int, string, ... -> need to get value (string) from config
		valueChildStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(rowID, rowData, reqFieldChild, fileParameters)
		if len(errorRowsAfterGet) > 0 {
			errorRows = append(errorRows, errorRowsAfterGet...)
			continue
		}
		if isByPassField { // no need to convert
			continue
		}

		// 2.4. Get real value from string value
		realValueChild, err := basejobmanager.ConvertToRealValue(reqFieldChild.Type, valueChildStr, reqFieldChild.ValueDependsOnKey)
		if err != nil {
			errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
		} else {
			if realValueChild != nil {
				childMap[reqFieldChild.Field] = realValueChild
			}
			// config will be converted to Json string, then save to DB -> delete to reduce size of json string
			delete(arrayItemMap, fieldNameChild)
		}
	}

	// 3. Return
	return arrayItemMap, childMap, errorRows
}

func getValueStrByRequestFieldMD(rowID int, rowData []string, reqField *configloader.RequestFieldMD, fileParameters map[string]interface{}) (string, bool, []ErrorRow) {
	// 1. If type is array or object, not get value
	if reqField.Type == configloader.TypeArray || reqField.Type == configloader.TypeObject {
		return "", true, nil
	}

	// 2. Get value in String type
	var valueStr string
	var errorRows []ErrorRow
	isByPassField := false
	switch reqField.ValueDependsOn {
	case configloader.ValueDependsOnExcel:
		cellValue, errorRowsExel := validateAndGetValueForRequestFieldExcel(rowID, rowData, reqField)
		if len(errorRowsExel) != 0 {
			errorRows = append(errorRows, errorRowsExel...)
		} else {
			valueStr = cellValue
		}
	case configloader.ValueDependsOnParam:
		paramValue, errorRowsParams := validateAndGetValueForFieldParam(rowID, reqField, fileParameters)
		if len(errorRowsParams) != 0 {
			errorRows = append(errorRows, errorRowsParams...)
		} else {
			valueStr = paramValue
		}
	case configloader.ValueDependsOnNone:
		valueStr = reqField.Value
	case configloader.ValueDependsOnTask:
		isByPassField = true // If value depends on Previous Task -> not get value
		valueDependsOnKeyMatched, errorRow := validateAndMatchJsonPath(rowID, rowData, reqField.ValueDependsOnKey)
		if errorRow != nil {
			errorRows = append(errorRows, *errorRow)
		} else {
			reqField.ValueDependsOnKey = valueDependsOnKeyMatched
		}
		valueStr = reqField.Value
	default:
		errMsg := fmt.Sprintf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
		errorRows = append(errorRows, ErrorRow{rowID, errMsg})

	}
	return valueStr, isByPassField, errorRows
}

func validateAndGetValueForFieldParam(rowID int, reqField *configloader.RequestFieldMD, fileParameters map[string]interface{}) (string, []ErrorRow) {
	var errorRows []ErrorRow
	paramKey := reqField.ValueDependsOnKey

	// Validate Require
	paramValue, existed := fileParameters[paramKey]
	if !existed {
		reason := fmt.Sprintf("%s %s", errConfigMissingParam, paramKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
	}

	// convert param to string
	paramValueStr := ""
	if reqField.Type == configloader.TypeJson {
		jsonStr, _ := json.Marshal(paramValue)
		paramValueStr = string(jsonStr)
	} else {
		paramValueStr = fmt.Sprintf("%v", paramValue)
	}

	// check required
	if reqField.Required && len(paramValueStr) == 0 {
		reason := fmt.Sprintf("%s %s", errRowMissingData, paramKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
		return "", errorRows
	}

	reqField.Value = paramValueStr
	return paramValueStr, errorRows
}

func validateAndGetValueForRequestFieldExcel(rowID int, rowData []string, reqField *configloader.RequestFieldMD) (string, []ErrorRow) {
	var errorRows []ErrorRow
	columnKey := reqField.ValueDependsOnKey
	columnIndex := int(strings.ToUpper(columnKey)[0]) - int('A') // get first character then

	// Validate Require
	if reqField.Required &&
		(columnIndex >= len(rowData) || // column request out of range
			len(strings.TrimSpace(rowData[columnIndex])) == 0) { // column is required by value is empty
		reason := fmt.Sprintf("%s %s", errRowMissingDataColumn, columnKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
		return "", errorRows
	}

	cellValue := ""
	if columnIndex < len(rowData) { // if out of range -> default is empty
		cellValue = strings.TrimSpace(rowData[columnIndex])
	}

	reqField.Value = cellValue

	return cellValue, errorRows
}

// validateAndMatchJsonPath ... support validate json path and update json path if it contains variable
// for example:
//   - Json path: data.transactions.#(name=="{{ $A }}").id
//   - Excel data has $A = quy
//     -> output: data.transactions.#(name=="quy").id
func validateAndMatchJsonPath(rowID int, rowData []string, jsonPath string) (string, *ErrorRow) {
	// 1. Extract data with format like `{{ $A }}`
	matchers := regexDoubleBrace.FindStringSubmatch(jsonPath)

	// 2. Return if not match
	if len(matchers) != 2 {
		return jsonPath, nil
	}

	// 3. Validate and Replace value
	{
		valuePatternWithDoubleBrace := matchers[0]
		valuePattern := strings.TrimSpace(matchers[1])

		// 3.1. Get column key: $A -> A
		if !strings.HasPrefix(valuePattern, dto.PrefixMappingRequest) || len(valuePattern) != 2 {
			errorRow := ErrorRow{RowId: rowID, Reason: errConfigMapping}
			logger.Errorf("validateResponseCode ... error %s -> %s", errConfigMapping, valuePatternWithDoubleBrace)
			return "", &errorRow
		}
		columnKey := string(valuePattern[1]) // if `$A` -> columnIndex = `A`
		columnIndex := int(strings.ToUpper(columnKey)[0]) - int('A')

		// 3.2. Validate value
		if columnIndex >= len(rowData) || // column request out of range
			len(strings.TrimSpace(rowData[columnIndex])) == 0 { // column is required by value is empty
			reason := fmt.Sprintf("%s %s", errRowMissingDataColumn, columnKey)
			errorRow := ErrorRow{RowId: rowID, Reason: reason}
			logger.Errorf("validateResponseCode ... error %+v", reason)
			return "", &errorRow
		}

		// 3.3. Replace value
		cellValue := strings.TrimSpace(rowData[columnIndex])
		jsonPath = strings.ReplaceAll(jsonPath, valuePatternWithDoubleBrace, cellValue)
	}

	// 4. return
	return jsonPath, nil
}
