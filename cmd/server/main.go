package main

import (
	"log"
	"os"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
