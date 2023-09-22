package executetask

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
)

type TaskGroupByRowType struct {
	fileID int
	rowID  int32
	tasks  []*fileprocessingrow.ProcessingFileRow
}

type BoundedParallelismParams struct {
	ctx                    context.Context
	taskGroupByRowChannels <-chan TaskGroupByRowType
	jobExecuteTask         *jobExecuteTask
}
