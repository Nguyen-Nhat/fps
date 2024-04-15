package filewriter

import (
	"bytes"
	"fmt"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	configmappingEnt "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
)

type FileWriter interface {
	UpdateDataInColumnOfFile(columnName string, columnData map[int]string) error
	OutputFileContentType() string
	GetFileBytes() (*bytes.Buffer, error)
}

func NewFileWriter(fileURL, sheetName string, dataIndexStart int,
	inputFileType string, outputFileTypeCfg configmappingEnt.OutputFileType) (FileWriter, error) {
	// 1. input-output pair is match?
	switch strings.ToUpper(inputFileType) {
	case constant.ExtFileCSV:
		if outputFileTypeCfg != constant.ExtFileCSV {
			return nil, fmt.Errorf("InputFileType %v and OutputFileType %v are not same", constant.ExtFileCSV, outputFileTypeCfg)
		}
	default:
		if outputFileTypeCfg != constant.ExtFileXLSX {
			return nil, fmt.Errorf("InputFileType %v and OutputFileType %v are not same", constant.ExtFileXLSX, outputFileTypeCfg)
		}
	}

	// 2. get File Writer implementation
	switch outputFileTypeCfg {
	case configmappingEnt.OutputFileTypeCSV:
		return NewCsvFileWriter(fileURL, dataIndexStart)
	case configmappingEnt.OutputFileTypeXLSX:
		return NewExcelFileWriter(fileURL, sheetName, dataIndexStart)
	default:
		return nil, fmt.Errorf("outputFileType %v is not supported", outputFileTypeCfg)
	}
}
