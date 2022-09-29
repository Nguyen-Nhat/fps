package user

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/user"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

const issue1231 = "LOY-1231"

var jiraTestDetailsForCreateUser = jiratest.Detail{
	IssueLinks:      []string{issue1231},
	Objective:       "Test User",
	Precondition:    "No precondition",
	Folder:          "HN17/Loyalty File Processing/User/Create",
	WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1231},
	ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
}

func TestReturnGrantPoint_All_Success(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailsForCreateUser
	detail.Name = "[createUser] return success when input valid"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init
	ctx := context.Background()
	userServer := initUserServerForTesting()

	// 2. Mock request
	req := user.CreateUserRequest{
		Name: "Quy",
	}

	// 3. Request server
	userRes, err := userServer.CreateUser(ctx, &req)

	// 4. Assert
	if err != nil {
		t.Errorf("Error create user: %v", err)
	}
	if req.Name == userRes.Name {
		t.Logf("Create success with name = %v", userRes.Name)
	} else {
		t.Errorf("Create failed with name = %v <> %v", req.Name, userRes.Name)
	}
}
