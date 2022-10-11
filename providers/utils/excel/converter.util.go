package excel

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"reflect"
	"regexp"
	"strconv"
)

func ConvertToStruct[META any, OUT any, C dto.Converter[META, OUT]](
	dataIndexStart int, metadata C, data [][]string,
) (*dto.Sheet[OUT], error) {
	// Check empty
	totalRows := len(data)
	fmt.Printf("Total rows = %v\n", totalRows)
	if totalRows < dataIndexStart {
		return nil, fmt.Errorf("sheet empty")
	}

	// Mapping Header with Column Index in file
	headers := metadata.GetHeaders()
	headerMap, err := mappingHeaderWithColumnIndexInFile(headers, data[0])
	if err != nil {
		return nil, err
	}

	// Explore each row data
	var dataRows []OUT
	var errorRows []dto.ErrorRow
	for rowId := dataIndexStart - 1; rowId < len(data); rowId++ { // each row
		mt := metadata         // clone object
		realRowId := rowId + 1 // real row id of data
		errorRow, isOk := readDataReflect(mt.GetMetadata(), headerMap, data[rowId])
		if isOk {
			dataRows = append(dataRows, mt.ToOutput(realRowId))
		} else {
			errorRow.RowId = realRowId
			errorRows = append(errorRows, *errorRow)
		}
	}

	return &dto.Sheet[OUT]{
		DataIndexStart: dataIndexStart,
		Data:           dataRows,
		ErrorRows:      errorRows,
	}, nil
}

func mappingHeaderWithColumnIndexInFile(headers []string, headersInFile []string) (map[string]int, error) {
	numberColumnAppear := 0
	var headerMap = make(map[string]int)
	for _, header := range headers {
		for j, columnName := range headersInFile {
			if columnName == header {
				headerMap[columnName] = j
				numberColumnAppear++
			}
		}
	}
	if numberColumnAppear != len(headers) {
		return nil, fmt.Errorf("not enough column in file")
	}
	return headerMap, nil
}

func readDataReflect[META any](metadata *META, headerMap map[string]int, rowData []string) (*dto.ErrorRow, bool) {
	el := reflect.ValueOf(metadata).Elem()

	fmt.Printf("numfield %v", el.NumField())

	for i := 0; i < el.NumField(); i++ {
		fieldName := el.Type().Field(i).Name
		varValue := el.Field(i).Interface()

		if reflect.TypeOf(varValue) == reflect.TypeOf(dto.CellData[string]{}) { // For String type
			cellData := varValue.(dto.CellData[string])
			errorMsg, hasError := validateAndSetRawValue(headerMap, &cellData, rowData)
			if hasError {
				return &dto.ErrorRow{Reason: fmt.Sprintf("%v %v", fieldName, errorMsg)}, false
			}

			cellData.Value = cellData.ValueRaw
			el.Field(i).Set(reflect.ValueOf(cellData))
		} else if reflect.TypeOf(varValue) == reflect.TypeOf(dto.CellData[int]{}) { // For Int type
			cellData := varValue.(dto.CellData[int])
			if len(cellData.ColumnName) == 0 { // in case no column name -> ignore
				continue
			}
			errorMsg, hasError := validateAndSetRawValue(headerMap, &cellData, rowData)
			if hasError {
				return &dto.ErrorRow{Reason: fmt.Sprintf("%v %v", fieldName, errorMsg)}, false
			}

			value, err := strconv.Atoi(cellData.ValueRaw)
			if err != nil {
				msg := fmt.Sprintf("%v is missing or wrong format", fieldName)
				return &dto.ErrorRow{Reason: msg}, false
			}
			cellData.Value = value
			el.Field(i).Set(reflect.ValueOf(cellData))
		} else {
			return &dto.ErrorRow{Reason: "Failed execute " + fieldName}, false
		}
	}

	return nil, true
}

func validateAndSetRawValue[V any](headerMap map[string]int, cellData *dto.CellData[V], rowData []string) (string, bool) {
	columnId, ok := headerMap[cellData.ColumnName]

	rawValue := ""
	if ok && columnId < len(rowData) {
		rawValue = rowData[columnId]
	}
	cellData.ValueRaw = rawValue

	constrains := cellData.Constrains

	if constrains.IsRequired && len(rawValue) == 0 {
		return "required", true
	}

	if constrains.MinLength != nil && len(rawValue) < *constrains.MinLength {
		return fmt.Sprintf("length must >= %v", *constrains.MinLength), true
	}

	if constrains.MaxLength != nil && *constrains.MaxLength < len(rawValue) {
		return fmt.Sprintf("length must <= %v", *constrains.MaxLength), true
	}

	if constrains.Regexp != nil && len(rawValue) > 0 {
		matches := regexp.MustCompile(*constrains.Regexp).FindAllString(rawValue, -1)
		if len(matches) <= 0 {
			return "wrong format", true
		}
	}

	return "", false
}
