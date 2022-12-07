package awardpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	fileawardpointrepo "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

const issue1220 = "LOY-1220"

func getDefaultJiraTestDetail() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1220},
		Objective:       "Test save award point from file",
		Precondition:    "File already upload to file server",
		Folder:          "HN17/Loyalty File Processing/Award point/Save by file",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1220},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

const (
	FileUrlTest          = "https://storage.googleapis.com/folder/file.xlsx"
	NoteTest             = "Sample not"
	MerchantIDTest int64 = 12387456865126
	UserSubTest          = "1"
	UserNameTest         = "Abc"
	UserEmailTest        = "abc@email.com"
)

func TestUnauthenticated__Return401Unauthenticated(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()
}

func TestValidateFileUrlEmpty__Return400FileUrlIsRequired(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileawardpoint.CreateFileAwardPointDetailRequest{
		MerchantID: MerchantIDTest,
		Note:       NoteTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileawardpoint.CreateFileAwardPointDetailRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileawardpoint.ErrFileUrlRequired)
}

func TestValidateMerchantIdEmpty__Return400MerchantIdIsRequired(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileawardpoint.CreateFileAwardPointDetailRequest{
		FileUrl: FileUrlTest,
		Note:    NoteTest,
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileawardpoint.CreateFileAwardPointDetailRequest{})

	// 4. Assert
	// Expect
	assert.ErrorIs(t, err, fileawardpoint.ErrMerchantIDRequired)
}

func TestValidateNoteMaxLength__Return400NoteTooLong(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	// 2. Mock request
	req := fileawardpoint.CreateFileAwardPointDetailRequest{
		MerchantID: MerchantIDTest,
		FileUrl:    FileUrlTest,
		Note:       strings.Repeat("a", 256),
	}
	jsonBody, _ := json.Marshal(req)
	httpRequest, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(jsonBody))
	httpRequest.Header.Set("Content-Type", "application/json")

	// 3. Request server
	err := render.Bind(httpRequest, &fileawardpoint.CreateFileAwardPointDetailRequest{})

	// 4. Assert
	assert.ErrorIs(t, err, fileawardpoint.ErrNoteOverMaxLength)
}

func TestAllInputValid__Return200FileAwardPointId(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)

	fileAwardPointServer := fileawardpoint.InitFileAwardPointServer(db)

	// 2. Mock request
	ctx = middleware.SetUserToContext(ctx, middleware.User{
		Sub:   UserSubTest,
		Name:  UserNameTest,
		Email: UserEmailTest,
	})
	req := fileawardpoint.CreateFileAwardPointDetailRequest{
		MerchantID: MerchantIDTest,
		FileUrl:    FileUrlTest,
		Note:       NoteTest,
	}

	// 3. Request server
	res, err := fileAwardPointServer.CreateFileAwardPoint(ctx, &req)

	// 4. Assert
	assert.Nil(t, err)
	assert.NotNil(t, res.Data.FileAwardPointID)
	assert.Greater(t, res.Data.FileAwardPointID, 0)
}

func TestAllInputValid__SaveMerchantIdFileUrlNoteToDB(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[saveAwardPoint] Send invalidToken, return 401"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init

	ctx := context.Background()
	db, entClient := common.PrepareDatabaseSqlite(ctx, t)

	fileAwardPointServer := fileawardpoint.InitFileAwardPointServer(db)

	// 2. Mock request
	ctx = middleware.SetUserToContext(ctx, middleware.User{
		Sub:   UserSubTest,
		Name:  UserNameTest,
		Email: UserEmailTest,
	})
	req := fileawardpoint.CreateFileAwardPointDetailRequest{
		MerchantID: MerchantIDTest,
		FileUrl:    FileUrlTest,
		Note:       NoteTest,
	}

	// 3. Request server
	res, err := fileAwardPointServer.CreateFileAwardPoint(ctx, &req)

	// 4. Assert
	assert.Nil(t, err)
	fap, err := entClient.FileAwardPoint.Query().Where(fileawardpointrepo.ID(res.Data.FileAwardPointID)).Only(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res.Data.FileAwardPointID, fap.ID)
	assert.Equal(t, FileUrlTest, fap.FileURL)
	assert.Equal(t, MerchantIDTest, fap.MerchantID)
	assert.Equal(t, NoteTest, fap.Note)
	assert.Equal(t, UserEmailTest, fap.CreatedBy)
	assert.Equal(t, UserEmailTest, fap.UpdatedBy)
}
