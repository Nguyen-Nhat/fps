package handlefileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
)

const (
	issue1297 = "LOY-1297"

	prefixStep1 = "[Job-ExecuteFile][Step1] "
)

func jiraTestDetailStep1() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1297},
		Objective:       "Test Job Execute ProcessingFile - Step 1",
		Precondition:    "File already upload to file server",
		Folder:          "HN17/Loyalty File Processing/Job/Execute ProcessingFile/Step 1",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1297},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}
