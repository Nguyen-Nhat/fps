package job

import (
	"sync"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/membertxn"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/loyalty"
	"github.com/xo/dburl"
)

var once sync.Once

type FileProcessingJob struct {
	fapService    fileawardpoint.Service
	memTxnService membertxn.Service
	cfg           config.JobConfig
	fileService   fileservice.IService
	loyalClient   loyalty.IClient
}

var awardPointJob *FileProcessingJob

func InitJob(cfg config.Config) {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		return
	}
	logger.Infof("Connected to db %v", cfg.Database.MySQL.DBName)

	// file award point
	fapRepo := fileawardpoint.NewRepo(db)
	fapService := fileawardpoint.NewService(fapRepo)

	// member transaction
	mtRepo := membertxn.NewRepo(db)
	memberTxnService := membertxn.NewService(mtRepo)

	// file service
	fileServiceClient := fileservice.NewClient(cfg.ProviderConfig.FileService)
	fileService := fileservice.NewService(fileServiceClient)

	// loyaltyClient
	loyaltyClient := loyalty.NewClient(cfg.ProviderConfig.Loyalty)

	if awardPointJob == nil {
		once.Do(func() {
			awardPointJob = &FileProcessingJob{
				// config
				cfg: cfg.JobConfig,
				// services ...
				fapService:    fapService,
				memTxnService: memberTxnService,
				fileService:   fileService,
				loyalClient:   loyaltyClient,
			}
		})
	}

	// run job method
	awardPointJob.StartGrantPointJob()
	initJobUpdateStatusFAP(cfg.JobConfig.UpdateStatusFAPJobConfig, fapService, memberTxnService, loyaltyClient, fileService)
}
