package flatten

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/workers"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

// jobFlattenManager ...
type jobFlattenManager struct {
	cfg     config.SchedulerConfig
	cronJob *cron.Cron
	// services
	fpService         fileprocessing.Service
	fprService        fileprocessingrow.Service
	fpRowGroupService fpRowGroup.Service
	fileService       fileservice.IService
	// services config
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

var jobFlattenMgr *jobFlattenManager
var once sync.Once

// NewJobFlattenManager ...
func NewJobFlattenManager(
	cfg config.SchedulerConfig,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fpRowGroupService fpRowGroup.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
	cfgTaskService configtask.Service,
) basejobmanager.CronJobManager {
	if jobFlattenMgr == nil {
		once.Do(func() {
			jobFlattenMgr = &jobFlattenManager{
				cfg:               cfg,
				fpService:         fpService,
				fprService:        fprService,
				fpRowGroupService: fpRowGroupService,
				fileService:       fileService,
				cfgMappingService: cfgMappingService,
				cfgTaskService:    cfgTaskService,
			}
		})
	}

	jobFlattenMgr.cronJob = basejobmanager.InitCron(jobFlattenMgr)

	return jobFlattenMgr
}

func (mgr *jobFlattenManager) Start() {
	mgr.cronJob.Start()
}

func (mgr *jobFlattenManager) GetJobName() string {
	return "Job Flatten ProcessingFile"
}

func (mgr *jobFlattenManager) GetSchedulerConfig() config.SchedulerConfig {
	return mgr.cfg
}

// Execute ...
// Logic:
//  1. Fetch all file that have status = INIT
//  2. Flattening each file
func (mgr *jobFlattenManager) Execute() {
	// 1. Fetch all file that have status = INIT
	ctx := context.Background()
	fpList, err := mgr.fpService.GetListFileByStatuses(ctx, []int16{fileprocessing.StatusInit})
	if err != nil {
		logger.Errorf("===== Cannot get list Processing File, got: %v", err)
		return
	}

	// 2. Check empty
	if len(fpList) == 0 {
		logger.InfoT("No INIT file for executing!")
		return
	}

	// 3. Flattening each file
	jobFlatten := newJobFlatten(
		mgr.fpService, mgr.fprService, mgr.fpRowGroupService,
		mgr.fileService,
		mgr.cfgMappingService, mgr.cfgTaskService)
	workerPool := workers.NewWorkerPool(mgr.cfg.NumDigesters)
	workerPool.Run()
	for _, fp := range fpList {
		tmpFp := fp
		workerPool.AddTask(func() {
			jobFlatten.Flatten(ctx, *tmpFp)
		})
	}
	workerPool.Close()
}
