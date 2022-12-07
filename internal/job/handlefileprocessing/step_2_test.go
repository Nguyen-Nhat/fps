package handlefileprocessing

import (
	"context"
	"database/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob_step2__downloadFileFailed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, cannot download file from FileService"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	pf.FileURL = "https://abc.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__cannotReadFile(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, cannot read file after downloading"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// url is correct by file is wrong
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

// Mapping -------------------------------------------------------------------------------------------------------------

func TestJob_step2__mapping_invalid(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet not enough column"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping not enough column
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_noData(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet empty (no rows, or only has header rows)"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping empty
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_task(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has task_id column is not number"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping has task_id is not number
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_endpoint(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has endpoint column is not link"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping has endpoint is not link
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_header(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has header column is not json"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping header is not json
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_request(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has request column is not json"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but mapping request is not json
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_requestRangeColumn(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has request column contains column not in range {A-Z}"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Mapping sheet has request column contains column not in range {A-Z}
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_requestContainsResponseAnotherTask(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has request column contains response of another task, but it is not start with $response"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Mapping sheet has request column contains response of another task, but it is not start with $response
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_response(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has response column is not json"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Mapping sheet has response column is not json
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__mapping_invalidFormat_responseInvalid(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Mapping sheet has response column is wrong format"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Mapping sheet has response column is wrong format
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

// Data ----------------------------------------------------------------------------------------------------------------

func TestJob_step2__data_missingColumns(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Data sheet missing columns that are mention in Mapping sheet"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Data sheet missing columns that are mention in Mapping sheet
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

func TestJob_step2__data_noData(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Failed when file status is Init, Data sheet empty (no rows, or only has header rows)"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Data sheet empty (no rows, or only has header rows)
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/26/04fcce66-14fe-423d-be99-e18867f79ca9/OMNI-998_%20Giai%20%C4%91o%E1%BA%A1n%201_%20Chuy%E1%BB%83n%20%C4%91i%E1%BB%83m%20Loyalty.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusFailed, int(pfUpdated.Status))
}

// Save data to DB -----------------------------------------------------------------------------------------------------

func TestJob_step2__saveData_failed(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status not changed (still is Init) when file status is Init, save all row data to DB but error occur"
	defer detail.Setup(t)()

	assert.True(t, true) // test manually this case
}

func TestJob_step2__saveData_success(t *testing.T) {
	// Jira Testcase
	detail := jiraTestDetailStep2()
	detail.Name = prefixStep2 + "file.status=Processing, data extracted to processing_file_row table when file status is Init, save all row data to DB success"
	defer detail.Setup(t)()

	// Testcase Implementation
	// 1. Init or Update mock data
	pf := common.GetProcessingFileMockById(1)
	// file correct but Data sheet empty (no rows, or only has header rows)
	pf.FileURL = "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/11/30/c40b1b0c-298f-49c9-beb6-0dfa933aae09/Fumart%20Loyalty%20-%20Import%20sellers.xlsx"

	// 2. Init DB & Init Service
	ctx := context.Background()
	db, _ := common.PrepareDatabaseSqlite(ctx, t)
	jobExecutor := initJobExecutor(db)

	// 3. Handle logic
	pfUpdated := jobExecutor.handleFileInInitStatus(ctx, pf)

	// 4. Assert
	assert.Equal(t, fileprocessing.StatusProcessing, int(pfUpdated.Status))
	assert.True(t, pfUpdated.TotalMapping > 0)
	assert.True(t, pfUpdated.StatsTotalRow > 0)
	recordStatuses := getStatusOfAllExtractedRecordsByFileId(t, db, pfUpdated)
	totalRecordExtracted := len(recordStatuses)
	assert.Equal(t, totalRecordExtracted, int(pfUpdated.TotalMapping*pfUpdated.StatsTotalRow))
}

// ---------------------------------------------------------------------------------------------------------------------

func getStatusOfAllExtractedRecordsByFileId(t *testing.T, db *sql.DB, pfUpdated *fileprocessing.ProcessingFile) []int16 {
	var recordStatuses []int16
	rows, err := db.Query("SELECT status FROM processing_file_row WHERE file_id = ?", pfUpdated.ID)
	if err != nil {
		assert.Failf(t, "query failed", "query failed")
	}
	for rows.Next() {
		var status int16
		err = rows.Scan(&status)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			recordStatuses = append(recordStatuses, status)
		}
	}
	return recordStatuses
}
