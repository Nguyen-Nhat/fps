package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

type ProcessingFile struct {
	ent.ProcessingFile
}

func Name() string {
	return "ProcessingFile"
}

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusFinished   = 4
)
