package excel

import (
	"fmt"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
)

func ConvertToStructByMap(dataIndexStart int, metadata []dto.CellData[string], data [][]string) (*dto.Sheet[map[string]string], error) {
	// 1. Check empty
	totalRows := len(data)
	logger.Infof("Total rows = %v\n", totalRows)
	if totalRows < dataIndexStart || dataIndexStart < 2 {
		return nil, ErrEmptySheet
	}

	// 2. Mapping Header with Column Index in file
	headerMap, err := mappingHeader(metadata)
	if err != nil {
		return nil, err
	}

	// 3. Validate header -> if list column not contain header then return error
	headerRow := data[0]
	for columnName, columnIndex := range headerMap {
		if columnIndex >= len(headerRow) || len(headerRow[columnIndex]) <= 0 {
			logger.Errorf("column %v not contains header!!!", columnName)
			return nil, fmt.Errorf("column %v not contains header", columnName)
		}
	}

	// 4. Explore each row data
	var dataRows []map[string]string
	var errorRows []dto.ErrorRow
	for rowId := dataIndexStart - 1; rowId < len(data); rowId++ { // each row
		mt := metadata         // clone object
		realRowId := rowId + 1 // real row id of data
		errorRow, isOk := readDataReflectMapping(mt, headerMap, data[rowId])
		if isOk {
			output := ToOutput(realRowId, mt)
			dataRows = append(dataRows, output)
		} else {
			errorRow.RowId = realRowId
			errorRows = append(errorRows, *errorRow)
		}
	}

	// 5. return
	return &dto.Sheet[map[string]string]{
		DataIndexStart: dataIndexStart,
		Data:           dataRows,
		ErrorRows:      errorRows,
	}, nil
}

func ToOutput(rowId int, cellData []dto.CellData[string]) map[string]string {
	output := make(map[string]string)
	for _, cell := range cellData {
		output[cell.ColumnName] = cell.Value
	}
	return output
}

func readDataReflectMapping(metadata []dto.CellData[string], headerMap map[string]int, rowData []string) (*dto.ErrorRow, bool) {

	for i, cellData := range metadata {
		errorMsg, hasError := validateAndSetRawValue(headerMap, &cellData, rowData)
		if hasError {
			return &dto.ErrorRow{Reason: fmt.Sprintf("Column %v error %v", cellData.ColumnName, errorMsg), RowData: rowData}, false
		}
		cellData.Value = cellData.ValueRaw
		metadata[i] = cellData
	}

	return nil, true
}

// mappingHeader ... return ( map[column]columnInd ). For example: map[A]=1, map[B]=2, ...
func mappingHeader(columnsData []dto.CellData[string]) (map[string]int, error) {
	var headerMapping = make(map[string]int)
	for _, data := range columnsData {
		columnName := data.ColumnName

		if len(columnName) != 1 {
			return nil, fmt.Errorf("only support column has 1 character, such as {A, B, C, ...}, invalid value is %v", columnName)
		}

		columnIndex := int(strings.ToUpper(columnName)[0]) - int('A') // get first character then
		headerMapping[columnName] = columnIndex

	}

	return headerMapping, nil
}
