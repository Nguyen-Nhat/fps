package main

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
)

func main() {
	cfg := config.Load()
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
