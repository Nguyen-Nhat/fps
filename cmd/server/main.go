package main

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job"
	"log"
	"os"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
	"git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/urfave/cli/v2"
)

func main() {
	cfg := config.Load()

	// Init logger
	_, err := logger.InitLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Cannot init logger, got err: %v", err)
	}

	// Init App for providing API
	app := &cli.App{
		Name:  "Loyalty File Processing Server",
		Usage: "...",
		Action: func(*cli.Context) error {
			// Init Job
			job.InitJob(cfg)

			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}
			return srv.Serve(cfg.Server.HTTP)
		},
		Commands: []*cli.Command{
			{
				Name:  job.Name,
				Usage: "File Processing Jobs",
				Subcommands: []*cli.Command{
					{
						Name:  job.ExecuteFile,
						Usage: "File Processing Jobs - Execute File",
						Action: func(ctx *cli.Context) error {
							job.InitJobExecuteFileThenRun(cfg)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
