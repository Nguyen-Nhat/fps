package executetask

import (
	"context"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/taskprovider"
)

type jobExecuteTask struct {
	fprService fileprocessingrow.Service
}

func newJobExecuteTask(fprService fileprocessingrow.Service) *jobExecuteTask {
	return &jobExecuteTask{fprService: fprService}
}

// ExecuteTask ...
// For each tasks in same rowIndex:
//  1. If task succeed, only get response, then go to next task
//  2. Else, Execute task (failed task is not in this context, so only Init task)
//  3. Update task status and save raw data for tracing
func (job *jobExecuteTask) ExecuteTask(ctx context.Context, fileID int, rowID int32, tasks []*fileprocessingrow.ProcessingFileRow) {
	logger.Infof("----- Execute fileID=%v, rowID=%v, with %v task(s)", fileID, rowID, len(tasks))

	providerClient := taskprovider.NewClientV1()
	previousResponse := make(map[int32]string) // map[task_index]=<response_string>
	for _, task := range tasks {
		// 1. If success, only get response, then go to next task
		if task.IsSuccessStatus() {
			previousResponse[task.TaskIndex] = task.TaskResponseRaw
			continue
		}

		// 2. Execute task
		startAt := time.Now()
		logger.Infof("---------- Execute fileID=%v, rowID=%v, taskID=%v", fileID, rowID, task.TaskIndex)
		curl, responseBody, isSuccess, messageRes := providerClient.Execute(int(task.TaskIndex), task.TaskMapping, previousResponse)

		// 3. Update task status and save raw data for tracing
		updateRequest := toResponseResult(curl, responseBody, messageRes, isSuccess, startAt)
		_, err := job.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
		if err != nil {
			logger.ErrorT("Update %v failed ---> ignore remaining tasks", fileprocessingrow.Name())
			break // error occur  -> break loop, finish execute task
		}
		if isSuccess { // task success -> put responseBody to previousResponse (map)
			previousResponse[task.TaskIndex] = responseBody
		} else {
			logger.ErrorT("Execute task failed ---> ignore remaining tasks")
			break // task failed  -> break loop, finish execute task
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func toResponseResult(curl string, responseBody string, messageRes string, isSuccess bool, startAt time.Time) fileprocessingrow.UpdateAfterExecutingByJob {
	// 1. get status
	var status int16
	if isSuccess {
		status = fileprocessingrow.StatusSuccess
	} else {
		status = fileprocessingrow.StatusFailed
	}

	// 2. Common value
	return fileprocessingrow.UpdateAfterExecutingByJob{
		RequestCurl:  curl,
		ResponseRaw:  responseBody,
		Status:       status,
		ErrorDisplay: messageRes,
		ExecutedTime: time.Since(startAt).Milliseconds(),
	}
}
