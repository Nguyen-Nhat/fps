package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusFinished   = 4
)

type ProcessingFile struct {
	ent.ProcessingFile
}

type ErrorDisplay string

// ---------------------------------------------------------------------------------------------------------------------

func Name() string {
	return "ProcessingFile"
}

func (pf *ProcessingFile) IsInitStatus() bool {
	return pf.Status == StatusInit
}

func (pf *ProcessingFile) IsProcessingStatus() bool {
	return pf.Status == StatusProcessing
}
