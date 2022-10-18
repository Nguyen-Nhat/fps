package awardpoint

import (
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	_ "github.com/go-sql-driver/mysql"
)

const issue1232 = "LOY-1232"

func getJiraTestJobSendAwardPoint() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1232},
		Objective:       "Test job send request award point to loyalty core",
		Precondition:    "Table file_award_pint have record",
		Folder:          "HN17/Loyalty File Processing/Job/Job send reward point",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1232},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

func TestDownloadFileInitAndProcessing__Success(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Download file have status INIT and PROCESSING from file service"
	detail.Objective += "</br>Test download file have status INIT and PROCESSING from file service"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}
func TestValidatePhoneNumberFormat(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Validate phone number format in file"
	detail.Objective += "</br>Test phone number format have 10 digit if start with 0 and 11 digit if start with 84"
	detail.Precondition += "</br>File already downloaded"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestValidateEmptyPhoneNumberFormat(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Validate empty phone number in file"
	detail.Objective += "</br>Test empty phone number in file"
	detail.Precondition += "</br>File already downloaded"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestValidateInactivePhoneNumberFormat(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Inactive phone number after call API core => set error in file result"
	detail.Objective += "</br>Test inactive phone number in file"
	detail.Precondition += "</br>File already downloaded"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestValidatePointNumber(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Validate point number in file, must not have empty, 0 or negative value"
	detail.Objective += "</br>Point must not have empty, 0 or negative value"
	detail.Precondition += "</br>File already downloaded"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestUploadFileResultToFileService(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Create and upload file validate result to file service"
	detail.Objective += "</br>Valid record must save to table member_transaction"
	detail.Precondition += "</br>Validate file process have completed"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestSaveToMemberTransactionDb(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Save member point award record to table member_transaction with status init"
	detail.Objective += "</br>Valid records must save to table member_transaction"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestSendAPIGrantPointToLoyaltyCoreSuccess(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Call to API grant point of Loyalty core success => Change status of member_transaction to processing"
	detail.Objective += "</br>Test Call api grant point of loyalty core success"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}

func TestSendAPIGrantPointToLoyaltyCoreFail(t *testing.T) {
	// Jira Testcase
	detail := getJiraTestJobSendAwardPoint()
	detail.Name = "[JobSendAwardPoint] Call to API grant point of Loyalty core false => Change status of member_transaction to fail and save error message to column error_message"
	detail.Objective += "</br> Call api grant point of loyalty core fail"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init DB & Init Server

	// 2. Mock request

	// 3. Request server

	// 4. Assert
}
