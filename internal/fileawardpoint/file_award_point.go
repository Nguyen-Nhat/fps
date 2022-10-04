package fileawardpoint

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// FileAwardPoint is model of table `file_award_point`
type FileAwardPoint struct {
	ent.FileAwardPoint
}

// Status ENUM ...
const (
	statusInit       = 0
	statusProcessing = 1
	statusSuccess    = 2
	statusFailed     = 3
)
