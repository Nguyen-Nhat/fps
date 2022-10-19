package job

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"github.com/stretchr/testify/assert"
	"testing"
)

const issue1225 = "LOY-1225"

func getDefaultJiraTestDetail() jiratest.Detail {
	return jiratest.Detail{
		IssueLinks:      []string{issue1225},
		Objective:       "Test save award point from file",
		Precondition:    "File already upload to file server",
		Folder:          "HN17/Loyalty File Processing/Job/Check Status FileAwardPoint",
		WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1225},
		ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?pageId=368453857"},
	}
}

func TestJobUpdateStatusFAP__file_noFAPProcessing(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP] do nothing when NO file award point has status=PROCESSING"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSTimeout(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] NOT update member_transaction record when it has status=Processing, call Loyalty timeout, check NOT timeout"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSTimeoutAndCheckTimeout(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] update member_transaction.status=timeout and file result contains this txn when it has status=Processing, call Loyalty timeout, check timeout"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSNoResponse(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] NOT update member_transaction record when it has status=Processing, call Loyalty but no response, check NOT timeout"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSNoResponseAndCheckTimeout(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] update member_transaction.status=Timeout and file result contains this txn when it has status=Processing, call Loyalty but no response, check timeout"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSAndErrorNotSuccess(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] update member_transaction.status=Failed and file result contains this txn when it has status=Processing, Loyalty return error code NOT SUCCESS"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__txn_fapHasTxnProcessing_callLSAndSuccess(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each txn] update member_transaction.status=Success and file result NOT contains this txn when it has status=Processing, Loyalty return error code = SUCCESS"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__file_fapHasTxnProcessing_txnNotFinish(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each file] NOT update file_award_point record when file award point status is Processing, it has txns that have not finished"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__file_fapHasTxnProcessing_allTxnFinish_uploadFileSuccess(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each file] update file_award_point.status=Finish when file award point status is Processing, all txns of it are finished, upload result file success"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}

func TestJobUpdateStatusFAP__file_fapHasTxnProcessing_allTxnFinish_uploadFileFailed(t *testing.T) {
	// Jira Testcase
	detail := getDefaultJiraTestDetail()
	detail.Name = "[jobUpdateStatusFAP][Each file] update file_award_point.status=Finish when file award point status is Processing, all txns of it are finished, upload result file failed"
	defer detail.Setup(t)()

	// todo testcase
	assert.Equal(t, true, true)
}
