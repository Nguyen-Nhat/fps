package executetask

import (
	"context"
	"encoding/json"

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
func (job *jobExecuteTask) ExecuteTask(ctx context.Context, fileID int, rowId int32, tasks []*fileprocessingrow.ProcessingFileRow) {
	logger.Infof("----- Execute fileId=%v, rowId=%v, with %v task(s)", fileID, rowId, len(tasks))

	providerClient := taskprovider.NewClientV1()
	previousResponse := make(map[int32]string) // map[task_index]=<response_string>
	for _, task := range tasks {
		// 1. If success, only get response, then go to next task
		if task.IsSuccessStatus() {
			previousResponse[task.TaskIndex] = task.TaskResponseRaw
			continue
		}

		// 2. Execute task
		logger.Infof("---------- Execute fileId=%v, rowId=%v, taskId=%v in ", fileID, rowId, task.TaskIndex)
		requestBody, curl, responseBody, isSuccess, messageRes := providerClient.Execute(int(task.TaskIndex), task.TaskMapping, previousResponse)

		// 3. Update task status and save raw data for tracing
		updateRequest := toResponseResult(requestBody, curl, responseBody, messageRes, isSuccess)
		_, err := job.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
		if err != nil {
			logger.ErrorT("Update %v failed ---> ignore remaining tasks", fileprocessingrow.Name())
			break
		}
		if isSuccess { // task success -> put responseBody to previousResponse (map)
			previousResponse[task.TaskIndex] = responseBody
		} else {
			break // task failed  -> break loop, finish execute task
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func toResponseResult(requestBody map[string]interface{}, curl string, responseBody string, messageRes string, isSuccess bool) fileprocessingrow.UpdateAfterExecutingByJob {
	// 1. Common value
	reqByte, _ := json.Marshal(requestBody)
	updateRequest := fileprocessingrow.UpdateAfterExecutingByJob{
		RequestCurl:  curl,
		RequestRaw:   string(reqByte),
		ResponseRaw:  responseBody,
		Status:       1,
		ErrorDisplay: messageRes,
	}

	// 2. Set status
	if isSuccess {
		updateRequest.Status = fileprocessingrow.StatusSuccess
	} else {
		updateRequest.Status = fileprocessingrow.StatusFailed
	}

	return updateRequest
}
