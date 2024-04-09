package handlefileprocessing

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

type JobHandleProcessingFileAll struct {
	cfg               config.FileProcessingConfig
	fpService         fileprocessing.Service
	fprService        fileprocessingrow.Service
	fileService       fileservice.IService
	cfgMappingService configmapping.Service
}

var fileProcessingJob *JobHandleProcessingFileAll
var once sync.Once

func InitJobHandleProcessingFileAll(
	cfg config.FileProcessingConfig,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
) *JobHandleProcessingFileAll {
	if fileProcessingJob == nil {
		once.Do(func() {
			fileProcessingJob = &JobHandleProcessingFileAll{
				cfg:               cfg,
				fpService:         fpService,
				fprService:        fprService,
				fileService:       fileService,
				cfgMappingService: cfgMappingService,
			}
			fileProcessingJob.initCron()
		})
	}
	return fileProcessingJob
}

func (j *JobHandleProcessingFileAll) initCron() {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	jobName := "Job Handle ProcessingFile"

	id, err := c.AddFunc(j.cfg.Schedule, func() {
		logger.Infof("\n")
		logger.Infof("========== Running %v: Start  ...", jobName)
		j.start()
		logger.Infof("========== Running %v: Finish ...\n", jobName)

	})
	if err != nil {
		logger.Errorf("Init Job failed: %v", err)
	}

	logger.Infof("Init Job success: ID = %v", id)

	c.Start()
}

/*
I. Get all ProcessingFile in DB
II. If empty -> finish job
III. Execute each Processing File
*/

func (j *JobHandleProcessingFileAll) start() {
	ctx := context.Background()

	// I. Get all ProcessingFile in DB
	processingFiles, err := j.fpService.GetListFileByStatuses(ctx, []int16{fileprocessing.StatusInit, fileprocessing.StatusProcessing})
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
		jobExecutor := newJobHandleProcessingFile(j.fpService, j.fprService, j.fileService, j.cfgMappingService)
		jobExecutor.Execute(ctx, file)
	}
}
