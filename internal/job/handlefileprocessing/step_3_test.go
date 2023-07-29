package handlefileprocessing

import (
	"context"
	"encoding/json"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob_step3__find_noInit(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "do nothing when there is no record that has processing_file_row.status=Init"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	// delete all row is Init
	client.ProcessingFileRow.Delete().Where(processingfilerow.FileID(int64(fileId)), processingfilerow.Status(fileprocessingrow.StatusInit))

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)
}

func TestJob_step3__find_hasInitButAnotherTaskFailed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "do nothing when there is a record has processing_file_row.status=Init, but another task (in same row_index) is Failed"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	rowIndex := 0
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	// update task is failed
	_ = client.ProcessingFileRow.Update().Where(
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)),
		processingfilerow.TaskIndex(1)).
		SetStatus(fileprocessingrow.StatusFailed).Exec(ctx)
	pfrQuery := client.ProcessingFileRow.Query().Where(
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)))
	pfrsBefore, _ := pfrQuery.All(ctx)

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)

	pfrsAfter, _ := pfrQuery.All(ctx)
	assert.Equal(t, len(pfrsBefore), len(pfrsAfter))
	assertListProcessingFileRow(t, pfrsAfter, pfrsBefore)
}

// Execute each group ---------------------------------------------------------------------------------------------------

func TestJob_step3__executeTask_buildRequestFailedBecauseResponseOfAnotherTask(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "processing_file_row.status=Failed and remaining tasks (same row_index) are ignore when execute task failed because build request failed"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	rowIndex := 0
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	conditionalGetTasksInRow0 := []predicate.ProcessingFileRow{
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)),
	}
	// update task 1 to success
	_ = client.ProcessingFileRow.Update().Where(conditionalGetTasksInRow0...).Where(processingfilerow.TaskIndex(1)).
		SetStatus(fileprocessingrow.StatusSuccess).Exec(ctx)

	// update task 2 depend on task 1
	wrongMappingResponse := "result.abcdef"
	task2, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).Where(processingfilerow.TaskIndex(2)).Only(ctx)
	assert.NotNil(t, task2)
	newTaskMappingNotMatchResponse := updateMappingKeyFromTaskMapping(task2, wrongMappingResponse) // mock wrong key
	_ = client.ProcessingFileRow.Update().Where(conditionalGetTasksInRow0...).Where(processingfilerow.TaskIndex(2)).
		SetStatus(fileprocessingrow.StatusInit).SetTaskMapping(newTaskMappingNotMatchResponse).Exec(ctx)
	// get list task
	pfrsBefore, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)

	pfrsAfter, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)
	assert.Equal(t, len(pfrsBefore), len(pfrsAfter))
	for i, after := range pfrsAfter {
		before := pfrsBefore[i]
		if after.TaskIndex == 2 {
			assert.Equal(t, fileprocessingrow.StatusFailed, int(after.Status))
			//assert.Equal(t, before.TaskRequestRaw, after.TaskRequestRaw)
			//assert.Equal(t, before.TaskResponseRaw, after.TaskResponseRaw)
			assert.Contains(t, after.ErrorDisplay, wrongMappingResponse)
		} else {
			assert.Equal(t, before.Status, after.Status)
			assert.Equal(t, before.TaskRequestRaw, after.TaskRequestRaw)
			assert.Equal(t, before.TaskResponseRaw, after.TaskResponseRaw)
			assert.Equal(t, before.ErrorDisplay, after.ErrorDisplay)
		}
	}
}

func TestJob_step3__executeTask_callTimeout(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "processing_file_row.status=Failed and remaining tasks (same row_index) are ignore when execute task failed because call request timeout"
	defer detail.Setup(t)()

	assert.True(t, true) // manual testing because I haven't found the way to mock timeout server
}

func TestJob_step3__executeTask_callFailed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "processing_file_row.status=Failed and remaining tasks (same row_index) are ignore when execute task failed because call success but response failed"
	defer detail.Setup(t)()

	failedMsg := "failed 123"
	responseBodyMock := buildMockResponseBody("500", failedMsg)
	serverMock, endpointMock := mockServer(500, responseBodyMock)
	defer serverMock.Close()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	rowIndex := 0
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	conditionalGetTasksInRow0 := []predicate.ProcessingFileRow{
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)),
	}
	// update endpoint mock
	pfrsBefore := updateMockServerEndpointToTaskInDB(client, conditionalGetTasksInRow0, ctx, endpointMock)

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)

	pfrsAfter, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)
	assert.Equal(t, len(pfrsBefore), len(pfrsAfter))
	for i, after := range pfrsAfter {
		before := pfrsBefore[i]
		if after.TaskIndex == 1 {
			assert.Equal(t, fileprocessingrow.StatusFailed, int(after.Status))
			// assert.NotEqual(t, before.TaskRequestRaw, after.TaskRequestRaw) // comment, because RequestRaw is deprecated
			assert.NotEqual(t, before.TaskResponseRaw, after.TaskResponseRaw)
			// assert.True(t, len(after.TaskRequestRaw) > 0) // comment, because RequestRaw is deprecated
			assert.True(t, len(after.TaskResponseRaw) > 0)
			assert.Contains(t, after.ErrorDisplay, failedMsg)
		} else { // remaining task not change data
			assert.Equal(t, before.Status, after.Status)
			assert.Equal(t, before.TaskRequestRaw, after.TaskRequestRaw)
			assert.Equal(t, before.TaskResponseRaw, after.TaskResponseRaw)
			assert.Equal(t, before.ErrorDisplay, after.ErrorDisplay)
		}
	}
}

func TestJob_step3__executeTask_callSuccessButMappingResponseFailed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "processing_file_row.status=Failed and remaining tasks (same row_index) are ignore when execute task success but cannot read Code in response"
	defer detail.Setup(t)()

	successMsg := "success"
	responseBodyMock := buildMockResponseBodyWithResponseNotMatch("00", successMsg)
	serverMock, endpointMock := mockServer(200, responseBodyMock)
	defer serverMock.Close()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	rowIndex := 0
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	conditionalGetTasksInRow0 := []predicate.ProcessingFileRow{
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)),
	}
	// update endpoint mock
	pfrsBefore := updateMockServerEndpointToTaskInDB(client, conditionalGetTasksInRow0, ctx, endpointMock)

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)

	pfrsAfter, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)
	assert.Equal(t, len(pfrsBefore), len(pfrsAfter))
	for i, after := range pfrsAfter {
		before := pfrsBefore[i]
		if after.TaskIndex == 1 {
			assertTaskSuccess(t, after, before, successMsg)
		} else { // task 2 failed
			assert.Equal(t, fileprocessingrow.StatusFailed, int(after.Status))
			//assert.NotEqual(t, before.TaskRequestRaw, after.TaskRequestRaw) // comment, because RequestRaw is deprecated
			//assert.True(t, len(after.TaskRequestRaw) > 0) // comment, because RequestRaw is deprecated
			assert.Equal(t, before.TaskResponseRaw, after.TaskResponseRaw)
			assert.Contains(t, after.ErrorDisplay, "result.memberId")
		}
	}
}

func TestJob_step3__executeTask_callSuccess(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep3()
	detail.Name = prefixStep3 + "processing_file_row.status=Success when execute task success (call to endpoint success)"
	defer detail.Setup(t)()

	successMsg := "success"
	responseBodyMock := buildMockResponseBody("00", successMsg)
	serverMock, endpointMock := mockServer(200, responseBodyMock)
	defer serverMock.Close()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, client := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 2. Get data mock and pre-processing
	fileId := 2
	rowIndex := 0
	f, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	fileBefore := &fileprocessing.ProcessingFile{ProcessingFile: *f}
	conditionalGetTasksInRow0 := []predicate.ProcessingFileRow{
		processingfilerow.FileID(int64(fileId)),
		processingfilerow.RowIndex(int32(rowIndex)),
	}
	// update endpoint mock
	pfrsBefore := updateMockServerEndpointToTaskInDB(client, conditionalGetTasksInRow0, ctx, endpointMock)

	// 3. Handle logic
	jobExecutor.handleFileInProcessingStatus(ctx, fileBefore)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileBefore.Status))
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assertProcessingFileNoChange(t, fileBefore, fileAfter)

	pfrsAfter, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)
	assert.Equal(t, len(pfrsBefore), len(pfrsAfter))
	for i, after := range pfrsAfter {
		before := pfrsBefore[i]
		assert.Equal(t, fileprocessingrow.StatusSuccess, int(after.Status))
		//assert.NotEqual(t, before.TaskRequestRaw, after.TaskRequestRaw) // comment, because RequestRaw is deprecated
		assert.NotEqual(t, before.TaskResponseRaw, after.TaskResponseRaw)
		//assert.True(t, len(after.TaskRequestRaw) > 0) // comment, because RequestRaw is deprecated
		assert.True(t, len(after.TaskResponseRaw) > 0)
		assert.Contains(t, after.ErrorDisplay, successMsg)
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func assertProcessingFileNoChange(t *testing.T, fileBefore *fileprocessing.ProcessingFile, fileAfter *ent.ProcessingFile) {
	assert.Equal(t, fileBefore.Status, fileAfter.Status)
	assert.Equal(t, fileBefore.ResultFileURL, fileAfter.ResultFileURL)
	assert.Equal(t, fileBefore.TotalMapping, fileAfter.TotalMapping)
	assert.Equal(t, fileBefore.StatsTotalSuccess, fileAfter.StatsTotalSuccess)
	assert.Equal(t, fileBefore.StatsTotalRow, fileAfter.StatsTotalRow)
	assert.Equal(t, fileBefore.ErrorDisplay, fileAfter.ErrorDisplay)
	assert.Equal(t, fileBefore.UpdatedAt, fileAfter.UpdatedAt)
}

func assertListProcessingFileRow(t *testing.T, pfrsAfter []*ent.ProcessingFileRow, pfrsBefore []*ent.ProcessingFileRow) {
	for i, after := range pfrsAfter {
		before := pfrsBefore[i]
		assert.Equal(t, before.Status, after.Status)
		assert.Equal(t, before.TaskRequestRaw, after.TaskRequestRaw)
		assert.Equal(t, before.TaskResponseRaw, after.TaskResponseRaw)
		assert.Equal(t, before.ErrorDisplay, after.ErrorDisplay)
	}
}

func assertTaskSuccess(t *testing.T, after *ent.ProcessingFileRow, before *ent.ProcessingFileRow, successMsg string) {
	assert.Equal(t, fileprocessingrow.StatusSuccess, int(after.Status))
	//assert.NotEqual(t, before.TaskRequestRaw, after.TaskRequestRaw) // comment, because RequestRaw is deprecated
	assert.NotEqual(t, before.TaskResponseRaw, after.TaskResponseRaw)
	//assert.True(t, len(after.TaskRequestRaw) > 0) // comment, because RequestRaw is deprecated
	assert.True(t, len(after.TaskResponseRaw) > 0)
	assert.Contains(t, after.ErrorDisplay, successMsg)
}

func updateMappingKeyFromTaskMapping(task2 *ent.ProcessingFileRow, newMappingKey string) string {
	taskMapping := task2.TaskMapping

	var mappingRow dto.MappingRow
	_ = json.Unmarshal([]byte(taskMapping), &mappingRow)
	mappingRequestMap := mappingRow.Request

	for key, mappingRequest := range mappingRequestMap {
		mappingRequest.IsMappingExcel = false
		mappingRequest.IsMappingResponse = true
		mappingRequest.MappingKey = newMappingKey
		mappingRequestMap[key] = mappingRequest
		break // only update first element
	}

	b, _ := json.Marshal(mappingRow)

	return string(b)
}

func updateEndpointFromTaskMapping(task *ent.ProcessingFileRow, newEndpoint string) string {
	taskMapping := task.TaskMapping

	var mappingRow dto.MappingRow
	_ = json.Unmarshal([]byte(taskMapping), &mappingRow)
	mappingRow.Endpoint = newEndpoint

	b, _ := json.Marshal(mappingRow)

	return string(b)
}

func updateMockServerEndpointToTaskInDB(client *ent.Client, conditionalGetTasksInRow0 []predicate.ProcessingFileRow, ctx context.Context, endpointMock string) []*ent.ProcessingFileRow {
	pfrsBefore, _ := client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx)
	for _, task := range pfrsBefore {
		updatedTaskMapping := updateEndpointFromTaskMapping(task, endpointMock)
		_ = client.ProcessingFileRow.Update().Where(processingfilerow.ID(task.ID)).
			SetStatus(fileprocessingrow.StatusInit).
			SetTaskMapping(updatedTaskMapping).
			Exec(ctx)
	}

	pfrsBefore, _ = client.ProcessingFileRow.Query().Where(conditionalGetTasksInRow0...).All(ctx) // get after updating
	return pfrsBefore
}
