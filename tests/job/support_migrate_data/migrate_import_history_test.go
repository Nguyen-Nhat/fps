package support_migrate_data

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xuri/excelize/v2"

	googlesheet "git.teko.vn/loyalty-system/loyalty-file-processing/internal/adapter/googlesheet"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/support_migrate_data"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"go.tekoapis.com/tekone/library/excelreader"
	"go.tekoapis.com/tekone/library/test/monkey"
)

type migrateImportHistorySuite struct {
	suite.Suite
	ctx         context.Context
	job         *support_migrate_data.JobMigrateImportHistoryImpl
	db          *sql.DB
	entClient   *ent.Client
	fixedTime   time.Time
	fpService   fileprocessing.Service
	fileService fileservice.IService

	sheetId            string
	forceDeleteClients []int32
}

func TestMigrateImportHistory(t *testing.T) {
	ts := &migrateImportHistorySuite{}
	ts.ctx = context.Background()
	ts.fixedTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	ts.db, ts.entClient = common.PrepareDatabaseSqlite(ts.ctx, ts.T(), true)
	ts.tearDown()

	fpRepo := fileprocessing.NewRepo(ts.db)
	fpServiceMock := fileprocessing.NewService(fpRepo)

	fileServiceClient := &fileservice.Client{}
	fileServiceMock := fileservice.NewService(fileServiceClient)

	googleSheetClient := &googlesheet.GoogleSheetImpl{}
	monkey.PatchInstanceMethod(reflect.TypeOf(googleSheetClient), "GetXlsxFile",
		func(c *googlesheet.GoogleSheetImpl, ctx context.Context, path string) (*excelize.File, error) {
			fileReader, err := os.ReadFile(ts.sheetId)
			assert.Nil(ts.T(), err)
			file, err := excelize.OpenReader(bytes.NewReader(fileReader))
			assert.Nil(ts.T(), err)
			reader := excelreader.NewExcelReader[support_migrate_data.ImportHistoryDataObject](file, excelreader.ExcelReaderOption{
				HeaderRowIndex: 1,
			})
			dataFrame, err := reader.ScanData(2)
			goldie.New(ts.T()).AssertJson(ts.T(), ts.sheetId, dataFrame)
			return file, nil
		})

	ts.job = &support_migrate_data.JobMigrateImportHistoryImpl{
		SheetId:            ts.sheetId,
		FpService:          fpServiceMock,
		FileService:        fileServiceMock,
		GoogleSheetAdapter: googleSheetClient,
	}

	defer func() {
		ts.tearDown()
	}()

	suite.Run(t, ts)
}

func (ts *migrateImportHistorySuite) tearDown() {
	ts.forceDeleteClients = []int32{}
	common.TruncateAllTables(ts.ctx, ts.db, ts.entClient)
}

func (ts *migrateImportHistorySuite) assert(wantFile string) {
	defer func() {
		ts.tearDown()
	}()
	g := goldie.New(ts.T())
	ts.job.ForceDeleteClients = ts.forceDeleteClients
	err := ts.job.Run(ts.ctx)
	if err != nil {
		g.AssertJson(ts.T(), wantFile, err.Error())
		assert.Equal(ts.T(), strings.HasPrefix(wantFile, "error"), true)
		return
	}
	assert.Nil(ts.T(), err)

	processingFile, err := ts.entClient.ProcessingFile.
		Query().
		All(ts.ctx)
	assert.Nil(ts.T(), err)
	for _, pf := range processingFile {
		pf.CreatedAt = ts.fixedTime
		pf.UpdatedAt = ts.fixedTime
	}

	g.AssertJson(ts.T(), wantFile, processingFile)
}

func (ts *migrateImportHistorySuite) Test200_HappyCase_ThenReturnSuccess() {
	ts.sheetId = "./xlsx/happy_case_migrate_import_history.xlsx"
	ts.assert("happy_case_migrate_import_history")
}

func (ts *migrateImportHistorySuite) Test200_HappyCase_DeleteDataBefore_ThenReturnSuccess() {
	common.MockProcessingFile(ts.ctx, ts.entClient)
	ts.forceDeleteClients = []int32{1, 10, 12, 16}
	ts.sheetId = "./xlsx/happy_case_migrate_import_history.xlsx"
	ts.assert("happy_case_migrate_import_history_delete_data_before")
}

func (ts *migrateImportHistorySuite) Test400_ErrorCase_InvalidData_ThenReturnError() {
	ts.sheetId = "./xlsx/error_case_invalid_data.xlsx"
	ts.assert("error_case_invalid_data")
}

func (ts *migrateImportHistorySuite) Test400_ErrorCase_InvalidTemplate_ThenReturnError() {
	ts.sheetId = "./xlsx/error_case_invalid_template.xlsx"
	ts.assert("error_case_invalid_template")
}
