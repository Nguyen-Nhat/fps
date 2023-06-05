package job

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/executetask"
	"github.com/urfave/cli/v2"
	"github.com/xo/dburl"
	"os"
	"os/signal"
	"syscall"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/flatten"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/handlefileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

func Command(cfg config.Config) *cli.Command {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		panic(err)
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

	// services about config
	cmRepo := configmapping.NewRepo(db)
	cmService := configmapping.NewService(cmRepo)
	ctRepo := configtask.NewRepo(db)
	ctService := configtask.NewService(ctRepo)

	return &cli.Command{
		Name:  "jobs",
		Usage: "Loyalty File Processing Jobs",
		Subcommands: []*cli.Command{
			{
				Name:  "old",
				Usage: "old job",
				Action: func(*cli.Context) error {
					// Init Job
					handlefileprocessing.InitJobHandleProcessingFileAll(cfg.JobConfig.FileProcessingConfig,
						fpService, fprService, fileService)

					waitForKillingSign()
					return nil
				},
			},
			{
				Name:  "process-file",
				Usage: "Consumer handle consume message from kafka",
				Subcommands: []*cli.Command{
					{
						Name:  "flatten",
						Usage: "flatten data in file processing",
						Action: func(*cli.Context) error {
							job := flatten.NewJobFlattenManager(cfg.JobConfig.FlattenConfig,
								fpService, fprService, fileService,
								cmService, ctService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
					{
						Name:  "execute-task",
						Usage: "execute task for file processing",
						Action: func(*cli.Context) error {
							job := executetask.NewJobExecuteTaskManager(cfg.JobConfig.FlattenConfig,
								fpService, fprService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
				},
			},
		},
	}
}

func waitForKillingSign() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	endSignal := <-sig
	logger.Infof("Job end due to signal: %s", endSignal.String())
}
