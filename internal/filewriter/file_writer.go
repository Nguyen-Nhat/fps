package filewriter

import (
	"bytes"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/errorz"
	configmappingEnt "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
)

type FileWriter interface {
	UpdateDataInColumnOfFile(columnName string, columnData map[int]string) error
	OutputFileContentType() string
	GetFileBytes() (*bytes.Buffer, error)
}

func NewFileWriter(fileURL, sheetName string, dataIndexStart int,
	inputFileType string, outputFileTypeCfg configmappingEnt.OutputFileType) (FileWriter, error) {
	// 1. if output file type is not set, then use input file type. Otherwise, use output file type
	outputFileType := strings.ToUpper(inputFileType)
	if len(outputFileTypeCfg) > 0 {
		switch outputFileTypeCfg {
		case configmappingEnt.OutputFileTypeCSV:
			outputFileType = constant.ExtFileCSV
		case configmappingEnt.OutputFileTypeXLSX:
			outputFileType = constant.ExtFileXLSX
		default:
			return nil, errorz.OutputFileTypeNotSupported(outputFileTypeCfg.String())
		}
	}

	// 1.1. prevent case export XLS
	if outputFileType == constant.ExtFileXLS {
		return nil, errorz.OutputFileTypeNotSupported(outputFileType)
	}

	// 2. get File Writer implementation
	switch strings.ToUpper(inputFileType) {
	case constant.ExtFileCSV:
		return NewCsvFileWriter(fileURL, sheetName, dataIndexStart, outputFileType)
	case constant.ExtFileXLSX:
		return NewExcelFileWriter(fileURL, sheetName, dataIndexStart, outputFileType)
	case constant.ExtFileXLS:
		return NewXlsFileWriter(fileURL, sheetName, dataIndexStart, outputFileType)
	default:
		return nil, errorz.InputFileTypeNotSupported(inputFileType)
	}
}
