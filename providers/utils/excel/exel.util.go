package excel

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const prefixMappingRequest = "$"
const FirstColumnKey = "A"

func LoadExcelByUrl(fileURL string, sheetName string) ([][]string, error) {
	var sheetNames []string
	if len(sheetName) > 0 {
		sheetNames = []string{sheetName}
	}

	sheetMap, err := LoadSheetsInExcelByUrl(fileURL, sheetNames)
	if err != nil {
		return nil, err
	}

	for k := range sheetMap {
		return sheetMap[k], nil
	}

	return nil, fmt.Errorf("sheet empty")
}

func LoadSheetsInExcelByUrl(fileURL string, sheetNameArr []string) (map[string][][]string, error) {
	// 1. Load data to bytes
	data, err := utils.GetDataFromURL(fileURL)
	if err != nil {
		return nil, err
	}

	// 2. Read file with Excelize
	excel, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		logger.Infof("Reader", err)
		return nil, err
	}

	sheetMap, err := validateAndGetDataInSheets(excel, sheetNameArr)
	if err != nil {
		return nil, err
	}

	// 6. Close file excel
	defer func() {
		if err = excel.Close(); err != nil {
			logger.Errorf("Cannot close excel file, got: %v", err)
		}
	}()

	// 7. Return
	return sheetMap, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func validateAndGetDataInSheets(excel *excelize.File, listSheetName []string) (map[string][][]string, error) {
	// 3. Check no sheet
	lst := excel.GetSheetList()
	if len(lst) == 0 {
		logger.Infof("Empty document")
		return nil, fmt.Errorf("file empty")
	}

	// 4. Get default sheetName if listSheetName is empty
	if len(listSheetName) == 0 { // Get First Sheet if no name
		sheetName := excel.GetSheetName(0)
		fmt.Printf("First sheet is %v\n", sheetName)
		listSheetName = []string{sheetName}
	}

	// 5. Get data in listSheetName
	sheetMap := make(map[string][][]string) // map of 2-dimensional arrays, with key is sheetName
	for _, sheetName := range listSheetName {
		sheetData, err := excel.GetRows(sheetName)
		if err != nil {
			logger.Errorf("Get data in sheet `%v` error %v", sheetName, err)
			return nil, err
		}
		sheetMap[sheetName] = sheetData
	}

	return sheetMap, nil
}

// IsColumnIndex ... return TRUE if columnName start with `$` another character is from A-Z
// Eg: $A, $D, ... -> TRUE
func IsColumnIndex(columnName string) bool {
	if len(columnName) < 2 || !strings.HasPrefix(columnName, prefixMappingRequest) {
		return false
	}
	for _, c := range columnName[1:] {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// GetColumnIndexInFile ... return the position of column by name
func GetColumnIndexInFile(exFile *excelize.File, sheetName, columnName string) (string, error) {
	// 1. Get data in Sheet
	sheetData, err := exFile.GetRows(sheetName)
	if err != nil || sheetData == nil {
		logger.ErrorT("Failed to load sheet %v in file", sheetName)
		return "", err
	}

	// 2. Get list Header then select one of them
	headers := sheetData[0]
	columnIndex := ""
	for index, header := range headers {
		if header == columnName {
			// no need to handle error because we will check length column index in below code
			columnIndex, _ = excelize.ColumnNumberToName(index)
			break
		}
	}

	// return correct value
	if columnIndex != "" {
		return columnIndex, nil
	}

	// return error
	errMsg := fmt.Sprintf("not found column %v in in sheet %v", columnName, sheetName)
	logger.ErrorT(errMsg)
	return "", fmt.Errorf("%s", errMsg)
}

// GetValueFromColumnKey ...
func GetValueFromColumnKey(columnKey string, data []string) string {
	if columnKey == "" {
		return ""
	}
	columnIndex, err := excelize.ColumnNameToNumber(strings.ToUpper(columnKey))
	if err != nil {
		return ""
	}
	if columnIndex <= len(data) { // column request out of range
		return strings.TrimSpace(data[columnIndex-1])
	}
	return ""
}
