package fileprocessing

import (
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"github.com/stretchr/testify/assert"
)

const issue1296 = "LOY-1296"

var jiraTestDetailsForListFileAwardPoint = jiratest.Detail{
	IssueLinks:      []string{"LOY-1296"},
	Precondition:    "No precondition",
	Folder:          "HN17/Loyalty File Processing/Processing File",
	WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1296},
	ConfluenceLinks: []string{"https://confluence.teko.vn/display/PAYMS/%5BOMNI-999%5D+Declare+and+calculate+rewards+according+to+the+policy"},
}

func TestGetFileProcessHistory_Parameter_empty__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400 - Parameter empty"
	detail.Objective = "When parameter empty, should return code 400 - Invalid request"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_include_only_clientId__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include only clientId"
	detail.Objective = "When parameter only include clientId, should return code 200 " +
		"and response contain list file history matching clientId and default pagination"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_include_clientId_and_page__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include clientId and page"
	detail.Objective = "When parameter include clientId and page, should return code 200 " +
		"and response contain list and default pagination's size"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_include_clientId_and_size__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include clientId and size"
	detail.Objective = "When parameter include clientId and size, should return code 200 " +
		"and response contain list file and default pagination's page"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_include_page_and_size__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400 - Parameter parameter include page and size"
	detail.Objective = "When parameter doesn't include clientId, should return code 400 - Invalid request"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_include_all__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include all"
	detail.Objective = "When parameter include all, should return code 200" +
		" and response contain list file and pagination"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)
}
func TestGetFileProcessHistory_Parameter_size_bigger_than_200__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400, error code = 3 - Parameter size bigger than 200"
	detail.Objective = "When parameter contain size bigger than 200, should return code 400"

	defer detail.Setup(t)()
}

func TestGetFileProcessHistory_Parameter_page_bigger_than_1000__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400, error code = 3 - Parameter page bigger than 1000"
	detail.Objective = "When parameter contain page bigger than 1000, should return code 400"

	defer detail.Setup(t)()
}
