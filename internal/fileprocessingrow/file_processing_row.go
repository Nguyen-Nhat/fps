package fileprocessingrow

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"strconv"
	"strings"
)

// Status ENUM ...
const (
	StatusInit    = 1
	StatusFailed  = 3
	StatusSuccess = 4
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
func (s *CustomStatisticModel) GetErrorDisplay() string {
	errorDisplay := ""

	taskErrorDisplays := strings.Split(s.ErrorDisplays, ",")
	taskStatuses := strings.Split(s.Statuses, ",")

	for i := 0; i < len(taskStatuses); i++ {
		status := taskStatuses[i]
		statusInt, _ := strconv.Atoi(status)
		if statusInt == StatusFailed {
			if len(errorDisplay) == 0 {
				errorDisplay = taskErrorDisplays[i]
			} else {
				errorDisplay = errorDisplay + ", " + taskErrorDisplays[i]
			}
		}
	}

	return errorDisplay
}
