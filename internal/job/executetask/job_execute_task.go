package executetask

import (
	"context"
	"fmt"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/taskprovider"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
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
		startAt := time.Now()

		// 1. If success, only get response, then go to next task
		if task.IsSuccessStatus() {
			previousResponse[task.TaskIndex] = task.TaskResponseRaw
			continue
		}

		// 2. Build & Map request
		logger.Infof("---------- Execute fileID=%v, rowID=%v, taskID=%v", fileID, rowID, task.TaskIndex)
		configTask, err := convertConfigMappingAndMapDataFromPreviousResponse(task, previousResponse)
		if err != nil {
			updateRequest := toResponseResult("", "", err.Error(), fileprocessingrow.StatusFailed, startAt)
			_, _ = job.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
			break // task failed  -> break loop, finish execute task
		}

		// 3. Check case row group
		if configTask.RowGroup.IsSupportGrouping() {
			updateRequest := toResponseResult("", "", "", fileprocessingrow.StatusWaitForGrouping, startAt)
			_, _ = job.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
			break // need to handle in Job Execute Row Group -> finish execute this task, and this row
		}

		// 4. Execute task
		curl, responseBody, isSuccess, messageRes := providerClient.Execute(configTask)

		// 5. Update task status and save raw data for tracing
		statusTask := fileprocessingrow.StatusFailed
		if isSuccess {
			statusTask = fileprocessingrow.StatusSuccess
		}
		updateRequest := toResponseResult(curl, responseBody, messageRes, int16(statusTask), startAt)
		_, err = job.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
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

func convertConfigMappingAndMapDataFromPreviousResponse(
	task *fileprocessingrow.ProcessingFileRow,
	previousResponse map[int32]string) (configloader.ConfigTaskMD, error) {
	// 1. Load Data and Mapping
	configMapping, err := converter.StringJsonToStruct("config mapping", task.TaskMapping, configloader.ConfigMappingMD{})
	if err != nil {
		return configloader.ConfigTaskMD{}, fmt.Errorf("failed to load config map")
	}

	// 2. Map data then Build request
	configTask, err := mapDataByPreviousResponse(int(task.TaskIndex), *configMapping, previousResponse)
	return configTask, err
}

func toResponseResult(curl string, responseBody string, messageRes string, status int16, startAt time.Time) fileprocessingrow.UpdateAfterExecutingByJob {
	// 2. Common value
	return fileprocessingrow.UpdateAfterExecutingByJob{
		RequestCurl:  curl,
		ResponseRaw:  responseBody,
		Status:       status,
		ErrorDisplay: messageRes,
		ExecutedTime: time.Since(startAt).Milliseconds(),
	}
}
