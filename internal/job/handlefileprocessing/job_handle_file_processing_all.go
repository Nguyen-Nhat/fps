package handlefileprocessing

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

type JobHandleProcessingFileAll struct {
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	fileService fileservice.IService
}

func StartJobHandleProcessingFileAll(
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
) JobHandleProcessingFileAll {
	return JobHandleProcessingFileAll{
		fpService:   fpService,
		fprService:  fprService,
		fileService: fileService,
	}
}

func (j *JobHandleProcessingFileAll) StartJobForTesting() {
	jobName := "Job Handle ProcessingFile"
	logger.Infof("========== Running %v: Start  ...", jobName)
	j.start()
	logger.Infof("========== Running %v: Finish ...", jobName)
}

/*
I. Get all ProcessingFile in DB
II. If empty -> finish job
III. Execute each Processing File
*/

func (j *JobHandleProcessingFileAll) start() {
	ctx := context.Background()

	// I. Get all ProcessingFile in DB
	processingFiles, err := j.fpService.GetListFileAwardPointByStatuses(ctx, []int16{fileprocessing.StatusInit, fileprocessing.StatusProcessing})
	if err != nil {
		logger.Errorf("===== Cannot get list Processing File, got: %v", err)
		return
	}

	// II. Check empty
	if len(processingFiles) == 0 {
		logger.InfoT("No Init or Processing file for executing!")
		return
	}

	// III. Execute each Processing File
	for _, file := range processingFiles {
		jobExecutor := newJobHandleProcessingFile(j.fpService, j.fprService, j.fileService)
		jobExecutor.Execute(ctx, file)
	}
}
