package fileprocessingrow

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

func toProcessingFileRowArr(request []CreateProcessingFileRowJob) []ProcessingFileRow {
	var res []ProcessingFileRow
	for _, req := range request {
		res = append(res, toProcessingFileRow(req))
	}
	return res
}

func toProcessingFileRow(request CreateProcessingFileRowJob) ProcessingFileRow {
	return ProcessingFileRow{
		ProcessingFileRow: ent.ProcessingFileRow{
			FileID:      int64(request.FileId),
			RowIndex:    int32(request.RowIndex),
			RowDataRaw:  request.RowDataRaw,
			TaskIndex:   int32(request.TaskIndex),
			TaskMapping: request.TaskMapping,
			Status:      StatusInit,
		},
	}
}