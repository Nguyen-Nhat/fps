package fpRowGroup

import (
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

func toProcessingFileRowGroupArr(request []CreateRowGroupJob) []ProcessingFileRowGroup {
	var res []ProcessingFileRowGroup
	for _, req := range request {
		res = append(res, toProcessingFileRowGroup(req))
	}
	return res
}

func toProcessingFileRowGroup(request CreateRowGroupJob) ProcessingFileRowGroup {
	return ProcessingFileRowGroup{
		ProcessingFileRowGroup: ent.ProcessingFileRowGroup{
			FileID:       int64(request.FileID),
			TaskIndex:    int32(request.TaskIndex),
			GroupByValue: request.GroupByValue,
			TotalRows:    int32(request.TotalRows),
			RowIndexList: request.RowIndexList,
			Status:       StatusInit,
			ExecutedTime: -1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
}
