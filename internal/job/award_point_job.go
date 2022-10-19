package job

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/membertxn"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/loyalty"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/robfig/cron/v3"
)

var minPoint = 1

var fileAwardPointMetadata = dto.FileAwardPointMetadata{
	Phone: dto.CellData[string]{
		ColumnName: "Phone number (*)",
		Constrains: dto.Constrains{IsRequired: true, Regexp: constant.PhoneNumberRegex},
	},
	Point: dto.CellData[int]{
		ColumnName: "Points (*)",
		Constrains: dto.Constrains{IsRequired: true, Min: &minPoint},
	},
	Note: dto.CellData[string]{
		ColumnName: "Note",
		Constrains: dto.Constrains{IsRequired: false},
	},
}

func (f FileProcessingJob) StartGrantPointJob() bool {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	id, err := c.AddFunc(f.cfg.AwardPointJobConfig.Schedule, func() {
		logger.Infof("========== Running job GrantPointForMember Start  ...")
		f.grantPointForAllMemberTxn(context.Background())
		logger.Infof("========== Running job GrantPointForMember Finish ...")
	})
	if err != nil {
		logger.Errorf("Init Job failed: %v", err)
	}

	logger.Infof("Init Job success: ID = %v", id)

	c.Start()

	return false
}

func (f FileProcessingJob) grantPointForAllMemberTxn(ctx context.Context) {
	// Get list file award point
	fileAwardPoints, err := f.fapService.GetListFileAwardPointByStatuses(ctx, []int16{fileawardpoint.StatusInit, fileawardpoint.StatusProcessing})
	if err != nil {
		logger.Errorf("Cannot get list file award point, got: %v", err)
		return
	}

	// Check empty
	if len(fileAwardPoints) == 0 {
		logger.Infof("No init or processing file for executing!")
	}

	// Execute each file award point
	for _, fap := range fileAwardPoints {
		logger.Infof("Processing file award point ID = %v", fap.ID)

		var validTxnRecords []membertxn.MemberTxnDTO
		newFileResultUrl := fap.ResultFileURL
		switch fap.Status {
		case fileawardpoint.StatusInit:
			{
				validTxnRecords, newFileResultUrl, err = f.handleCaseInitAwardPoint(ctx, fap)
				if err != nil {
					logger.Errorf("Cannot handle init file award point %#v, got %v", fap, err)
					continue
				}
			}
		case fileawardpoint.StatusProcessing:
			{
				validTxnRecords, err = f.handleCaseProcessingAwardPoint(ctx, fap)
				if err != nil {
					logger.Errorf("Cannot handle processing file award point %#v, got %v", fap, err)
					continue
				}
			}
		}

		// Grant point for all member
		logger.Infof("Granting point for %d phone number", len(validTxnRecords))
		var errorResults []dto.FileAwardPointResultRow
		for _, record := range validTxnRecords {
			// Grand point for each member txn
			errorDisplay, err := f.grantPointForEachMemberTxn(ctx, fap, record)
			if err != nil {
				logger.Errorf("Grand point for member transaction %#v fail, got %v", record, err)
			}
			if len(errorDisplay) > 0 {
				errorResult := dto.FileAwardPointResultRow{
					Phone: record.Phone, Point: int(record.Point), Note: record.TxnDesc, Error: errorDisplay,
				}
				errorResults = append(errorResults, errorResult)
			}
		}

		// Save error grant point to result file
		if len(errorResults) > 0 {
			resultFileUrl, err := f.fileService.AppendErrorAndUploadFileAwardPointResult(errorResults, newFileResultUrl)
			if err != nil {
				logger.Errorf("===== Upload result filed failed")
			} else {
				logger.Infof("===== Update Result File URL to %v", resultFileUrl)
				_, err = f.fapService.UpdateResultFileUrlOne(ctx, fap.ID, resultFileUrl)
				if err != nil {
					logger.Errorf("===== Update Result File URL failed")
				}
			}
		}
	}
}

// handleCaseInitAwardPoint handle logic for file with init status
// 1. Download and extract data
// 2. Update total row
// 3. Insert member transaction to db
// 4. Update file status to processing or failed if errorRow == totalRow
// 5. Upload validation result file
func (f FileProcessingJob) handleCaseInitAwardPoint(ctx context.Context, fap *fileawardpoint.FileAwardPoint) ([]membertxn.MemberTxnDTO, string, error) {
	// 1. Download and extract data
	sheetData, err := excel.LoadExcelByUrl(fap.FileURL)
	if err != nil {
		logger.Errorf("Cannot get data from file url %v, got %v", fap.FileURL, err)
		return nil, "", err
	}

	indexStart := 3

	sheet, err := excel.ConvertToStruct[
		dto.FileAwardPointMetadata,
		dto.FileAwardPointRow,
		dto.Converter[dto.FileAwardPointMetadata, dto.FileAwardPointRow],
	](indexStart, &fileAwardPointMetadata, sheetData)

	if err != nil {
		logger.Errorf("Cannot convert data from file url %v, got %v", fap.FileURL, err)
		return nil, "", err
	}

	var (
		totalValidRow   = len(sheet.Data)
		totalInvalidRow = len(sheet.ErrorRows)
	)

	logger.Infof("Total valid row: %d", totalValidRow)
	logger.Infof("Total invalid row: %d", totalInvalidRow)

	// 2. Update total row
	_, err = f.fapService.UpdateTotalRowOne(ctx, fap.ID, totalValidRow+totalInvalidRow)
	if err != nil {
		logger.Errorf("Cannot update file award point total row: got: %v", err)
		return nil, "", err
	}

	var validTxnRecords []membertxn.MemberTxnDTO
	// 3. Insert member transaction to db
	for rowIndex, rowData := range sheet.Data {
		logger.Infof("Processing row index: %v data: %+v", rowIndex, rowData)
		createdMemberTxn, err := f.memTxnService.Create(context.Background(), membertxn.MemberTxnDTO{
			FileAwardPointID: int64(fap.ID),
			Point:            int64(rowData.Point),
			Phone:            rowData.Phone,
			TxnDesc:          rowData.Note,
			Status:           membertxn.StatusInit,
			SentTime:         time.Now(),
		})

		if err != nil {
			logger.Errorf("Cannot create member transaction %+v, got %v", rowData, err)
		}

		validTxnRecords = append(validTxnRecords, membertxn.MemberTxnDTO{
			ID:               int64(createdMemberTxn.ID),
			FileAwardPointID: int64(fap.ID),
			Point:            int64(rowData.Point),
			Phone:            rowData.Phone,
			Status:           membertxn.StatusInit,
			Error:            rowData.Error,
		})
	}

	// 4. Change file status to processing or failed if errorRow == totalRow
	if len(sheet.Data) == 0 {
		_, err = f.fapService.UpdateStatusOne(ctx, fap.ID, fileawardpoint.StatusFailed)
	} else {
		_, err = f.fapService.UpdateStatusOne(ctx, fap.ID, fileawardpoint.StatusProcessing)
	}

	if err != nil {
		logger.Errorf("Cannot update file award point status: got: %v", err)
		return nil, "", err
	}

	// 5. Upload validation result file
	fileName := utils.ExtractFileName(fap.DisplayName)
	resultFileName := fileName.FullNameWithSuffix("_result")
	newFileResultUrl, err := f.fileService.UploadFileAwardPointError(sheet.ErrorRows, resultFileName)
	if err != nil {
		logger.Errorf("Cannot upload result file, got %v", err)
		return nil, "", err
	}

	_, err = f.fapService.UpdateResultFileUrlOne(ctx, fap.ID, newFileResultUrl)
	if err != nil {
		logger.Errorf("Cannot save file award point result URL, got: %v", err)
		return nil, newFileResultUrl, err
	}
	return validTxnRecords, newFileResultUrl, nil
}

// handleCaseProcessingAwardPoint handle logic for file with processing status
// 1. Read member transaction record with status init from db
func (f FileProcessingJob) handleCaseProcessingAwardPoint(ctx context.Context, fap *fileawardpoint.FileAwardPoint) ([]membertxn.MemberTxnDTO, error) {
	memberTxnRecord, err := f.memTxnService.GetByFileAwardPointIDStatuses(ctx, int32(fap.ID), []int16{membertxn.StatusInit})
	if err != nil {
		logger.Errorf("GetByFileAwardPointIDStatuses got err: %v", err)
		return nil, err
	}
	return dto.MapMemberTxnToFileAwardPointRow(memberTxnRecord), nil
}

// grantPointForEachMemberTxn Call Loyalty core for granting point
// 1. Generate refID
// 2. Call API grant point in loyalty core
// 3. Update status member txn record
func (f FileProcessingJob) grantPointForEachMemberTxn(ctx context.Context, fap *fileawardpoint.FileAwardPoint, record membertxn.MemberTxnDTO) (string, error) {
	logger.Infof("Granting point for phone number %s", record.Phone)
	// 1. Generate refID
	refID := strings.ToUpper(utils.RandStringBytes(12))
	refID = fmt.Sprintf("%v%v", record.ID, refID)

	// 2. Call API grant point in loyalty core
	res, err := f.loyalClient.GrantPoint(loyalty.GrantPointRequest{
		MerchantID: strconv.FormatInt(fap.MerchantID, 10),
		RefId:      refID,
		Phone:      record.Phone,
		Point:      int32(record.Point),
		TxnDesc:    record.TxnDesc,
	})

	if err != nil {
		logger.Errorf("Call API grant point for number %v with info failed, %v got err: %v", record.Phone, record, err)
		return "Nạp điểm lỗi. Thử lại sau!", err
	}

	// 3. Update status member txn record
	if res.IsSuccess() {
		logger.Infof("Grant Point success for %v", record.Phone)
		loyaltyTxnId, _ := strconv.Atoi(res.Data.Transaction.TxnID)
		_, err = f.memTxnService.UpdateOne(ctx, membertxn.UpdateMemberTxnDTO{
			ID:           record.ID,
			RefID:        refID,
			SentTime:     time.Now().Truncate(time.Second),
			Status:       membertxn.StatusProcessing,
			Error:        res.Message,
			LoyaltyTxnID: int64(loyaltyTxnId),
		})
		return "", err
	} else {
		logger.Infof("Grant Point failed for %v: %v", record.Phone, res)
		_, err = f.memTxnService.UpdateOne(ctx, membertxn.UpdateMemberTxnDTO{
			ID:       record.ID,
			RefID:    refID,
			SentTime: time.Now(),
			Status:   membertxn.StatusFailed,
			Error:    res.Message,
		})

		// build error display -> save to result file
		errorDisplay := ""
		if res.IsNotEnoughBalance() {
			errorDisplay = "Nạp điểm lỗi do hệ thống thiếu điểm"
		} else {
			errorDisplay = fmt.Sprintf("Nạp điểm lỗi: %v", res.Message)
		}
		return errorDisplay, err
	}
}
