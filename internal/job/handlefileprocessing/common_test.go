package handlefileprocessing

import (
	"database/sql"
	"fmt"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"net/http"
	"net/http/httptest"
)

const (
	issue1297 = "LOY-1297"

	prefixStep2 = "[Job-ExecuteFile][Step2] "
	prefixStep3 = "[Job-ExecuteFile][Step3] "
	prefixStep4 = "[Job-ExecuteFile][Step4] "
)

func jiraTestDetailStep2() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1297},
		Objective:       "Test Job Execute ProcessingFile - Step 2",
		Precondition:    "File already upload to file server, file is Init status",
		Folder:          "HN17/Loyalty File Processing/Job/Execute ProcessingFile/Step 2",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1297},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

func jiraTestDetailStep3() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1297},
		Objective:       "Test Job Execute ProcessingFile - Step 3",
		Precondition:    "File already upload to file server, file is Processing status",
		Folder:          "HN17/Loyalty File Processing/Job/Execute ProcessingFile/Step 3",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1297},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

func jiraTestDetailStep4() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1297},
		Objective:       "Test Job Execute ProcessingFile - Step 4",
		Precondition:    "File already upload to file server",
		Folder:          "HN17/Loyalty File Processing/Job/Execute ProcessingFile/Step 4",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1297},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

// Function ------------------------------------------------------------------------------------------------------------

func initJobExecutor(db *sql.DB) *jobHandleProcessingFileImpl {
	fileServiceConfig := config.FileServiceConfig{
		Endpoint: "https://files.dev.tekoapis.net",
		Paths: config.FileServicePaths{
			UploadDoc: "/upload/doc",
		},
	}

	// file processing
	fpRepo := fileprocessing.NewRepo(db)
	fpService := fileprocessing.NewService(fpRepo)

	// file processing row
	fprRepo := fileprocessingrow.NewRepo(db)
	fprService := fileprocessingrow.NewService(fprRepo)

	// file service
	fileServiceClient := fileservice.NewClient(fileServiceConfig)
	fileService := fileservice.NewService(fileServiceClient)

	// job executor
	return newJobHandleProcessingFile(fpService, fprService, fileService)
}

func mockServer(httpStatus int, payload string) (*httptest.Server, string) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		_, err := w.Write([]byte(payload))
		if err != nil {
			panic("===== Setting up test failed")
		}
	}))

	return server, server.URL
}

func buildMockResponseBody(successCode string, successMsg string) string {
	template := `
   	{
		"code": "%v",
		"message": "%v",
      	"result": {
			"memberId": "1234"
      	}
	}
    `
	responseBodyMock := fmt.Sprintf(template, successCode, successMsg)
	return responseBodyMock
}

func buildMockResponseBodyWithResponseNotMatch(successCode string, successMsg string) string {
	template := `
   	{
		"code": "%v",
		"message": "%v",
      	"result": {
			"memberIdddddd": "1234"
      	}
	}
    `
	responseBodyMock := fmt.Sprintf(template, successCode, successMsg)
	return responseBodyMock
}
