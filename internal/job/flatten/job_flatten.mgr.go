package flatten

import (
	"context"
	"github.com/robfig/cron/v3"
	"sync"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

// JobFlattenManager ...
type JobFlattenManager struct {
	cfg config.SchedulerConfig
	// services
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	fileService fileservice.IService
	// services config
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

var jobFlattenMgr *JobFlattenManager
var once sync.Once

// NewJobFlattenManager ...
func NewJobFlattenManager(
	cfg config.SchedulerConfig,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
	cfgTaskService configtask.Service,
) *JobFlattenManager {
	if jobFlattenMgr == nil {
		once.Do(func() {
			jobFlattenMgr = &JobFlattenManager{
				cfg:               cfg,
				fpService:         fpService,
				fprService:        fprService,
				fileService:       fileService,
				cfgMappingService: cfgMappingService,
				cfgTaskService:    cfgTaskService,
			}
			jobFlattenMgr.initCron()
		})
	}
	return jobFlattenMgr
}

// Execute ...
// Logic:
//  1. Fetch all file that have status = INIT
//  2. Flattening each file
func (mgr *JobFlattenManager) Execute() {
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
	jobFlatten := newJobFlatten(mgr.cfg, mgr.fpService, mgr.fprService, mgr.fileService, mgr.cfgMappingService, mgr.cfgTaskService)
	for _, fp := range fpList { // todo: can use multi thread for improving performance
		jobFlatten.Flatten(ctx, *fp)
	}
}

func (mgr *JobFlattenManager) initCron() {
	cronFlatten := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	jobName := "Job Flatten ProcessingFile"

	id, err := cronFlatten.AddFunc(mgr.cfg.Schedule, func() {
		logger.Infof("\n")
		logger.Infof("========== Running %v: Start  ...", jobName)
		mgr.Execute()
		logger.Infof("========== Running %v: Finish ...\n", jobName)
	})
	if err != nil {
		logger.Errorf("Init Job failed: %v", err)
	}

	logger.Infof("Init Job success: ID = %v", id)

	cronFlatten.Start()
}
