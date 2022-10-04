package job

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/membertxn"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/xo/dburl"
	"sync"
)

var once sync.Once

type FileProcessingJob struct {
	fapService    fileawardpoint.Service
	memTxnService membertxn.Service
	cfg           config.JobConfig
}

var awardPointJob *FileProcessingJob

func InitJob(cfg config.Config) {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		return
	}
	logger.Info("Connected to db")

	// file award point
	fapRepo := fileawardpoint.NewRepo(db)
	fapService := fileawardpoint.NewService(fapRepo)

	// member transaction
	mtRepo := membertxn.NewRepo(db)
	memberTxnService := membertxn.NewService(mtRepo)

	if awardPointJob == nil {
		once.Do(func() {
			awardPointJob = &FileProcessingJob{
				// config
				cfg: cfg.JobConfig,
				// services ...
				fapService:    fapService,
				memTxnService: memberTxnService,
			}
		})
	}

	// run job method
	awardPointJob.grantPointForEachMemberInFileAwardPoint()
}
