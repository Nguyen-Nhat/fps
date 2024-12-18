package executetask

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
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/basejobmanager"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/workers"
)

// JobExecuteTaskManager ...
type jobExecuteTaskManager struct {
	cfg     config.SchedulerConfig
	cronJob *cron.Cron
	// services
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	slackClient slack.Client
}

var jobExecuteTaskMgr *jobExecuteTaskManager
var once sync.Once

// NewJobExecuteTaskManager ...
func NewJobExecuteTaskManager(
	cfg config.Config,
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
) basejobmanager.CronJobManager {
	if jobExecuteTaskMgr == nil {
		once.Do(func() {
			slackClient := slack.NewSlackClient(cfg.SlackWebhook)

			jobExecuteTaskMgr = &jobExecuteTaskManager{
				cfg:         cfg.JobConfig.ExecuteTaskConfig,
				fpService:   fpService,
				fprService:  fprService,
				slackClient: slackClient,
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
	var workerPools []*workers.WorkerPool
	for _, file := range fpList {
		// 3.1. Get all task of file.ID, group by rowIndex
		taskGroupByRow, _ := mgr.fprService.GetAllRowsNeedToExecuteByJob(ctx, file.ID, 5000)
		if len(taskGroupByRow) == 0 {
			logger.ErrorT("No row need to execute for fileId=%v", file.ID)
			continue
		}

		// 3.2. Init worker pool
		numDigesters := mgr.cfg.GetNumDigesters(int(file.ClientID))
		logger.Infof("----- Init worker pool with size %d for file %d, total taskGroupByRow is %d\n", numDigesters, file.ID, len(taskGroupByRow))
		jobExecTask := newJobExecuteTask(mgr.fprService)
		workerPool := workers.NewWorkerPool(numDigesters)
		workerPool.Run()

		// 3.3. Handle tasks in each rowIndex by adding it into worker
		for rowId, tasks := range taskGroupByRow {
			tmpRowId := rowId
			tmpTasks := tasks
			workerPool.AddTask(func() {
				jobExecTask.ExecuteTask(ctx, file.ID, tmpRowId, tmpTasks)
			})
		}

		// 3.4. Close worker to stop receiving new tasks
		workerPool.Close()

		// 3.5. Add to list worker pool
		workerPools = append(workerPools, &workerPool)
	}

	// 4. Wait for all tasks completed in each Worker Pool
	for _, wp := range workerPools {
		wp.Wait()
	}

}
