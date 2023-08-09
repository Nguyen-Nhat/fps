package fpRowGroup

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusSuccess    = 4
)

type ProcessingFileRowGroup struct {
	ent.ProcessingFileRowGroup
}

// Method --------------------------------------------------------------------------------------------------------------

func Name() string {
	return "ProcessingFileRowGroup"
}

func (rg *ProcessingFileRowGroup) IsInitStatus() bool {
	return rg.Status == StatusInit
}

func (rg *ProcessingFileRowGroup) IsSuccessStatus() bool {
	return rg.Status == StatusSuccess
}

func (rg *ProcessingFileRowGroup) IsFailedStatus() bool {
	return rg.Status == StatusFailed
}
