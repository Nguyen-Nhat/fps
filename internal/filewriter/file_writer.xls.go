package filewriter

import (
	"bytes"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/xls"
)

type xlsFileWriter struct {
	fileData       [][]string
	sheetName      string
	dataIndexStart int

	outputFileType string
}

func NewXlsFileWriter(fileURL, sheetName string, dataIndexStart int, outputFileType string) (FileWriter, error) {
	// 1. Load file
	fileData, err := xls.LoadXlsByUrl(fileURL, sheetName)
	if err != nil {
		return nil, err
	}

	// 2. Return instant
	return &xlsFileWriter{
		fileData:       fileData,
		sheetName:      sheetName,
		dataIndexStart: dataIndexStart,
		outputFileType: outputFileType,
	}, nil
}

// UpdateDataInColumnOfFile ... write {columnData} into column {columnName}
func (x *xlsFileWriter) UpdateDataInColumnOfFile(columnName string, columnData map[int]string) error {
	allRowsData := x.fileData
	// 1. Get column index
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

	// 2. Set data into file data
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
		if rowID < x.dataIndexStart-1 {
			continue
		}

		// Update data
		if data, existed := columnData[rowID-x.dataIndexStart+1]; existed {
			resultData[rowID][columnIndex-1] = data
		}
	}

	// 3. Replace value for file data
	x.fileData = resultData
	return nil
}

func (x *xlsFileWriter) OutputFileContentType() string {
	return getOutputFileContentType(x.outputFileType)
}

func (x *xlsFileWriter) GetFileBytes() (*bytes.Buffer, error) {
	switch x.outputFileType {
	case constant.ExtFileCSV:
		return writeDataToCsv(x.fileData)
	default:
		return writeDataToXlsx(x.fileData, x.sheetName)
	}
}

// private method ------------------------------------------------------------------------------------------------------
