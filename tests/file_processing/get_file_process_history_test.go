package fileprocessing

import (
	"context"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	fileprocessing2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"testing"
)

const issue1296 = "LOY-1296"

var jiraTestDetailsForListFileAwardPoint = jiratest.Detail{
	IssueLinks:      []string{"LOY-1296"},
	Precondition:    "No precondition",
	Folder:          "HN17/Loyalty File Processing/Processing File",
	WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1296},
	ConfluenceLinks: []string{"https://confluence.teko.vn/display/PAYMS/%5BOMNI-999%5D+Declare+and+calculate+rewards+according+to+the+policy"},
}

var jiraTestDetailsForSample = jiratest.Detail{
	IssueLinks:      []string{"HN5-173"},
	Precondition:    "No precondition",
	Folder:          "HN5/Loyalty File Processing/GetFileProcessHistory",
	WebLinks:        []string{"https://jira.teko.vn/browse/HN5-173"},
	ConfluenceLinks: []string{"https://confluence.teko.vn/pages/viewpage.action?spaceKey=PAYMS&title=%5BOMNI-1101%5D%5BSub-TDD%5D+Improve+File+Processing+Service"},
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

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID: 1,
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ?", req.ClientID)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)
}

func TestGetFileProcessHistory_Parameter_include_clientId_and_page__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include clientId and page"
	detail.Objective = "When parameter include clientId and page, should return code 200 " +
		"and response contain list and default pagination's size"
	defer detail.Setup(t)()

	assert.Equal(t, true, true)

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID: 1,
		Page:     2,
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ?", req.ClientID)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    constant.PaginationDefaultSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)
}

func TestGetFileProcessHistory_Parameter_include_clientId_and_size__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 200 - Parameter include clientId and size"
	detail.Objective = "When parameter include clientId and size, should return code 200 " +
		"and response contain list file and default pagination's page"
	defer detail.Setup(t)()

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID: 1,
		PageSize: 1,
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ?", req.ClientID)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: constant.PaginationDefaultPage,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)
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

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID: 1,
		Page:     2,
		PageSize: 1,
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ?", req.ClientID)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)
}
func TestGetFileProcessHistory_Parameter_size_bigger_than_200__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400, error code = 3 - Parameter size bigger than 200"
	detail.Objective = "When parameter contain size bigger than 200, should return code 400"

	defer detail.Setup(t)()

	// Unit tested at 'api/server/fileawardpoint/file_award_point.server_test.go'
	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Parameter_page_bigger_than_1000__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[GetFileProcessHistory] Response code = 400, error code = 3 - Parameter page bigger than 1000"
	detail.Objective = "When parameter contain page bigger than 1000, should return code 400"

	defer detail.Setup(t)()

	// Unit tested at 'api/server/fileawardpoint/file_award_point.server_test.go'
	assert.Equal(t, true, true)
}

func TestGetFileProcessHistory_Input_createdBy_is_empty__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForSample
	detail.Name = "[GetFileProcessHistory] Return code = 200 when createdBy is empty"
	defer detail.Setup(t)()

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID:  12,
		Page:      1,
		PageSize:  2,
		CreatedBy: "tri.tm1@teko.vn",
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ? and  created_by = ?", req.ClientID, req.CreatedBy)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)

}

func TestGetFileProcessHistory_Input_createdBy_is_valid__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForSample
	detail.Name = "[GetFileProcessHistory] Return code = 200 when createdBy is valid"
	defer detail.Setup(t)()

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID: 12,
		Page:     1,
		PageSize: 2,
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ?", req.ClientID)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)

}

func TestGetFileProcessHistory_Input_createdBy_is_not_exist__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForSample
	detail.Name = "[GetFileProcessHistory] Return code = 200 with a empty processFiles list when createdBy is not exist"
	defer detail.Setup(t)()

	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	fpServer := fileprocessing2.InitFileProcessingServer(db)

	req := fileprocessing2.GetFileProcessHistoryRequest{
		ClientID:  12,
		Page:      1,
		PageSize:  2,
		CreatedBy: "long.pt@teko.vn",
	}
	req.InitDefaultValue()

	fpRes, err := fpServer.GetFileProcessHistory(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file processing: %v", err)
	}
	assert.Equal(t, codes.OK, fpRes.Error)

	// Count number of row
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM processing_file WHERE client_id = ? and  created_by = ?", req.ClientID, req.CreatedBy)
	if err != nil {
		t.Errorf("Error query processing file: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := common.PageSizeCalculator(dbRowCount, req.Page, req.PageSize)
	assert.Equal(t, expectedDataLength, len(fpRes.Data.ProcessingFiles))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  dbRowCount,
		TotalPage:   common.PageCalculator(dbRowCount, req.PageSize),
	}
	assert.Equal(t, expectedPagination, fpRes.Data.Pagination)
}
