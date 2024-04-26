package filewriter

import (
	"bytes"
	"fmt"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	configmappingEnt "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
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
		case configmappingEnt.OutputFileTypeXLS:
			outputFileType = constant.ExtFileXLS
		default:
			return nil, fmt.Errorf("outputFileType %v is not supported", outputFileTypeCfg)
		}
	}

	// 2. get File Writer implementation
	switch outputFileType {
	case constant.ExtFileCSV:
		return NewCsvFileWriter(fileURL, dataIndexStart)
	case constant.ExtFileXLSX:
		return NewExcelFileWriter(fileURL, sheetName, dataIndexStart, utils.XlsxContentType)
	case constant.ExtFileXLS:
		return NewExcelFileWriter(fileURL, sheetName, dataIndexStart, utils.XlsContentType)
	default:
		return nil, fmt.Errorf("outputFileType %v is not supported", outputFileTypeCfg)
	}
}
