package flatten

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

// defaultFileHeaderAtRowID ... this value is hardcode, please move it to config_mapping in the future
const defaultFileHeaderAtRowID = 0

var regexDoubleBrace = regexp.MustCompile(`\{\{(.*?)}}`)

// validateImportingData ...
// check: empty, invalid data type, constrains
func validateImportingData(ctx context.Context, sheetData [][]string, cfgMapping configloader.ConfigMappingMD) ([]*configloader.ConfigMappingMD, []ErrorRow, error) {
	// 1. Empty or no data (at start row)
	dataStartAt := cfgMapping.DataStartAtRow
	if dataStartAt <= 1 { // must >= 2
		msg := i18n.GetMessageCtx(ctx, "errConfigMapping")
		logger.ErrorT("%v: DataStartAtRow = %s", msg, dataStartAt)
		return nil, nil, errors.New(msg)
	}
	if len(sheetData) == 0 || len(sheetData) < dataStartAt {
		return []*configloader.ConfigMappingMD{}, nil, nil
	}
	fileHeader := sheetData[defaultFileHeaderAtRowID]

	// 2. Validate
	var errorRows []ErrorRow
	var configMappings []*configloader.ConfigMappingMD

	for id := dataStartAt - 1; id < len(sheetData); id++ {
		rowID := id - dataStartAt + 1 // rowID is index of data (not include header), start from 1
		cfgMappingWithConvertedData, errorRowsInRow := validateImportingDataRowAndCloneConfigMapping(ctx, rowID, fileHeader, sheetData[id], cfgMapping)

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
func validateImportingDataRowAndCloneConfigMapping(ctx context.Context, rowID int, fileHeader []string, rowData []string, configMapping configloader.ConfigMappingMD) (configloader.ConfigMappingMD, []ErrorRow) {
	var errorRows []ErrorRow

	// 0. Check row empty
	if len(rowData) == 0 {
		return configloader.ConfigMappingMD{}, []ErrorRow{{RowId: rowID, Reason: i18n.GetMessageCtx(ctx, "errRowNoData")}}
	}

	// 1. Get value for each RequestField in each Task
	var tasksUpdated []configloader.ConfigTaskMD
	fileParameters := configMapping.FileParameters
	for _, orgTask := range configMapping.Tasks {
		task := orgTask.Clone()
		// 1.1. RequestField in Request Header
		for fieldName, reqField := range task.RequestHeaderMap {
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(ctx, rowID, rowData, reqField, fileParameters, fileHeader)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.1.2. Get real value
			valueDependsOnKey := getValueDependOnKeyExcel(reqField.ValueDependsOn, reqField.ValueDependsOnKey, fileHeader)
			realValue, err := basejobmanager.ConvertToRealValue(ctx, reqField.Type, valueStr, valueDependsOnKey)
			if err != nil {
				errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
			} else {
				if realValue != nil {
					task.RequestHeader[reqField.Field] = realValue
				}
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestHeaderMap, fieldName)
			}
		}

		// 1.2. RequestField in Path Params
		for fieldName, reqField := range task.PathParamsMap {
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(ctx, rowID, rowData, reqField, fileParameters, fileHeader)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.2.2. Get real value
			valueDependsOnKey := getValueDependOnKeyExcel(reqField.ValueDependsOn, reqField.ValueDependsOnKey, fileHeader)
			realValue, err := basejobmanager.ConvertToRealValue(ctx, reqField.Type, valueStr, valueDependsOnKey)
			if err != nil {
				errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
			} else {
				if realValue != nil &&
					len(fmt.Sprintf("%+v", realValue)) > 0 { // ignore case empty value
					task.PathParams[reqField.Field] = realValue
				}
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.PathParamsMap, fieldName)
			}
		}

		// 1.3. RequestField in Request Params
		for fieldName, reqField := range task.RequestParamsMap {
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(ctx, rowID, rowData, reqField, fileParameters, fileHeader)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.3.1. Get real value
			valueDependsOnKey := getValueDependOnKeyExcel(reqField.ValueDependsOn, reqField.ValueDependsOnKey, fileHeader)
			realValue, err := basejobmanager.ConvertToRealValue(ctx, reqField.Type, valueStr, valueDependsOnKey)
			if err != nil {
				errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
			} else {
				if realValue != nil &&
					len(fmt.Sprintf("%+v", realValue)) > 0 { // ignore case empty value
					task.RequestParams[reqField.Field] = realValue
				}
				// config will be converted to Json string, then save to DB -> delete to reduce size of json string
				delete(task.RequestParamsMap, fieldName)
			}
		}

		// 1.4. RequestField in Request Body (support ArrayItem)
		for fieldName, reqField := range task.RequestBodyMap {
			// 1.4.1. Validate ArrayItemMap
			if len(reqField.ArrayItemMap) > 0 {
				arrayItemMapUpdated, childMap, errorRowsForArrayItem := validateArrayItemMap(ctx, rowID, rowData, reqField.ArrayItemMap, fileParameters, fileHeader)
				if len(errorRowsForArrayItem) > 0 {
					errorRows = append(errorRows, errorRowsForArrayItem...)
					continue
				}

				reqField.ArrayItemMap = arrayItemMapUpdated
				task.RequestBody[fieldName] = []map[string]interface{}{childMap}
				if len(childMap) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
					if val, ok := childMap[""]; ok {
						task.RequestBody[fieldName] = []interface{}{val}
					}
				} else if len(childMap) == 0 { // childMap is empty -> remove field
					delete(task.RequestBody, fieldName)
				}

				if len(arrayItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove task.RequestBodyMap field value
					delete(task.RequestBodyMap, fieldName)
				}
			}

			// 1.4.2. Validate field
			valueStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(ctx, rowID, rowData, reqField, fileParameters, fileHeader)
			if len(errorRowsAfterGet) > 0 {
				errorRows = append(errorRows, errorRowsAfterGet...)
				continue
			}
			if isByPassField { // no need to convert
				continue
			}

			// 1.4.2. Get real value
			valueDependsOnKey := getValueDependOnKeyExcel(reqField.ValueDependsOn, reqField.ValueDependsOnKey, fileHeader)
			realValue, err := basejobmanager.ConvertToRealValue(ctx, reqField.Type, valueStr, valueDependsOnKey)
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

		// 1.5. Validate ResponseCode config
		resultAfterMatch, errorRow := validateAndMatchJsonPath(ctx, rowID, rowData, task.Response.Code.MustHaveValueInPath, fileHeader)
		if errorRow != nil {
			errorRows = append(errorRows, *errorRow)
		} else {
			task.Response.Code.MustHaveValueInPath = resultAfterMatch
		}

		// 1.6. Set value for remaining data
		task.ImportRowHeader = fileHeader
		task.ImportRowData = rowData
		task.ImportRowIndex = rowID
		tasksUpdated = append(tasksUpdated, task)
	}

	// 2. Update Tasks with updated value
	configMapping.Tasks = tasksUpdated

	// 3. Return
	return configMapping, errorRows
}

func getValueDependOnKeyExcel(valueDependsOn configloader.ValueDependsOn, valueDependsOnKey string, fileHeader []string) string {
	if valueDependsOn != configloader.ValueDependsOnExcel {
		return valueDependsOnKey
	}
	if columnIndex, err := excelize.ColumnNameToNumber(valueDependsOnKey); err == nil && len(fileHeader) > columnIndex-1 {
		return fileHeader[columnIndex-1]
	}
	return valueDependsOnKey
}

func validateArrayItemMap(ctx context.Context, rowID int, rowData []string, requestFieldsMap map[string]*configloader.RequestFieldMD, fileParameters map[string]interface{}, fileHeader []string) (
	map[string]*configloader.RequestFieldMD, map[string]interface{}, []ErrorRow) {
	// 1. Init value
	var errorRows []ErrorRow
	childMap := make(map[string]interface{})

	// 2. Explore each item in array
	for fieldNameChild, reqFieldChild := range requestFieldsMap {
		// 2.1. Case field type is Array
		if len(reqFieldChild.ArrayItemMap) > 0 {
			// 2.1.1. Validate
			arrayItemMapUpdated, childMapInArr, errorRowsForArrayItem := validateArrayItemMap(ctx, rowID, rowData, reqFieldChild.ArrayItemMap, fileParameters, fileHeader)
			if len(errorRowsForArrayItem) > 0 {
				errorRows = append(errorRows, errorRowsForArrayItem...)
				continue
			}

			// 2.1.2. Update childMap and remove field if it was mapped data
			reqFieldChild.ArrayItemMap = arrayItemMapUpdated
			childMap[fieldNameChild] = []map[string]interface{}{childMapInArr}
			if len(childMapInArr) == 1 { // in case int[], string[], ... -> remove key empty, then convert map to array
				if val, ok := childMapInArr[""]; ok {
					childMap[fieldNameChild] = []interface{}{val}
				}
			} else if len(childMapInArr) == 0 { // childMapInArr is empty -> remove field
				delete(childMap, fieldNameChild)
			}
			if len(arrayItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove that field in arrayItemMap
				delete(requestFieldsMap, fieldNameChild)
			}

			// 2.1.3. Continue
			continue
		}

		// 2.2. Case field type is Object
		if len(reqFieldChild.ItemsMap) > 0 {
			// 2.2.1. Validate
			objectItemMapUpdated, childMapInObj, errorRowsForArrayItem := validateArrayItemMap(ctx, rowID, rowData, reqFieldChild.ItemsMap, fileParameters, fileHeader)
			if len(errorRowsForArrayItem) > 0 {
				errorRows = append(errorRows, errorRowsForArrayItem...)
				continue
			}

			// 2.2.2. Update childMap and remove field if it was mapped data
			reqFieldChild.ItemsMap = objectItemMapUpdated
			childMap[fieldNameChild] = childMapInObj
			if len(objectItemMapUpdated) == 0 { // if no remaining item that hasn't mapped -> remove that field in ItemsMap
				delete(requestFieldsMap, fieldNameChild)
			}

			// 2.2.3. Continue
			continue
		}

		// 2.3. In Normal case, field maybe int, string, ... -> need to get value (string) from config
		valueChildStr, isByPassField, errorRowsAfterGet := getValueStrByRequestFieldMD(ctx, rowID, rowData, reqFieldChild, fileParameters, fileHeader)
		if len(errorRowsAfterGet) > 0 {
			errorRows = append(errorRows, errorRowsAfterGet...)
			continue
		}
		if isByPassField { // no need to convert
			continue
		}

		// 2.4. Get real value from string value
		valueDependsOnKey := getValueDependOnKeyExcel(reqFieldChild.ValueDependsOn, reqFieldChild.ValueDependsOnKey, fileHeader)
		realValueChild, err := basejobmanager.ConvertToRealValue(ctx, reqFieldChild.Type, valueChildStr, valueDependsOnKey)
		if err != nil {
			errorRows = append(errorRows, ErrorRow{rowID, err.Error()})
		} else {
			if realValueChild != nil {
				childMap[reqFieldChild.Field] = realValueChild
			}
			// config will be converted to Json string, then save to DB -> delete to reduce size of json string
			delete(requestFieldsMap, fieldNameChild)
		}
	}

	// 3. Return
	return requestFieldsMap, childMap, errorRows
}

func getValueStrByRequestFieldMD(ctx context.Context, rowID int, rowData []string, reqField *configloader.RequestFieldMD, fileParameters map[string]interface{}, fileHeader []string) (string, bool, []ErrorRow) {
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
		cellValue, errorRowsExel := validateAndGetValueForRequestFieldExcel(ctx, rowID, rowData, reqField, fileHeader)
		if len(errorRowsExel) != 0 {
			errorRows = append(errorRows, errorRowsExel...)
		} else {
			valueStr = cellValue
		}
	case configloader.ValueDependsOnParam:
		paramValue, errorRowsParams := validateAndGetValueForFieldParam(ctx, rowID, reqField, fileParameters)
		if len(errorRowsParams) != 0 {
			errorRows = append(errorRows, errorRowsParams...)
		} else {
			valueStr = paramValue
		}
	case configloader.ValueDependsOnNone:
		valueStr = reqField.Value
	case configloader.ValueDependsOnTask:
		isByPassField = true // If value depends on Previous Task -> not get value
		valueDependsOnKeyMatched, errorRow := validateAndMatchJsonPath(ctx, rowID, rowData, reqField.ValueDependsOnKey, fileHeader)
		if errorRow != nil {
			errorRows = append(errorRows, *errorRow)
		} else {
			reqField.ValueDependsOnKey = valueDependsOnKeyMatched
		}
		valueStr = reqField.Value
	case configloader.ValueDependsOnFunc:
		isByPassField = true // If value depends on Function -> not get value, we will execute func later
		reqField.ValueDependsOnFunc.ParamsMapped = mapValueForCustomFunctionParams(reqField.ValueDependsOnFunc.ParamsRaw, rowData, fileParameters)
	case configloader.ValueDependsOnDb:
		isByPassField = true // If value depends on DB -> not get value, we will get after inserted task to DB
		valueStr = reqField.Value
	default:
		errMsg := fmt.Sprintf("cannot convert ValueDependsOn=%s", reqField.ValueDependsOn)
		errorRows = append(errorRows, ErrorRow{rowID, errMsg})

	}
	return valueStr, isByPassField, errorRows
}

func validateAndGetValueForFieldParam(ctx context.Context, rowID int, reqField *configloader.RequestFieldMD, fileParameters map[string]interface{}) (string, []ErrorRow) {
	var errorRows []ErrorRow
	paramKey := reqField.ValueDependsOnKey

	// Validate Require
	paramValue, existed := fileParameters[paramKey]
	if !existed && reqField.Required {
		reason := i18n.GetMessageCtx(ctx, "errConfigMissingParam", "name", paramKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
	}

	// convert param to string
	paramValueStr := constant.EmptyString
	if reqField.Type == configloader.TypeJson {
		jsonStr, _ := json.Marshal(paramValue)
		paramValueStr = string(jsonStr)
	} else {
		paramValueStr = fmt.Sprintf("%v", paramValue)
	}

	// set default value if empty
	if paramValueStr == constant.EmptyString {
		paramValueStr = reqField.DefaultValuePattern
	}

	// check required
	if reqField.Required && len(paramValueStr) == 0 {
		reason := i18n.GetMessageCtx(ctx, "errRowMissingDataColumn", "name", paramKey)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
		return "", errorRows
	}

	reqField.Value = paramValueStr
	return paramValueStr, errorRows
}

func validateAndGetValueForRequestFieldExcel(ctx context.Context, rowID int, rowData []string, reqField *configloader.RequestFieldMD, fileHeader []string) (string, []ErrorRow) {
	var errorRows []ErrorRow
	columnKey := reqField.ValueDependsOnKey
	columnIndex, err := excelize.ColumnNameToNumber(columnKey)
	if err != nil {
		logger.Errorf("validateResponseCode ... error %+v", err)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: err.Error()})
		return "", errorRows
	}
	// ColumnNameToNumber return value in range {1...}, but we expect columnIndex belongs to {0...}
	columnIndex--

	cellValue := constant.EmptyString
	if columnIndex < len(rowData) {
		cellValue = strings.TrimSpace(rowData[columnIndex])
	}
	if cellValue == constant.EmptyString {
		cellValue = reqField.DefaultValuePattern
	}

	// Validate Require
	if reqField.Required && cellValue == constant.EmptyString {
		columnName := columnKey
		if columnIndex < len(fileHeader) {
			columnName = fileHeader[columnIndex]
		}
		reason := i18n.GetMessageCtx(ctx, "errRowMissingDataColumn", "name", columnName)
		errorRows = append(errorRows, ErrorRow{RowId: rowID, Reason: reason})
		return "", errorRows
	}

	reqField.Value = cellValue

	return cellValue, errorRows
}

// validateAndMatchJsonPath ... support validate json path and update json path if it contains variable
// for example:
//   - Json path: data.transactions.#(name=="{{ $A }}").id
//   - Excel data has $A = quy
//     -> output: data.transactions.#(name=="quy").id
func validateAndMatchJsonPath(ctx context.Context, rowID int, rowData []string, jsonPath string, fileHeader []string) (string, *ErrorRow) {
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
		if !excel.IsColumnIndex(valuePattern) {
			errConfigMapping := i18n.GetMessageCtx(ctx, "errConfigMapping")
			errorRow := ErrorRow{RowId: rowID, Reason: errConfigMapping}
			logger.Errorf("validateResponseCode ... error %s -> %s", errConfigMapping, valuePatternWithDoubleBrace)
			return "", &errorRow
		}
		columnKey := valuePattern[1:] // if `$A` -> columnIndex = `A`
		columnIndex, err := excelize.ColumnNameToNumber(columnKey)
		if err != nil {
			errorRow := ErrorRow{RowId: rowID, Reason: err.Error()}
			logger.Errorf("validateResponseCode ... error %+v", err)
			return "", &errorRow
		}

		// 3.2. Validate value
		if columnIndex > len(rowData) || // column request out of range
			len(strings.TrimSpace(rowData[columnIndex-1])) == 0 { // column is required by value is empty
			columnName := columnKey
			if columnIndex < len(fileHeader) {
				columnName = fileHeader[columnIndex-1]
			}
			errRowMissingDataColumn := i18n.GetMessageCtx(ctx, "errRowMissingDataColumn", "name", columnName)
			reason := fmt.Sprintf("%s %s", errRowMissingDataColumn, columnName)
			errorRow := ErrorRow{RowId: rowID, Reason: reason}
			logger.Errorf("validateResponseCode ... error %+v", reason)
			return "", &errorRow
		}

		// 3.3. Replace value
		cellValue := strings.TrimSpace(rowData[columnIndex-1])
		jsonPath = strings.ReplaceAll(jsonPath, valuePatternWithDoubleBrace, cellValue)
	}

	// 4. return
	return jsonPath, nil
}

func mapValueForCustomFunctionParams(paramsRaw []string, rowData []string, fileParameters map[string]interface{}) []string {
	// 1. Init value for ParamsMapped
	paramsMapped := make([]string, len(paramsRaw))

	// 2. Explore all raw params and map data if we can
	for id, paramFuncPattern := range paramsRaw {
		paramMapped := paramFuncPattern

		// case get value column excel: $A, $B, ...
		if excel.IsColumnIndex(paramFuncPattern) {
			columnKey := paramFuncPattern[1:]
			paramMapped = excel.GetValueFromColumnKey(columnKey, rowData)
		} else
		// case get value from FileParameters
		if strings.Contains(paramFuncPattern, configloader.PrefixMappingRequestParameter) {
			paramKey := strings.TrimPrefix(paramFuncPattern, configloader.PrefixMappingRequestParameter+".")
			if paramValue, existed := fileParameters[paramKey]; existed {
				paramMapped = fmt.Sprintf("%+v", paramValue)
			}
		}

		// override mapped value
		paramsMapped[id] = paramMapped
	}

	return paramsMapped
}
