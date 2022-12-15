package main

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job"
	"log"
	"os"
	"os/signal"
	"syscall"

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
			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}
			return srv.Serve(cfg.Server.HTTP)
		},
		Commands: []*cli.Command{
			{
				Name:  "jobs",
				Usage: "Loyalty File Processing Jobs",
				Action: func(*cli.Context) error {
					// Init Job
					job.InitJob(cfg)

					sig := make(chan os.Signal, 1)
					signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
					endSignal := <-sig
					logger.Infof("Job end due to signal: %s", endSignal.String())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
