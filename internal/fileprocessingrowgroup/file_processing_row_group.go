package fpRowGroup

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// Status ENUM ...
const (
	StatusInit             = 1
	StatusFailed           = 3
	StatusSuccess          = 4
	StatusCalledApiFail    = 5
	StatusCalledApiSuccess = 6
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

func (rg *ProcessingFileRowGroup) IsCalledAPI() bool {
	return rg.Status == StatusCalledApiFail || rg.Status == StatusCalledApiSuccess
}

func (rg *ProcessingFileRowGroup) IsCalledApiSuccess() bool {
	return rg.Status == StatusCalledApiSuccess
}
