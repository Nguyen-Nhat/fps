package filewriter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	defaultSheetName          = "Data"
	defaultSheetNameFromExcel = "Sheet1"
)

func loadDataFromURL(fileURL string) ([]byte, error) {
	httpClient := &http.Client{}
	r, err := httpClient.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	return io.ReadAll(r.Body)
}

func writeDataToCsv(data [][]string) (*bytes.Buffer, error) {
	// 1. Create CSV file
	filePath := fmt.Sprintf("%d.csv", time.Now().UnixNano())
	csvFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	// 2. Write data to CSV file
	writer := csv.NewWriter(csvFile)
	defer func() {
		writer.Flush()
		err = csvFile.Close()
		if err != nil {
			logger.ErrorT("Failed to close file %v", filePath)
		}
		err = os.Remove(filePath)
		if err != nil {
			logger.ErrorT("Failed to remove file %v", filePath)
		}
	}()
	err = writer.WriteAll(data)
	if err != nil {
		return nil, err
	}

	// 3. Return bytes data
	byteData, err := os.ReadFile(filePath)
	return bytes.NewBuffer(byteData), err
}

func writeDataToXlsx(data [][]string, sheetName string) (*bytes.Buffer, error) {
	exFile := excelize.NewFile()

	// 1. Create a new sheet.
	if sheetName == constant.EmptyString {
		sheetName = defaultSheetName
	}

	// 2. Rename default sheet name
	err := exFile.SetSheetName(defaultSheetNameFromExcel, sheetName)
	if err != nil {
		return nil, err
	}

	// 3. Set value of a cell.
	for rowIndex, rowData := range data {
		cell, err := excelize.CoordinatesToCellName(1, rowIndex+1)
		if err != nil {
			return nil, err
		}
		err = exFile.SetSheetRow(sheetName, cell, &rowData)
		if err != nil {
			return nil, err
		}
	}

	return exFile.WriteToBuffer()
}

func getOutputFileContentType(outputFileType string) string {
	switch outputFileType {
	case constant.ExtFileCSV:
		return utils.CsvContentType
	case constant.ExtFileXLSX:
		return utils.XlsxContentType
	}
	return constant.EmptyString
}
