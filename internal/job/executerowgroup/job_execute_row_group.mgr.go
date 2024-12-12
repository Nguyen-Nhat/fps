package executerowgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/adapter/slack"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

// jobExecuteRowGroupManager ...
type jobExecuteRowGroupManager struct {
	cfg     config.SchedulerConfig
	cronJob *cron.Cron
	// services
	fpService         fileprocessing.Service
	fprService        fileprocessingrow.Service
	fpRowGroupService fpRowGroup.Service
	slackClient       slack.Client
}

var jobExecuteRowGroupMgr *jobExecuteRowGroupManager
var once sync.Once

// NewJobExecuteRowGroupManager ...
func NewJobExecuteRowGroupManager(
	cfg config.Config,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fpRowGroupService fpRowGroup.Service,
) basejobmanager.CronJobManager {
	if jobExecuteRowGroupMgr == nil {
		once.Do(func() {
			slackClient := slack.NewSlackClient(cfg.SlackWebhook)

			jobExecuteRowGroupMgr = &jobExecuteRowGroupManager{
				cfg:               cfg.JobConfig.ExecuteGroupTaskConfig,
				fpService:         fpService,
				fprService:        fprService,
				fpRowGroupService: fpRowGroupService,
				slackClient:       slackClient,
			}
		})
	}

	jobExecuteRowGroupMgr.cronJob = basejobmanager.InitCron(jobExecuteRowGroupMgr)

	return jobExecuteRowGroupMgr
}

func (mgr *jobExecuteRowGroupManager) Start() {
	mgr.cronJob.Start()
}

func (mgr *jobExecuteRowGroupManager) GetJobName() string {
	return "Job Execute Group Tasks ProcessingFile"
}

func (mgr *jobExecuteRowGroupManager) GetSchedulerConfig() config.SchedulerConfig {
	return mgr.cfg
}

// Execute ...
// Logic:
//  1. Fetch all file that have status = PROCESSING
//  2. If no file, do nothing
//  3. File all Row Group that have status in (INIT, CALLED_API)
func (mgr *jobExecuteRowGroupManager) Execute() {
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

	// 3. Execute Tasks in each file
	jobExecuteRowGroup := newJobExecuteRowGroup(mgr.fprService, mgr.fpRowGroupService)
	for _, file := range fpList {
		// 3.1. Get all task of file.ID, group by rowIndex
		rowGroupMap, _ := mgr.fpRowGroupService.FindRowGroupForJobExecute(ctx, file.ID)
		if len(rowGroupMap) == 0 {
			logger.ErrorT("No row need to execute for fileId=%v", file.ID)
			continue
		}

		// 3.2. Handle tasks in each rowIndex
		// todo: can use multi thread for improving performance
		for taskIndex, rowGroups := range rowGroupMap {
			for _, rowGroup := range rowGroups {
				jobExecuteRowGroup.ExecuteRowGroup(ctx, file.ID, taskIndex, *rowGroup)
			}
		}

	}

}
