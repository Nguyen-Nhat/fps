package handlefileprocessing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
)

func TestJob_step1__returnAllFileInitOrProcessing(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep1()
	detail.Name = prefixStep1 + "get all file that are Init or Processing status when run step 1"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpRepo := fileprocessing.NewRepo(db)
	fpService := fileprocessing.NewService(fpRepo)

	// 2. Call method that need to test
	fps, err := fpService.GetListFileByStatuses(ctx, []int16{fileprocessing.StatusInit, fileprocessing.StatusProcessing})

	// 3. Assert
	assert.Nil(t, err)
	assert.NotNil(t, fps)
	assert.Len(t, fps, 7, "Expect get 7 records: 6 Init, 1 Processing")
}
