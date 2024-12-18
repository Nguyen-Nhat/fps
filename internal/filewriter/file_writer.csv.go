package filewriter

import (
	"bytes"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	csvUtil "git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/csv"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

type csvFileWriter struct {
	fileData       [][]string
	dataIndexStart int
	sheetName      string

	outputFileType string
}

func NewCsvFileWriter(fileURL, sheetName string, dataIndexStart int, outputFileType string) (FileWriter, error) {
	// 1. Load CSV data
	allRowsData, err := csvUtil.LoadCSVByURL(fileURL)
	if err != nil {
		logger.ErrorT("Failed to get all rows in file")
		return nil, err
	}

	// 2. Return instant
	return &csvFileWriter{
		fileData:       allRowsData,
		dataIndexStart: dataIndexStart,
		sheetName:      sheetName,
		outputFileType: outputFileType,
	}, nil
}

// UpdateDataInColumnOfFile ... write {columnData} into column {columnName}
func (c *csvFileWriter) UpdateDataInColumnOfFile(columnName string, columnData map[int]string) error {
	// 1. Detect column_index base on column_name
	allRowsData := c.fileData
	var columnIndex int
	var err error
	if excel.IsColumnIndex(columnName) {
		columnIndex, err = excelize.ColumnNameToNumber(columnName[1:])
		if err != nil {
			columnIndex = 1 // default is first column (index from 1 -> n)
			logger.Warnf("----> Force column %v is in `%v` column", columnName, columnIndex)
		}
	} else {
		if index, existed := utils.IndexOf(allRowsData[0], columnName); existed {
			columnIndex = index + 1
		} else {
			columnIndex = 1 // if we can not find column name, return the first column index
		}
	}

	// 2. Set data into CSV data
	resultData := make([][]string, len(allRowsData))
	for rowID, rowData := range allRowsData {
		// init array, extend size if columnIndex is big
		resultData[rowID] = make([]string, len(rowData))
		if columnIndex > len(rowData) {
			resultData[rowID] = make([]string, columnIndex)
		}

		// copy data of row
		copy(resultData[rowID], rowData)

		// ignore if row is not Data
		if rowID < c.dataIndexStart-1 {
			continue
		}

		// Update data
		if data, existed := columnData[rowID-c.dataIndexStart+1]; existed {
			resultData[rowID][columnIndex-1] = data
		}
	}

	// 3. Replace value for file data
	c.fileData = resultData
	return nil
}

func (c *csvFileWriter) OutputFileContentType() string {
	return getOutputFileContentType(c.outputFileType)
}

func (c *csvFileWriter) GetFileBytes() (*bytes.Buffer, error) {
	switch c.outputFileType {
	case constant.ExtFileXLSX:
		return writeDataToXlsx(c.fileData, c.sheetName)
	default:
		return writeDataToCsv(c.fileData)
	}
}

// private method ------------------------------------------------------------------------------------------------------
