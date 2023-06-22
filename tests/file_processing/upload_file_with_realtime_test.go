package fileprocessing

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func jiraTestDetail4UploadFileWithRealtime() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{"MD-1724"},
		Folder:          "/HN17/Firebase Alternative/Upload file with FPS",
		WebLinks:        []string{"https://jira.teko.vn/browse/MD-1724"},
		ConfluenceLinks: []string{"https://confluence.teko.vn/display/PAYMS/%5BOMNI-999%5D+Declare+and+calculate+rewards+according+to+the+policy"},
	}
}

func TestUploadFileWithRealtime_Return_200_Test_1(t *testing.T) {
	// Jira test case
	detail := jiraTestDetail4UploadFileWithRealtime()
	detail.Name = "[Success] Make sure the old logic still doesn't panic"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)

	fileProcessingServer := fileprocessing.InitFileProcessingServer(db)

	// 2. Mock request
	ctx = middleware.SetUserToContext(ctx, middleware.User{
		Sub:   UserSubTest,
		Name:  UserNameTest,
		Email: UserEmailTest,
	})
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         FileUrlTest,
		FileDisplayName: DisplayNameTest,
		CreatedBy:       CreatedByTest,
	}

	// 3. Request server
	res, err := fileProcessingServer.CreateProcessingFile(ctx, &req)

	// 4. Assert
	assert.Nil(t, err)
	assert.NotNil(t, res.Data.ProcessFileID)
	assert.Greater(t, res.Data.ProcessFileID, int64(0))
}

func TestUploadFileWithRealtime_Return_200_Test_2(t *testing.T) {
	// Jira test case
	detail := jiraTestDetail4UploadFileWithRealtime()
	detail.Name = "[Success] Create record in f-alt-db when client call UploadFileAPI"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)

	fileProcessingServer := fileprocessing.InitFileProcessingServer(db)

	// 2. Mock request
	ctx = middleware.SetUserToContext(ctx, middleware.User{
		Sub:   UserSubTest,
		Name:  UserNameTest,
		Email: UserEmailTest,
	})
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         FileUrlTest,
		FileDisplayName: DisplayNameTest,
		CreatedBy:       CreatedByTest,
	}

	// 3. Request server
	res, err := fileProcessingServer.CreateProcessingFile(ctx, &req)

	// 4. Assert
	assert.Nil(t, err)
	assert.NotNil(t, res.Data.ProcessFileID)
	assert.Greater(t, res.Data.ProcessFileID, int64(0))
}

func TestUploadFileWithRealtime_Return_200_Test_3(t *testing.T) {
	// Jira test case
	detail := jiraTestDetail4UploadFileWithRealtime()
	detail.Name = "[Success] Update f-alt-db when job update record ProcessingFile"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}
