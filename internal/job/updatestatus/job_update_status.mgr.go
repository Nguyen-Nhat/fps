package updatestatus

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/adapter/slack"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

// jobUpdateStatusManager ...
type jobUpdateStatusManager struct {
	cfg     config.SchedulerConfig
	cronJob *cron.Cron
	// services
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	fileService fileservice.IService
	// services config
	cfgMappingService configmapping.Service
	slackClient       slack.Client
}

var jobUpdateStatusMgr *jobUpdateStatusManager
var once sync.Once

// NewJobUpdateStatusManager ...
func NewJobUpdateStatusManager(
	cfg config.Config,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
) basejobmanager.CronJobManager {
	if jobUpdateStatusMgr == nil {
		once.Do(func() {
			slackClient := slack.NewSlackClient(cfg.SlackWebhook)

			jobUpdateStatusMgr = &jobUpdateStatusManager{
				cfg:               cfg.JobConfig.UpdateStatusConfig,
				fpService:         fpService,
				fprService:        fprService,
				fileService:       fileService,
				cfgMappingService: cfgMappingService,
				slackClient:       slackClient,
			}
		})
	}

	jobUpdateStatusMgr.cronJob = basejobmanager.InitCron(jobUpdateStatusMgr)

	return jobUpdateStatusMgr
}

func (mgr *jobUpdateStatusManager) Start() {
	mgr.cronJob.Start()
}

func (mgr *jobUpdateStatusManager) GetJobName() string {
	return "Job Update Status ProcessingFile"
}

func (mgr *jobUpdateStatusManager) GetSchedulerConfig() config.SchedulerConfig {
	return mgr.cfg
}

// Execute ...
// Logic:
//  1. Fetch all file that have status = PROCESSING
//  2. Update status of each file
func (mgr *jobUpdateStatusManager) Execute() {
	// 1. Fetch all file that have status = PROCESSING
	ctx := context.Background()
	fpList, err := mgr.fpService.GetListFileByStatuses(ctx, []int16{fileprocessing.StatusProcessing})
	if err != nil {
		logger.Errorf("===== Cannot get list Processing File, got: %v", err)
		return
	}

	// 2. Check empty
	if len(fpList) == 0 {
		logger.InfoT("No PROCESSING file for executing!")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("===== Recovered from a panic %v", r)
			debug.PrintStack()

			fields := map[string]string{
				"Job": mgr.GetJobName(),
			}
			go func(newCtx context.Context) {
				mgr.slackClient.SendError(newCtx, slack.ErrorMsgPanic, nil, fmt.Errorf("%v", r), fields)
			}(context.Background())
		}
	}()

	// 3. Flattening each file
	jobFlatten := newJobUpdateStatus(mgr.fpService, mgr.fprService, mgr.fileService, mgr.cfgMappingService)
	for _, fp := range fpList {
		jobFlatten.UpdateStatus(ctx, *fp)
	}
}
