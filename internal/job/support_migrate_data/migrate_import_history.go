package support_migrate_data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	googlesheet "git.teko.vn/loyalty-system/loyalty-file-processing/internal/adapter/googlesheet"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"go.tekoapis.com/tekone/library/excelreader"
)

const (
	dataStartAtRow = 2
	resultHeader   = "result"
)

type JobMigrateImportHistoryImpl struct {
	SheetId            string
	ForceDeleteClients []int32
	FpService          fileprocessing.Service
	FileService        fileservice.IService
	GoogleSheetAdapter googlesheet.GoogleSheetAdapter
}

type ImportHistoryDataObject struct {
	ClientId          int32    `name:"client_id"`
	TenantId          string   `name:"tenant_id"`
	MerchantId        string   `name:"merchant_id"`
	DisplayName       string   `name:"display_name"`
	ExtFileRequest    string   `name:"ext_file_request"`
	FileUrl           string   `name:"file_url"`
	ResultFileUrl     string   `name:"result_file_url"`
	Status            int32    `name:"status"`
	StatsTotalRow     int32    `name:"stats_total_row"`
	StatsTotalSuccess int32    `name:"stats_total_success"`
	ErrorDisplay      string   `name:"error_display"`
	CreatedBy         string   `name:"created_by"`
	CreatedAt         int64    `name:"created_at"`
	UpdatedAt         int64    `name:"updated_at"`
	Index             int      `xlsx:"ind"`
	ErrList           []string `xlsx:"err"`
}

func NewJobMigrateImportHistory(
	sheetId string,
	forceDeleteClients []int32,
	fpService fileprocessing.Service,
	fileService fileservice.IService) (*JobMigrateImportHistoryImpl, error) {
	return &JobMigrateImportHistoryImpl{
		SheetId:            sheetId,
		ForceDeleteClients: forceDeleteClients,
		FpService:          fpService,
		FileService:        fileService,
		GoogleSheetAdapter: googlesheet.NewGoogleSheetClient(),
	}, nil
}

func (s *JobMigrateImportHistoryImpl) Run(ctx context.Context) (err error) {
	startTime := time.Now()
	logger.Infof("JobMigrateImportHistoryImpl | Start job at %v", startTime)
	defer func() {
		r := recover()
		msg := constant.EmptyString
		if r != nil {
			msg += fmt.Sprintf("panic: %s ", r)
		}
		if err != nil {
			msg += fmt.Sprintf("err: %s", err.Error())
		}
		if msg == constant.EmptyString {
			msg = constant.MessageSuccess
		}
		logger.Infof("JobMigrateImportHistoryImpl | End job at %v, duration: %v, result: %s", time.Now(), time.Since(startTime), msg)
	}()

	file, err := s.GoogleSheetAdapter.GetXlsxFile(ctx, s.SheetId)
	reader := excelreader.NewExcelReader[ImportHistoryDataObject](file, excelreader.ExcelReaderOption{
		HeaderRowIndex: 1,
		Context:        &ctx,
	})

	dataFrame, err := reader.ScanData(dataStartAtRow)
	if err != nil {
		logger.Errorf("JobMigrateImportHistoryImpl | Err when scan data: %v", err)
		return err
	}

	logger.Info("JobMigrateImportHistoryImpl | Scan data done")

	if !reader.IsValid(dataFrame) {
		_, _, err := reader.WriteResult(dataFrame, resultHeader)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when get write result %v", err)
			return err
		}
		fileName := fmt.Sprintf("migrate_import_history_%s.xlsx", uuid.New().String())
		filePath := fmt.Sprintf("./%s", fileName)
		err = reader.WriteFile(filePath)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when write file %v", err)
			return err
		}
		defer func() {
			err = os.Remove(filePath)
			if err != nil {
				logger.Errorf("JobMigrateImportHistoryImpl | Err when remove file %v", err)
			}
		}()

		fileResult, err := os.Open(filePath)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when open file %v", err)
			return err
		}
		defer func() {
			err = fileResult.Close()
			if err != nil {
				logger.Errorf("JobMigrateImportHistoryImpl | Err when close file %v", err)
			}
		}()

		byteData := &bytes.Buffer{}
		_, err = byteData.ReadFrom(fileResult)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when read from file %v", err)
			return err
		}

		// Upload file service
		url, err := s.FileService.UploadFileWithBytesData(byteData, utils.XlsxContentType, fileName)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when upload file %v", err)
			return err
		}

		logger.Errorf("JobMigrateImportHistoryImpl | Not valid excel file %v", url)
		return errors.New(constant.ExcelMsgInvalidFile)
	}

	listProcessingFile := s.convertDataFrame2ProcessingFile(dataFrame)
	logger.Infof("JobMigrateImportHistoryImpl | Validate data frame done, total record: %d", len(listProcessingFile))

	if len(s.ForceDeleteClients) > 0 {
		logger.Infof("JobMigrateImportHistoryImpl | Delete processing file with client ids: %v", s.ForceDeleteClients)
		err = s.FpService.Delete(ctx, s.ForceDeleteClients)
		if err != nil {
			logger.Errorf("JobMigrateImportHistoryImpl | Err when delete: %v", err)
			return err
		}
	}

	logger.Infof("JobMigrateImportHistoryImpl | Bulk insert processing file, total record: %d", len(listProcessingFile))
	err = s.FpService.BulkInsertProcessingFile(ctx, listProcessingFile)
	if err != nil {
		logger.Errorf("JobMigrateImportHistoryImpl | Err when bulk insert processing file: %v", err)
		return err
	}

	return nil
}

func (s *JobMigrateImportHistoryImpl) convertDataFrame2ProcessingFile(dataFrame []*ImportHistoryDataObject) []fileprocessing.ProcessingFile {
	listProcessingFile := make([]fileprocessing.ProcessingFile, len(dataFrame))
	for idx, row := range dataFrame {
		listProcessingFile[idx] = fileprocessing.ProcessingFile{
			ProcessingFile: ent.ProcessingFile{
				ClientID:          row.ClientId,
				DisplayName:       row.DisplayName,
				ExtFileRequest:    row.ExtFileRequest,
				FileURL:           row.FileUrl,
				ResultFileURL:     row.ResultFileUrl,
				Status:            int16(row.Status),
				StatsTotalRow:     row.StatsTotalRow,
				StatsTotalSuccess: row.StatsTotalSuccess,
				ErrorDisplay:      row.ErrorDisplay,
				CreatedBy:         row.CreatedBy,
				CreatedAt:         time.Unix(row.CreatedAt, 0),
				UpdatedAt:         time.Unix(row.UpdatedAt, 0),
				MerchantID:        row.MerchantId,
				TenantID:          row.TenantId,
			},
		}
	}
	return listProcessingFile
}
