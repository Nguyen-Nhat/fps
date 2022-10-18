package job

import (
	"context"
	"fmt"
	"strconv"
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

func (a FileProcessingJob) StartGrantPointJob() bool {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	id, err := c.AddFunc(a.cfg.AwardPointJobConfig.Schedule, func() {
		logger.Infof("Running job GrantPointForMember Start  ...")
		a.grantPointForAllMemberTxn(context.Background())
		logger.Infof("Running job GrantPointForMember Finish ...")
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
		switch fap.Status {
		case fileawardpoint.StatusInit:
			{
				validTxnRecords, err = f.handleCaseInitAwardPoint(ctx, fap)
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
		for _, record := range validTxnRecords {
			// Grand point for each member txn
			err = f.grantPointForEachMemberTxn(ctx, fap, record)
			if err != nil {
				logger.Errorf("Grand point for member transaction %#v fail, got %v", record, err)
				continue
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
func (f FileProcessingJob) handleCaseInitAwardPoint(ctx context.Context, fap *fileawardpoint.FileAwardPoint) ([]membertxn.MemberTxnDTO, error) {
	// 1. Download and extract data
	sheetData, err := excel.LoadExcelByUrl(fap.FileURL)
	if err != nil {
		logger.Errorf("Cannot get data from file url %v, got %v", fap.FileURL, err)
		return nil, err
	}

	indexStart := 3

	sheet, err := excel.ConvertToStruct[
		dto.FileAwardPointMetadata,
		dto.FileAwardPointRow,
		dto.Converter[dto.FileAwardPointMetadata, dto.FileAwardPointRow],
	](indexStart, &fileAwardPointMetadata, sheetData)

	if err != nil {
		logger.Errorf("Cannot convert data from file url %v, got %v", fap.FileURL, err)
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	// 5. Upload validation result file
	fileName := utils.ExtractFileName(fap.DisplayName)
	newFileResultUrl, err := f.fileService.UploadFileAwardPointError(sheet.ErrorRows,
		fmt.Sprintf("%s_result.%s", fileName.Name, fileName.Extension))
	if err != nil {
		logger.Errorf("Cannot upload result file, got %v", err)
		return nil, err
	}

	_, err = f.fapService.UpdateResultFileUrlOne(ctx, fap.ID, newFileResultUrl)
	if err != nil {
		logger.Errorf("Cannot save file award point result URL, got: %v", err)
		return nil, err
	}
	return validTxnRecords, nil
}

// handleCaseProcessingAwardPoint handle logic for file with processing status
// 1. Read member transaction record with status init from db
func (f FileProcessingJob) handleCaseProcessingAwardPoint(ctx context.Context, fap *fileawardpoint.FileAwardPoint) ([]membertxn.MemberTxnDTO, error) {
	memberTxnRecord, err := f.memTxnService.GetByFileAwardPointIDStatuses(context.Background(), int32(fap.ID), []int16{membertxn.StatusInit})
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
func (f FileProcessingJob) grantPointForEachMemberTxn(ctx context.Context, fap *fileawardpoint.FileAwardPoint, record membertxn.MemberTxnDTO) error {
	logger.Infof("Granting point for phone number %s", record.Phone)
	// 1. Generate refID
	refID, err := utils.GenerateRandomString(16)
	if err != nil {
		logger.Errorf("Cannot generate random id, got %v", err)
		return err
	}

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
		return err
	}

	// 3. Update status member txn record
	if res.Code == constant.LoyaltyCoreCodeSuccess {
		_, err = f.memTxnService.UpdateOne(ctx, membertxn.UpdateMemberTxnDTO{
			ID:           int64(record.ID),
			RefID:        refID,
			SentTime:     time.Now().Truncate(time.Second),
			Status:       membertxn.StatusProcessing,
			Error:        res.Message,
			LoyaltyTxnID: res.Data.Transaction.TxnID,
		})
	} else {
		_, err = f.memTxnService.UpdateOne(ctx, membertxn.UpdateMemberTxnDTO{
			ID:       int64(record.ID),
			RefID:    refID,
			SentTime: time.Now(),
			Status:   membertxn.StatusFailed,
			Error:    res.Message,
		})
	}
	return err
}
