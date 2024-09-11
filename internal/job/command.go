package job

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/xo/dburl"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/executerowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/executetask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/flatten"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/handlefileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/support_migrate_data"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job/updatestatus"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	googleSheetIdArg      = "google-sheet-id"
	forceDeleteClientsArg = "force-delete-clients"
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

	// file processing row group
	fpRowGroupRepo := fpRowGroup.NewRepo(db)
	fpRowGroupService := fpRowGroup.NewService(fpRowGroupRepo)

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
				Usage: "old job deprecated",
				Action: func(*cli.Context) error {
					// Init Job
					handlefileprocessing.InitJobHandleProcessingFileAll(cfg.JobConfig.FileProcessingConfig,
						fpService, fprService, fileService, cmService)

					waitForKillingSign()
					return nil
				},
			},
			{
				Name:  "process-file",
				Usage: "Handle process file uploaded, with many sub command",
				Subcommands: []*cli.Command{
					{
						Name:  "flatten",
						Usage: "Read file in processing_file table and flatten data in file upload to each task in row to processing_file_row or processing_file_row_group table",
						Action: func(*cli.Context) error {
							job := flatten.NewJobFlattenManager(cfg,
								fpService, fprService, fpRowGroupService, fileService,
								cmService, ctService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
					{
						Name:  "execute-task",
						Usage: "Read processing_file_row table and execute task and update status",
						Action: func(*cli.Context) error {
							job := executetask.NewJobExecuteTaskManager(cfg,
								fpService, fprService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
					{
						Name:  "execute-row-group",
						Usage: "If config task has config group will read processing_file_row_group table to execute group task for file processing and update status",
						Action: func(*cli.Context) error {
							job := executerowgroup.NewJobExecuteRowGroupManager(cfg,
								fpService, fprService, fpRowGroupService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
					{
						Name:  "update-status",
						Usage: "Checking status of each task in processing_file_row table and gen file result",
						Action: func(*cli.Context) error {
							job := updatestatus.NewJobUpdateStatusManager(cfg,
								fpService, fprService, fileService,
								cmService)
							job.Start()

							waitForKillingSign()
							return nil
						},
					},
				},
			},
			{
				Name:  "migrate-import-history",
				Usage: "Migrate Import History",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     googleSheetIdArg,
						Usage:    "Google Sheet ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     forceDeleteClientsArg,
						Usage:    "Force delete clients",
						Required: false,
					},
				},
				Action: func(cliCtx *cli.Context) error {
					sheetId := cliCtx.String(googleSheetIdArg)
					forceDeleteClientStr := cliCtx.String(forceDeleteClientsArg)
					forceDeleteClientInts := make([]int32, 0)
					if forceDeleteClientStr != constant.EmptyString {
						forceDeleteClientInts, err = utils.String2ArrayInt32(forceDeleteClientStr, constant.SplitByComma)
						if err != nil {
							logger.Errorf("Fail to convert force delete clients, got: %v", err)
							return err
						}
					}

					job, err := support_migrate_data.NewJobMigrateImportHistory(sheetId, forceDeleteClientInts, fpService, fileService)
					if err != nil {
						logger.Errorf("Fail to create job, got: %v", err)
						return err
					}

					err = job.Run(cliCtx.Context)
					if err != nil {
						logger.Errorf("Fail to run job, got: %v", err)
						return err
					}
					return nil
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
