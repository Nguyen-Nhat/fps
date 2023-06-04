package configloader

import (
	"errors"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
)

// excelConfigLoader ...
type excelConfigLoader struct {
}

// ---------------------------------------------------------------------------------------------------------------------

func (l *excelConfigLoader) Load(file fileprocessing.ProcessingFile) (ConfigMappingMD, error) {
	return ConfigMappingMD{}, errors.New("not support")
}
