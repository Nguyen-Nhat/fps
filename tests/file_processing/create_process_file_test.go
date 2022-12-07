package fileprocessing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	fileprocessingrepo "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfile"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

const issue1295 = "LOY-1295"

const (
	ClientIDTest    int32 = 12345678
	DisplayNameTest       = "a.xlsx"
	FileUrlTest           = "https://storage.googleapis.com/folder/file.xlsx"
	CreatedByTest         = "nguyen.ng@teko.vn"
	UserSubTest           = "1"
	UserNameTest          = "Abc"
	UserEmailTest         = "abc@email.com"
)

func getDefaultJiraTestDetail() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1295},
		Objective:       "Test API create process file",
		Precondition:    "File already upload to file server",
		Folder:          "HN17/Loyalty File Processing/Process file/Save by file",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1295},
		ConfluenceLinks: []string{"https://confluence.teko.vn/display/PAYMS/%5BOMNI-999%5D+Declare+and+calculate+rewards+according+to+the+policy"},
	}
}

func TestUnauthenticated__Return401Unauthenticated(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send invalidToken, return 401"
	defer detail.Setup(t)()
}

func TestValidateFileUrlEmpty__Return400FileUrlIsRequired(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send empty fileUrl, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         "",
		FileDisplayName: DisplayNameTest,
		CreatedBy:       CreatedByTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrFileUrlRequired)
}

func TestValidateClientIdEmpty__Return400ClientIdIsRequired(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send clientId empty, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		FileURL:         FileUrlTest,
		FileDisplayName: DisplayNameTest,
		CreatedBy:       CreatedByTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrClientIDRequired)
}

func TestValidateCreatedByEmpty__Return400CreatedByIsRequired(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send createdBy empty, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         FileUrlTest,
		FileDisplayName: DisplayNameTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrCreatedByRequired)
}

func TestValidateDisplayNameOverMaxLength__Return400OverMaxLength(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send displayName over max length 255 character, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         FileUrlTest,
		FileDisplayName: strings.Repeat("a", 256),
		CreatedBy:       CreatedByTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrDisplayNameOverMaxLength)
}
func TestValidateCreatedByOverMaxLength__Return400OverMaxLength(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send created by over max length 255 character, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         FileUrlTest,
		FileDisplayName: DisplayNameTest,
		CreatedBy:       strings.Repeat("a", 256),
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrCreatedByOverMaxLength)
}
func TestValidateFileUrlOverMaxLength__Return400OverMaxLength(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send fileUrl over max length 255 character, return 400"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileprocessing.CreateFileProcessingRequest{
		ClientID:        ClientIDTest,
		FileURL:         strings.Repeat("a", 256),
		FileDisplayName: DisplayNameTest,
		CreatedBy:       CreatedByTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileprocessing.CreateFileProcessingRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileprocessing.ErrFileUrlOverMaxLength)
}

func TestAllInputValid__Return200FileProcessingId(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send all valid input, return 200"
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

func TestAllInputValid__SaveClientIdFileUrlDisplayNameToDB(t *testing.T) {
	// Jira test case
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveProcessingFile] Send all valid input, save clientId, fileUrl, displayName to DB"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init
	ctx := context.Background()
	db, entClient := common.PrepareDatabaseSqlite(ctx, t)

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
	fp, err := entClient.ProcessingFile.Query().Where(fileprocessingrepo.ID(int(res.Data.ProcessFileID))).Only(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res.Data.ProcessFileID, int64(fp.ID))
	assert.Equal(t, FileUrlTest, fp.FileURL)
	assert.Equal(t, ClientIDTest, fp.ClientID)
	assert.Equal(t, DisplayNameTest, fp.DisplayName)
	assert.Equal(t, CreatedByTest, fp.CreatedBy)
}
