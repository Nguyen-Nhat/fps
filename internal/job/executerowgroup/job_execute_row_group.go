package executerowgroup

import (
	"context"
	"fmt"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/taskprovider"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

type jobExecuteRowGroup struct {
	fprService        fileprocessingrow.Service
	fpRowGroupService fpRowGroup.Service
}

func newJobExecuteRowGroup(
	fprService fileprocessingrow.Service,
	fpRowGroupService fpRowGroup.Service) *jobExecuteRowGroup {
	return &jobExecuteRowGroup{
		fprService:        fprService,
		fpRowGroupService: fpRowGroupService,
	}
}

// ExecuteRowGroup ...
func (job *jobExecuteRowGroup) ExecuteRowGroup(ctx context.Context, fileID int, taskIndex int, rowGroup fpRowGroup.ProcessingFileRowGroup) {
	startAt := time.Now()
	msgLog := fmt.Sprintf("fileID=%v, taskIndex=%v, rowGroup=%v", fileID, taskIndex, rowGroup.GroupByValue)
	logger.Infof("----- Execute %s", msgLog)

	// 1. Case status = INIT
	if rowGroup.IsInitStatus() {
		rowGroupUpdated, isFinished := job.executeRowGroupCaseStatusIsInit(ctx, fileID, taskIndex, rowGroup, msgLog, startAt)
		if isFinished {
			return // finish
		}

		rowGroup = *rowGroupUpdated
	}

	// 2. Case status = CALLED_API
	if rowGroup.IsCalledAPI() {
		job.executeRowGroupCaseStatusIsCalledAPI(ctx, fileID, taskIndex, rowGroup, msgLog, startAt)
	}
}

// Private method of Job -----------------------------------------------------------------------------------------------

func (job *jobExecuteRowGroup) executeRowGroupCaseStatusIsInit(ctx context.Context, fileID int, taskIndex int, rowGroup fpRowGroup.ProcessingFileRowGroup,
	msgLog string, startAt time.Time) (*fpRowGroup.ProcessingFileRowGroup, bool) {
	groupByValue := rowGroup.GroupByValue

	// 1.1. Get all task that have same groupValue
	tasks, err := job.fprService.GetAllTasksForJobExecuteRowGroup(ctx, fileID, taskIndex, groupByValue)
	if err != nil {
		logger.ErrorT("Failed %s, got error %+v -----> Finish", msgLog, err)
		return nil, true // finish
	}

	// 1.2. If the number of tasks is smaller than totalRows, mean that there is at least one task NOT ready for executing Row Group
	// -> finish this rowGroup, we will execute it at next cycle
	if len(tasks) < int(rowGroup.TotalRows) {
		logger.InfoT("Finish executing %s, because we expected %d tasks but there are only %d tasks are waiting", msgLog, rowGroup.TotalRows, len(tasks))
		return nil, true // finish
	}

	// 1.3. Collect data of all tasks then build request
	configTask, err := mergeTasksToConfigTask(tasks, taskIndex)
	if err != nil {
		updateRequest := toResponseResultRowGroup("", "", err.Error(), false, startAt)
		rowGroupUpdated, err := job.fpRowGroupService.UpdateAfterExecutingByJob(ctx, rowGroup.ID, updateRequest)
		if err != nil {
			return nil, true // finish
		} else {
			return rowGroupUpdated, false
		}
	}

	// 1.4. Call API
	providerClient := taskprovider.NewClientV1()
	curl, responseBody, isSuccess, messageRes := providerClient.Execute(configTask)

	// 1.5. Update task status
	updateRequest := toResponseResultRowGroup(curl, responseBody, messageRes, isSuccess, startAt)
	rowGroupUpdated, err := job.fpRowGroupService.UpdateAfterExecutingByJob(ctx, rowGroup.ID, updateRequest)
	if err != nil {
		logger.ErrorT("Update %v failed", fileprocessingrow.Name())
		return nil, true // finish
	} else {
		return rowGroupUpdated, false
	}
}

func (job *jobExecuteRowGroup) executeRowGroupCaseStatusIsCalledAPI(ctx context.Context, fileID int, taskIndex int, rowGroup fpRowGroup.ProcessingFileRowGroup, msgLog string, startAt time.Time) {
	tasks, err := job.fprService.GetAllTasksForJobExecuteRowGroup(ctx, fileID, taskIndex, rowGroup.GroupByValue)
	if err != nil {
		logger.ErrorT("Failed %s, got error %+v -----> Finish", msgLog, err)
		return // finish
	}

	var pfrIDs []int
	for _, task := range tasks {
		pfrIDs = append(pfrIDs, task.ID)
	}

	var pfrStatus = false
	if rowGroup.IsCalledApiSuccess() {
		pfrStatus = true
	}
	updateRequest := toResponseResult(rowGroup.GroupResponseRaw, rowGroup.ErrorDisplay, pfrStatus, startAt)
	_ = job.fprService.UpdateAfterExecutingByJobForListIDs(ctx, pfrIDs, updateRequest)
}

// Private method ------------------------------------------------------------------------------------------------------

func mergeTasksToConfigTask(tasks []*fileprocessingrow.ProcessingFileRow, taskIndex int) (configloader.ConfigTaskMD, error) {
	// 1. Init
	requestParams := make(map[string]interface{})
	requestBody := make(map[string]interface{})
	var configTask configloader.ConfigTaskMD

	// 2. Explore each task
	for _, task := range tasks {
		// 2.1. Load Data and Mapping
		configMapping, err := converter.StringJsonToStruct("config mapping", task.TaskMapping, configloader.ConfigMappingMD{})
		if err != nil {
			return configloader.ConfigTaskMD{}, fmt.Errorf("internal error")
		}

		// 2.2. Merge request params
		configTask = configMapping.GetConfigTaskMD(taskIndex)
		requestParamsMerged, err := mergeMapInterface(requestParams, configTask.RequestParams)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		}
		requestParams = requestParamsMerged

		// 2.3. Merge request body
		requestBodyMerged, err := mergeMapInterface(requestBody, configTask.RequestBody)
		if err != nil {
			return configloader.ConfigTaskMD{}, err
		}
		requestBody = requestBodyMerged
	}

	// 3. Update request params and request body of configTask with merged data
	configTask.RequestParams = requestParams
	configTask.RequestBody = requestBody

	// 4. Return
	return configTask, nil
}

func toResponseResultRowGroup(curl string, responseRaw string, messageRes string, isSuccess bool, startAt time.Time) fpRowGroup.UpdateAfterExecutingByJob {
	// 1. get status
	var status int16
	if isSuccess {
		status = fpRowGroup.StatusCalledApiSuccess
	} else {
		status = fpRowGroup.StatusCalledApiFail
	}

	// 2. Common value
	return fpRowGroup.UpdateAfterExecutingByJob{
		RequestCurl:  curl,
		ResponseRaw:  responseRaw,
		Status:       status,
		ErrorDisplay: messageRes,
		ExecutedTime: time.Since(startAt).Milliseconds(),
	}
}

func toResponseResult(responseRaw string, messageRes string, isSuccess bool, startAt time.Time) fileprocessingrow.UpdateAfterExecutingByJob {
	// 1. get status
	var status int16
	if isSuccess {
		status = fileprocessingrow.StatusSuccess
	} else {
		status = fileprocessingrow.StatusFailed
	}

	// 2. Common value
	return fileprocessingrow.UpdateAfterExecutingByJob{
		ResponseRaw:  responseRaw,
		Status:       status,
		ErrorDisplay: messageRes,
		ExecutedTime: time.Since(startAt).Milliseconds(),
	}
}
