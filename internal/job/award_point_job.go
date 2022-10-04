package job

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/robfig/cron/v3"
)

func (a FileProcessingJob) grantPointForEachMemberInFileAwardPoint() bool {
	c := cron.New()

	id, err := c.AddFunc(a.cfg.AwardPointJobConfig.Schedule, func() {
		logger.Infof("Running job ... Start  ...")
		a.doSth()
		logger.Infof("Running job ... Finish ...")
	})
	if err != nil {
		logger.Errorf("Init Job failed: %v", err)
	}

	logger.Infof("Init Job success: ID = %v", id)

	c.Start()

	return false
}

func (a FileProcessingJob) doSth() {
	// Get Transaction
	txnArr, err := a.memTxnService.GetByFileAwardPointId(context.Background(), 1)
	if err != nil {
		return
	}

	// Check empty
	if len(txnArr) == 0 {
		logger.Infof("No Transaction for executing!")
	}

	// Execute each transaction
	for _, txn := range txnArr {
		logger.Infof("Processing Transaction ID = %v, data = %v", txn.ID, txn)
	}
}
