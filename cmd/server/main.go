package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
	"git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
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
			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}
			return srv.Serve(cfg.Server.HTTP)
		},
		Commands: []*cli.Command{
			job.Command(cfg),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
