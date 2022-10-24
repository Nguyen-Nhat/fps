package job

import (
	"context"
	"fmt"
	"sync"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/membertxn"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/loyalty"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/robfig/cron/v3"
)

var onceUSJ sync.Once

type (
	resultMessage string

	UpdateStatusFAPJob interface {
		InitCron()
		Start()
	}

	updateStatusFAPJobImpl struct {
		cronJob       *cron.Cron
		cfg           config.UpdateStatusFAPJobConfig
		fapService    fileawardpoint.Service
		memTxnService membertxn.Service
		loyaltyClient loyalty.IClient
		fileService   fileservice.IService
	}
)

const (
	rmNone            resultMessage = ""
	rmSuccess         resultMessage = "Thành công"
	rmExpiresChecking resultMessage = "Thời gian xử lý quá lâu. Vùi lòng kiểm tra lại kết quả giao dịch của KH"
	rmTxnFailed       resultMessage = "Giao dịch xử lý bị lỗi."
)

var updateStatusJob UpdateStatusFAPJob

func initJobUpdateStatusFAP(
	cfg config.UpdateStatusFAPJobConfig,
	fabService fileawardpoint.Service,
	memberTxnService membertxn.Service,
	loyaltyClient loyalty.IClient,
	fileService fileservice.IService,
) UpdateStatusFAPJob {
	if updateStatusJob == nil {
		onceUSJ.Do(func() {
			updateStatusJob = &updateStatusFAPJobImpl{
				cfg:           cfg,
				fapService:    fabService,
				memTxnService: memberTxnService,
				loyaltyClient: loyaltyClient,
				fileService:   fileService,
			}
			updateStatusJob.InitCron()
		})
	}
	return updateStatusJob
}

func (j *updateStatusFAPJobImpl) InitCron() {
	c := cron.New()
	jobName := "Job Check Status of File Award Point"

	// 2. Declare Job schedule
	id, err := c.AddFunc(j.cfg.Schedule, func() {
		logger.Infof("========== Running %v: Start  ...", jobName)
		j.Start()
		logger.Infof("========== Running %v: Finish ...", jobName)
	})
	if err != nil {
		logger.Errorf("Init %v Failed: %v", jobName, err)
	}

	logger.Infof("Init %v Success: ID = %v", jobName, id)
	c.Start()
}

// Start ...
//  1. Lấy những file với status=Processing (bảng file_award_point)
//  2. Check từng file:
//     2.1. Lấy những giao dịch của file đó có status=Processing (bảng member_transaction)
//     2.2. Check trạng thái từng giao dịch: Call Loyalty check trạng thái
//     2.2.1. Nếu call bị lỗi (timeout hoặc lỗi), chưa hết tgian check => bỏ qua, check giao dịch tiếp theo (GD này sẽ được check ở chu kỳ job tiếp theo)
//     2.2.2. Nếu call bị lỗi (timeout hoặc lỗi), hết tgian check      => cập nhật status=Timeout
//     2.2.3. Nếu call thành công, response code<>success, chưa hết tgian check => bỏ qua, check giao dịch tiếp theo (GD này sẽ được check ở chu kỳ job tiếp theo)
//     2.2.4. Nếu call thành công, response code<>success, hết tgian check      => cập nhật status=Timeout
//     2.2.5. Nếu call thành công, response code=success, status KHÔNG là THÀNH CÔNG  => cập nhật status=Failed
//     2.2.6. Nếu call thành công, response code=success, status là THÀNH CÔNG        => cập nhật status=Success
//     2.3. Gửi kết quả vào file result:
//     2.3.1. Build file result từ danh sách các giao dịch lỗi ở trên
//     2.3.2. Upload file lên File service
//     2.4. Nếu tất cả giao dịch của file đó đã kết thúc xử lý (status=Success/Failed) => cập nhật status của file thành Finished
//  3. Kết thúc
func (j *updateStatusFAPJobImpl) Start() {
	// 1. Get All FAPs with status=Processing
	ctx := context.Background()
	expiresTime := j.cfg.MaxCheckingTime
	faps, err := j.fapService.GetListFileAwardPointByStatuses(ctx, []int16{fileawardpoint.StatusProcessing})
	if err != nil {
		logger.Errorf("Failed to get %v with status = %v, error= %v", fileawardpoint.Name(), fileawardpoint.StatusProcessing, err)
		return // finish job
	}
	if len(faps) == 0 { // Check empty
		logger.Infof("No %v need to check status", fileawardpoint.Name())
	}

	// 2. Check for each FAP
	for _, fap := range faps {
		// 2.1. Fetch member transaction
		txnArr, err := j.memTxnService.GetTxnStillProcessing(ctx, fap.ID)
		if err != nil {
			logger.Errorf("Failed to get %v processing: %v", membertxn.Name(), err)
			continue // next FAP
		}
		if len(txnArr) == 0 { // check empty
			logger.Infof("No %v for executing!", membertxn.Name())
			//continue // no finished, need to handle at step 2.4
		}

		// 2.2. Execute each transaction
		var failedRows []dto.FileAwardPointResultRow
		for _, txn := range txnArr {
			logger.Infof("===== Processing Transaction ID = %v, data = %v", txn.ID, txn)
			isExpiresChecking := txn.IsCheckExpires(expiresTime)
			resMsg, err := j.updateStatusForMemberTransaction(ctx, txn, isExpiresChecking)
			if err != nil {
				logger.Errorf("Error when get txn from loyalty: %v", err)
			}

			// create result row
			switch resMsg {
			case rmTxnFailed, rmExpiresChecking:
				failedRows = append(failedRows, toFileAwardPointResultRow(txn, resMsg))
			case rmNone, rmSuccess:
				// do nothing
			}
		}

		// 2.3. Send result
		// 2.3.1. Build file
		// 2.3.2. Upload file
		if len(failedRows) == 0 {
			logger.Infof("===== No failed rows, no need to send result file!")
		} else {
			resultFileUrl, err := j.fileService.AppendErrorAndUploadFileAwardPointResult(failedRows, fap.ResultFileURL)
			if err != nil {
				logger.Errorf("===== Upload result filed failed, so we do not update Result File URL")
			} else {
				logger.Infof("===== Update Result File URL to %v", resultFileUrl)
				_, err = j.fapService.UpdateResultFileUrlOne(ctx, fap.ID, resultFileUrl)
				if err != nil {
					logger.Errorf("===== Update Result File URL failed")
				}
			}
		}

		// 2.4. Update status & statistics if finished
		totalTxn, totalSuccessTxn, totalHaveNotTerminated := j.memTxnService.GetStatistic(ctx, fap.ID)
		logger.Infof("===== Statistics for %v: totalTxn=%v, totalTxnSuccess=%v", fileawardpoint.Name(), totalTxn, totalSuccessTxn)
		// in case total success changed
		if int32(totalSuccessTxn) != fap.StatsTotalSuccess {
			_, err := j.fapService.UpdateStatsTotalSuccessOne(ctx, fap.ID, totalSuccessTxn)
			if err != nil {
				logger.Errorf("===== Update total success: %v", err)
			}
		}
		// in case finished
		// if total success = 0, set status to failed
		// else set status to finished
		if totalHaveNotTerminated == 0 {
			if totalSuccessTxn == 0 {
				logger.Infof("===== Update Status to Failed (%v)", fileawardpoint.StatusFailed)
				_, err := j.fapService.UpdateStatusOne(ctx, fap.ID, fileawardpoint.StatusFailed)
				if err != nil {
					logger.Errorf("===== Update Status failed: %v", err)
				}
				return
			}

			logger.Infof("===== Update Status to Finished (%v)", fileawardpoint.StatusFinished)
			_, err := j.fapService.UpdateStatusOne(ctx, fap.ID, fileawardpoint.StatusFinished)
			if err != nil {
				logger.Errorf("===== Update Status failed: %v", err)
			}
		}
	}
}

func (j *updateStatusFAPJobImpl) updateStatusForMemberTransaction(ctx context.Context, txn membertxn.MemberTransaction, isExpiresChecking bool) (resultMessage, error) {
	resMsg := rmNone
	res, err := j.loyaltyClient.GetTransactionByID(txn.LoyaltyTxnID)

	// 2.2.1. Error but not expires check => do nothing
	if err != nil && !isExpiresChecking {
		return rmNone, err
	}
	// 2.2.1. Error and expires check => update status=Timeout
	if err != nil && isExpiresChecking {
		resMsg = rmExpiresChecking
		j.memTxnService.UpdateStatusAndError(ctx, txn.ID, membertxn.StatusTimeout, string(resMsg))
	}

	// 2.2.3. Call Success, error code <> Success, not expires check => do nothing
	if res.IsFailed() && !isExpiresChecking {
		return rmNone, fmt.Errorf("call Loyalty failed, have not expires checking")
	}

	// 2.2.4. Call Success, error code <> Success, expires check => update status=Timeout
	if res.IsFailed() && isExpiresChecking {
		resMsg = rmExpiresChecking
		j.memTxnService.UpdateStatusAndError(ctx, txn.ID, membertxn.StatusTimeout, string(resMsg))
	}

	// 2.2.5 & 2.2.6: Call Success, error code = Success
	if res.IsSuccess() {
		// no transaction data => do nothing
		if res.Data == nil || len(res.Data.Transactions) == 0 {
			if !isExpiresChecking {
				return rmNone, fmt.Errorf("no transaction data in Loyalty response: %v", res)
			} else {
				resMsg = rmExpiresChecking
				j.memTxnService.UpdateStatusAndError(ctx, txn.ID, membertxn.StatusTimeout, string(resMsg))
			}
		} else { // get & check txn info
			loyaltyTxn := res.Data.Transactions[0]
			if loyaltyTxn.IsSuccess() { // 2.2.6. Call Success, error code = Success, status = Success
				resMsg = rmSuccess
				j.memTxnService.UpdateStatusAndError(ctx, txn.ID, membertxn.StatusSuccess, string(resMsg))
			} else { // 2.2.5. Call Success, error code = Success, status = Failed
				resMsg = rmTxnFailed
				j.memTxnService.UpdateStatusAndError(ctx, txn.ID, membertxn.StatusFailed, string(resMsg))
			}
		}
	}

	return resMsg, nil
}

func toFileAwardPointResultRow(txn membertxn.MemberTransaction, resMsg resultMessage) dto.FileAwardPointResultRow {
	return dto.FileAwardPointResultRow{
		Phone: txn.Phone,
		Point: int(txn.Point),
		Note:  txn.TxnDesc,
		Error: string(resMsg),
	}
}
