package job

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/handlefileprocessing"
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

type MainJob struct {
	fapService    fileawardpoint.Service
	memTxnService membertxn.Service
	cfg           config.JobConfig
	fileService   fileservice.IService
	loyalClient   loyalty.IClient
}

var mainJob *MainJob

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

	// New job
	// file processing
	fpRepo := fileprocessing.NewRepo(db)
	fpService := fileprocessing.NewService(fpRepo)

	// file processing row
	fprRepo := fileprocessingrow.NewRepo(db)
	fprService := fileprocessingrow.NewService(fprRepo)

	if mainJob == nil {
		once.Do(func() {
			mainJob = &MainJob{
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
	mainJob.StartGrantPointJob()
	initJobUpdateStatusFAP(cfg.JobConfig.UpdateStatusFAPJobConfig, fapService, memberTxnService, loyaltyClient, fileService)

	// New job method
	handlefileprocessing.InitJobHandleProcessingFileAll(cfg.JobConfig.FileProcessingConfig, fpService, fprService, fileService)
}
