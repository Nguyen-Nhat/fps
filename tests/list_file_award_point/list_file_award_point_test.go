package list_file_award_point

import (
	"context"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	fileawardpoint2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/jiratest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"math"
	"testing"
)

const issue1221 = "LOY-1221"

var jiraTestDetailsForListFileAwardPoint = jiratest.Detail{
	IssueLinks:      []string{"LOY-1221"},
	Precondition:    "No precondition",
	Folder:          "HN17/Loyalty File Processing/API List File Award Point",
	WebLinks:        []string{"https://jira.teko.vn/browse/" + issue1221},
	ConfluenceLinks: []string{"https://confluence.teko.vn/display/PAYMS/%5Bv4%5D+Upsert+Credential"},
}

func pageCalculator(total int, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func pageSizeCalculator(total int, page int, pageSize int) int {
	var expectedPageSize int
	if (page-1)*pageSize > total {
		expectedPageSize = 0
	} else if temp := page * pageSize; temp > total {
		expectedPageSize = total - temp + pageSize
	} else {
		expectedPageSize = pageSize
	}
	return expectedPageSize
}

func TestAPIListFile_Parameter_empty__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 200 - Parameter empty"
	detail.Objective = "When parameter empty, should return code 200 and response contain list all file"
	defer detail.Setup(t)()

	ctx := context.Background()
	db := common.PrepareDatabase(ctx)
	fapServer := fileawardpoint2.InitFileAwardPointServer(db)

	req := fileawardpoint.GetListFileAwardPointDTO{}
	req.InitDefaultValue()

	fapRes, err := fapServer.GetList(ctx, &req)

	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	assert.Equal(t, codes.OK, fapRes.Error)

	// Count number of row in table file_award_point
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM file_award_point")
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := pageSizeCalculator(dbRowCount, req.Page, req.Size)
	assert.Equal(t, expectedDataLength, len(fapRes.Data.FileAwardPoints))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    expectedDataLength,
		TotalItems:  dbRowCount,
		TotalPage:   pageCalculator(dbRowCount, req.Size),
	}
	assert.Equal(t, expectedPagination, fapRes.Data.Pagination)
}

func TestAPIListFile_Parameter_include_only_merchantId__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 200 - Parameter only include merchantId"
	detail.Objective = "When parameter only include merchantId, should return code 200 and response contain list file correspond to merchantId"

	defer detail.Setup(t)()

	ctx := context.Background()
	db := common.PrepareDatabase(ctx)
	fapServer := fileawardpoint2.InitFileAwardPointServer(db)

	req := fileawardpoint.GetListFileAwardPointDTO{
		MerchantId: 1,
	}
	req.InitDefaultValue()

	fapRes, err := fapServer.GetList(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	assert.Equal(t, codes.OK, fapRes.Error)

	// Count number of row in table file_award_point
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM file_award_point WHERE merchant_id = ?", req.MerchantId)
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := pageSizeCalculator(dbRowCount, req.Page, req.Size)
	assert.Equal(t, expectedDataLength, len(fapRes.Data.FileAwardPoints))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    expectedDataLength,
		TotalItems:  dbRowCount,
		TotalPage:   pageCalculator(dbRowCount, req.Size),
	}
	assert.Equal(t, expectedPagination, fapRes.Data.Pagination)
}

func TestAPIListFile_Parameter_include_all__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 200 - Parameter include merchantId, page, size"
	detail.Objective = "When parameter include merchantId, page, size, should return code 200 and response contain list file "

	defer detail.Setup(t)()

	ctx := context.Background()
	db := common.PrepareDatabase(ctx)
	fapServer := fileawardpoint2.InitFileAwardPointServer(db)

	req := fileawardpoint.GetListFileAwardPointDTO{
		MerchantId: 1,
		Page:       1,
		Size:       2,
	}

	fapRes, err := fapServer.GetList(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	assert.Equal(t, codes.OK, fapRes.Error)

	// Count number of row in table file_award_point
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM file_award_point WHERE merchant_id = ?", req.MerchantId)
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := pageSizeCalculator(dbRowCount, req.Page, req.Size)
	assert.Equal(t, expectedDataLength, len(fapRes.Data.FileAwardPoints))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    expectedDataLength,
		TotalItems:  dbRowCount,
		TotalPage:   pageCalculator(dbRowCount, req.Size),
	}
	assert.Equal(t, expectedPagination, fapRes.Data.Pagination)
}

func TestAPIListFile_Parameter_only_include_page_and_size__Should_return_code_200(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 200 - Parameter only include page and size"
	detail.Objective = "When parameter only include page, size, should return code 200 and response contain"

	defer detail.Setup(t)()

	ctx := context.Background()
	db := common.PrepareDatabase(ctx)
	fapServer := fileawardpoint2.InitFileAwardPointServer(db)

	req := fileawardpoint.GetListFileAwardPointDTO{
		Page: 2,
		Size: 1,
	}

	fapRes, err := fapServer.GetList(ctx, &req)
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	assert.Equal(t, codes.OK, fapRes.Error)

	// Count number of row in table file_award_point
	var dbRowCount int
	rows, err := db.Query("SELECT COUNT(*) FROM file_award_point")
	if err != nil {
		t.Errorf("Error get list file award point: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbRowCount)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedDataLength := pageSizeCalculator(dbRowCount, req.Page, req.Size)
	assert.Equal(t, expectedDataLength, len(fapRes.Data.FileAwardPoints))

	// Assert pagination
	expectedPagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    expectedDataLength,
		TotalItems:  dbRowCount,
		TotalPage:   pageCalculator(dbRowCount, req.Size),
	}
	assert.Equal(t, expectedPagination, fapRes.Data.Pagination)
}

func TestAPIListFile_Parameter_size_bigger_than_200__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 400, error code = 3 - Parameter size bigger than 200"
	detail.Objective = "When parameter contain size bigger than 200, should return code 400"

	defer detail.Setup(t)()

	// Unit tested at 'api/server/fileawardpoint/file_award_point.server_test.go'
	t.Log("OK")
}

func TestAPIListFile_Parameter_page_bigger_than_1000__Should_return_code_400(t *testing.T) {
	detail := jiraTestDetailsForListFileAwardPoint
	detail.Name = "[ListFileAwardPoint] Response code = 400, error code = 3 - Parameter page bigger than 1000"
	detail.Objective = "When parameter contain page bigger than 1000, should return code 400"

	defer detail.Setup(t)()

	// Unit tested at 'api/server/fileawardpoint/file_award_point.server_test.go'
	t.Log("OK")
}
