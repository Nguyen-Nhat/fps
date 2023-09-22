package updatestatus

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
)

type BoundedParallelismParams struct {
	ctx                   context.Context
	jobFlatten            *jobUpdateStatus
	fileProcessingChannel <-chan fileprocessing.ProcessingFile
	numDigesters          int
}
