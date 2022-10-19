package fileawardpoint

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// FileAwardPoint is model of table `file_award_point`
type FileAwardPoint struct {
	ent.FileAwardPoint
}

func Name() string {
	return "FileAwardPoint"
}

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusFinished   = 4
)
