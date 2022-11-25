package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
)

func mapStatus(statusInDB int16) string {
	switch statusInDB {
	case fileprocessing.StatusInit:
		return FpStatusInit
	case fileprocessing.StatusProcessing:
		return FpStatusProcessing
	case fileprocessing.StatusFailed:
		return FpStatusFailed
	case fileprocessing.StatusFinished:
		return FpStatusFinished
	default:
		return ""
	}
}
