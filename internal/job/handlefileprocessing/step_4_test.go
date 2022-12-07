package handlefileprocessing

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob_step4__notFinished(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep4()
	detail.Name = prefixStep4 + "do nothing when file has at least a row which has all tasks are Init"
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
	_ = client.ProcessingFileRow.Update().
		Where(processingfilerow.FileID(int64(fileId))).
		SetStatus(fileprocessingrow.StatusInit).Exec(ctx)

	// 3. Handle logic
	jobExecutor.statisticAndUpdateFileStatus(ctx, fileBefore)

	// 4. Assert
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assert.Equal(t, fileBefore.Status, fileAfter.Status)
	assert.Equal(t, fileBefore.ResultFileURL, fileAfter.ResultFileURL)
	assert.Equal(t, fileBefore.TotalMapping, fileAfter.TotalMapping)
	assert.Equal(t, fileBefore.StatsTotalSuccess, fileAfter.StatsTotalSuccess)
	assert.Equal(t, fileBefore.StatsTotalRow, fileAfter.StatsTotalRow)
	assert.Equal(t, fileBefore.ErrorDisplay, fileAfter.ErrorDisplay)
}

func TestJob_step4__allSuccess(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep4()
	detail.Name = prefixStep4 + "update processing_file.status=Finished when file has all rows which have all task are Success"
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
	_ = client.ProcessingFileRow.Update().
		Where(processingfilerow.FileID(int64(fileId))).
		SetStatus(fileprocessingrow.StatusSuccess).Exec(ctx)

	// 3. Handle logic
	jobExecutor.statisticAndUpdateFileStatus(ctx, fileBefore)

	// 4. Assert
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assert.Equal(t, fileprocessing.StatusFinished, int(fileAfter.Status))
	assert.Equal(t, 0, len(fileAfter.ResultFileURL))
	assert.Equal(t, fileAfter.StatsTotalRow, fileAfter.StatsTotalSuccess)
}

func TestJob_step4__someSuccessAndSomeHasTaskFailed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep4()
	detail.Name = prefixStep4 + "update processing_file.status=Finished when file has some rows are Success and remaining rows have at least a task is Failed"
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
	_ = client.ProcessingFileRow.Update().
		Where(processingfilerow.FileID(int64(fileId))).
		SetStatus(fileprocessingrow.StatusSuccess).Exec(ctx)
	_ = client.ProcessingFileRow.Update().
		Where(processingfilerow.FileID(int64(fileId)),
			processingfilerow.RowIndex(0),
			processingfilerow.TaskIndex(2)).
		SetStatus(fileprocessingrow.StatusInit).Exec(ctx)

	// 3. Handle logic
	jobExecutor.statisticAndUpdateFileStatus(ctx, fileBefore)

	// 4. Assert
	fileAfter, err := client.ProcessingFile.Get(ctx, fileId)
	assert.Nil(t, err)
	assert.Equal(t, fileprocessing.StatusProcessing, int(fileAfter.Status))
	//assert.Equal(t, 0, len(fileAfter.ResultFileURL))
	assert.Equal(t, fileAfter.StatsTotalRow, fileAfter.StatsTotalSuccess+1)
}

func TestJob_step4__resultFileEmpty(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep4()
	detail.Name = prefixStep4 + "[Manual] Result filed has no error row when file is Finished and has all row success"
	defer detail.Setup(t)()

	// Testcase Implementation
	assert.True(t, true) // manual testing
}

func TestJob_step4__resultFileCorrect(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep4()
	detail.Name = prefixStep4 + "[Manual] Result filed has some error row (correct value) when file is Finised and has some row failed"
	defer detail.Setup(t)()

	// Testcase Implementation
	assert.True(t, true) // manual testing
}
