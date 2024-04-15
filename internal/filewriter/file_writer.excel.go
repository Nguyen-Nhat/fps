package filewriter

import (
	"bytes"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

type excelFileWriter struct {
	exFile         *excelize.File
	sheetName      string
	dataIndexStart int

	outputFileContentType string
}

func NewExcelFileWriter(fileURL, sheetName string, dataIndexStart int) (FileWriter, error) {
	// 1. Load file
	dataFile, err := loadDataFromURL(fileURL)
	if err != nil {
		return nil, err
	}

	// 2. Load excel
	exFile, err := excelize.OpenReader(bytes.NewReader(dataFile))
	if err != nil {
		logger.ErrorT("Failed to load file by url %v", fileURL)
		return nil, err
	}
	// -> Get first sheet if sheetName no data
	if sheetName == constant.EmptyString {
		sheetName = exFile.GetSheetName(0)
	}

	// 3. Return instant
	return &excelFileWriter{
		exFile:                exFile,
		sheetName:             sheetName,
		dataIndexStart:        dataIndexStart,
		outputFileContentType: utils.XlsxContentType,
	}, nil
}

// UpdateDataInColumnOfFile ... write {columnData} into column {columnName}
func (e *excelFileWriter) UpdateDataInColumnOfFile(columnName string, columnData map[int]string) error {
	// 1. Get column index
	var columnIndex string
	var err error
	if excel.IsColumnIndex(columnName) {
		columnIndex = columnName[1:]
	} else {
		columnIndex, err = excel.GetColumnIndexInFile(e.exFile, e.sheetName, columnName)
		if err != nil {
			columnIndex = excel.FirstColumnKey
			logger.Warnf("----> Force column %v is in `%v` column", columnName, columnIndex)
		}
	}

	// 2. Update column data
	for rowID, data := range columnData {
		axis := fmt.Sprintf("%v%v", columnIndex, e.dataIndexStart+rowID)
		err = e.exFile.SetCellValue(e.sheetName, axis, data)
		if err != nil {
			logger.Errorf("error when set value for cell %+v in sheet %+v, value = %+v", axis, e.sheetName, data)
		}
	}

	// 3. Return nil when no error
	return nil
}

func (e *excelFileWriter) OutputFileContentType() string {
	return e.outputFileContentType
}

func (e *excelFileWriter) GetFileBytes() (*bytes.Buffer, error) {
	return e.exFile.WriteToBuffer()
}

// private method ------------------------------------------------------------------------------------------------------
