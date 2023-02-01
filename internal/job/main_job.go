package job

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/handlefileprocessing"
	"sync"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"github.com/xo/dburl"
)

var once sync.Once

type MainJob struct {
	cfg         config.JobConfig
	fileService fileservice.IService
}

var mainJob *MainJob

func InitJob(cfg config.Config) {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		return
	}
	logger.Infof("Connected to db %v", cfg.Database.MySQL.DBName)

	// file service
	fileServiceClient := fileservice.NewClient(cfg.ProviderConfig.FileService)
	fileService := fileservice.NewService(fileServiceClient)

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
				fileService: fileService,
			}
		})
	}

	// New job method
	handlefileprocessing.InitJobHandleProcessingFileAll(cfg.JobConfig.FileProcessingConfig, fpService, fprService, fileService)
}
