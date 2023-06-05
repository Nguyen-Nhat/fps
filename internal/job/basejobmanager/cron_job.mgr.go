package basejobmanager

import (
	"github.com/robfig/cron/v3"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

// CronJobManager ... defining a job that run with cron expression
type CronJobManager interface {
	// Start ... Start Cron Job, it calls Execute() inside it
	Start()

	// Execute ... // Execute job directly
	Execute()

	GetJobName() string
	GetSchedulerConfig() config.SchedulerConfig
}

// InitCron ...
func InitCron(jobMgr CronJobManager) *cron.Cron {
	cronExecuteTask := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	jobName := jobMgr.GetJobName()
	id, err := cronExecuteTask.AddFunc(jobMgr.GetSchedulerConfig().Schedule, func() {
		logger.Infof("\n")
		logger.Infof("========== Running %v: Start  ...", jobName)
		jobMgr.Execute()
		logger.Infof("========== Running %v: Finish ...\n", jobName)
	})
	if err != nil {
		logger.Errorf("Init %s failed: %v", jobName, err)
	}

	logger.Infof("Init %s success: ID = %v", jobName, id)

	return cronExecuteTask
}
