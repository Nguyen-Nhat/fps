package fileprocessingrow

import (
	"strconv"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// Status ENUM ...
const (
	StatusInit            = 1
	StatusFailed          = 3
	StatusSuccess         = 4
	StatusWaitForGrouping = 5
	StatusRejected        = 6
	StatusWaitForAsync    = 7
)

type ProcessingFileRow struct {
	ent.ProcessingFileRow
}

// for custom struct, that is result of custom queries, complex queries
// prefix is `Custom...`
type (

	// CustomStatisticModel ...
	CustomStatisticModel struct {
		RowIndex      int
		Statuses      string
		Count         int
		ErrorDisplays string
	}
)

// Method --------------------------------------------------------------------------------------------------------------

func Name() string {
	return "ProcessingFileRow"
}

func (pf *ProcessingFileRow) IsInitStatus() bool {
	return pf.Status == StatusInit
}

func (pf *ProcessingFileRow) IsSuccessStatus() bool {
	return pf.Status == StatusSuccess
}

func (pf *ProcessingFileRow) IsFailedStatus() bool {
	return pf.Status == StatusFailed
}

func (pf *ProcessingFileRow) IsWaitForGroupingStatus() bool {
	return pf.Status == StatusWaitForGrouping
}

func (s *CustomStatisticModel) IsSuccessAll() bool {
	statusSuccessStr := strconv.Itoa(StatusSuccess)

	tmp := statusSuccessStr
	for i := 1; i < s.Count; i++ {
		tmp += "," + statusSuccessStr
	}
	// =>>> output = 4,4,4...,4

	return tmp == s.Statuses
}

func (s *CustomStatisticModel) IsContainsFailed() bool {
	statusFailedStr := strconv.Itoa(StatusFailed)
	return strings.Contains(s.Statuses, statusFailedStr)
}

// IsProcessed ... TRUE when contains terminated status (Success/Failed)
func (s *CustomStatisticModel) IsProcessed() bool {
	statusFailedStr := strconv.Itoa(StatusFailed)
	statusSuccessStr := strconv.Itoa(StatusSuccess)
	return strings.Contains(s.Statuses, statusFailedStr) || strings.Contains(s.Statuses, statusSuccessStr)
}

// IsWaiting ... TRUE when contains waiting status
func (s *CustomStatisticModel) IsWaiting() bool {
	statusFailedStr := strconv.Itoa(StatusWaitForGrouping)
	return strings.Contains(s.Statuses, statusFailedStr)
}
