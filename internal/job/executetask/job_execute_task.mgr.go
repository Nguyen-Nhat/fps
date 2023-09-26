package executetask

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/workers"
)

// JobExecuteTaskManager ...
type jobExecuteTaskManager struct {
	cfg     config.SchedulerConfig
	cronJob *cron.Cron
	// services
	fpService  fileprocessing.Service
	fprService fileprocessingrow.Service
}

var jobExecuteTaskMgr *jobExecuteTaskManager
var once sync.Once

// NewJobExecuteTaskManager ...
func NewJobExecuteTaskManager(
	cfg config.SchedulerConfig,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
) basejobmanager.CronJobManager {
	if jobExecuteTaskMgr == nil {
		once.Do(func() {
			jobExecuteTaskMgr = &jobExecuteTaskManager{
				cfg:        cfg,
				fpService:  fpService,
				fprService: fprService,
			}
		})
	}

	jobExecuteTaskMgr.cronJob = basejobmanager.InitCron(jobExecuteTaskMgr)

	return jobExecuteTaskMgr
}

func (mgr *jobExecuteTaskManager) Start() {
	mgr.cronJob.Start()
}

func (mgr *jobExecuteTaskManager) GetJobName() string {
	return "Job Execute Tasks ProcessingFile"
}

func (mgr *jobExecuteTaskManager) GetSchedulerConfig() config.SchedulerConfig {
	return mgr.cfg
}

// Execute ...
// Logic:
//  1. Fetch all file that have status = PROCESSING
//  2. If no file, do nothing
//  3. Group tasks by row_index, then execute each task order by task_id
func (mgr *jobExecuteTaskManager) Execute() {
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

	// 3. Execute Tasks in each file
	jobExecuteTask := newJobExecuteTask(mgr.fprService)
	workerPool := workers.NewWorkerPool(mgr.cfg.NumDigesters)
	workerPool.Run()
	for _, file := range fpList {
		// 3.1. Get all task of file.ID, group by rowIndex
		taskGroupByRow, _ := mgr.fprService.GetAllRowsNeedToExecuteByJob(ctx, file.ID, 5000)
		if len(taskGroupByRow) == 0 {
			logger.ErrorT("No row need to execute for fileId=%v", file.ID)
			continue
		}

		// 3.2. Handle tasks in each rowIndex
		for rowId, tasks := range taskGroupByRow {
			workerPool.AddTask(func() {
				jobExecuteTask.ExecuteTask(ctx, file.ID, rowId, tasks)
			})
		}
		workerPool.Close()
	}
}
