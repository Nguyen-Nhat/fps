package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
	"git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/database/migrate"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/faltservice"
)

func main() {
	cfg := config.Load()
	config.Cfg = cfg

	// Init logger
	_, err := logger.InitLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Cannot init logger, got err: %v", err)
	}

	// Init f-alt-service
	faltservice.InitParse(cfg)

	// Init App for providing API
	app := &cli.App{
		Name:  "Loyalty File Processing Server",
		Usage: "...",
		Action: func(*cli.Context) error {
			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}
			return srv.Serve(cfg.Server.HTTP)
		},
		Commands: []*cli.Command{
			job.Command(cfg),
			{
				Name:        "migrate",
				Usage:       "doing database migration",
				Subcommands: migrate.CliCommand(cfg.MigrationFolder, cfg.Database.MySQL.DatabaseTcpURI()),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
