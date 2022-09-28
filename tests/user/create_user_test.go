package user

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
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
	detail := jiraTestDetailsForCreateUser
	detail.Name = "[createUser] return success when input valid"
	defer detail.Setup(t)()

	t.Log("OK")
}
