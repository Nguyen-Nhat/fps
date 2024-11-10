package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	goSet "github.com/scylladb/go-set"
	"github.com/urfave/cli/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server"
	"git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/consumer"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/job"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/database/migrate"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/faltservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

var cfg config.Config

func main() {
	cfg = config.Load()
	config.Cfg = cfg

	// Init logger
	_, err := logger.InitLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Cannot init logger, got err: %v", err)
	}

	// Init i18n
	_, errI18n := i18n.LoadI18n(cfg.MessageFolder) // load i18n
	if errI18n != nil {
		log.Fatal(errI18n)
	}
	logger.Infof("Load i18n success! Say hello ...")
	logger.Infof("\t - [vi] %s", i18n.GetMessage("vi", "hello"))
	logger.Infof("\t - [vi] %s", i18n.GetMessageD("vi", "helloSomeone", "name", "FPS"))
	logger.Infof("\t - [en] %s", i18n.GetMessage("en", "hello"))
	logger.Infof("\t - [en] %s\n", i18n.GetMessageD("en", "helloSomeone", "name", "FPS"))

	// Init f-alt-service
	go func() {
		faltservice.InitParse(cfg)
	}()

	// Init App for providing API
	app := &cli.App{
		Name:  "Loyalty File Processing Server",
		Usage: "...",
		Commands: []*cli.Command{
			// all job config in here
			job.Command(cfg),
			// start server
			{
				Name:  "start",
				Usage: "start server",
				Action: func(*cli.Context) error {
					srv, err := server.NewServer(cfg)
					if err != nil {
						return err
					}
					return srv.Serve(cfg.Server.HTTP)
				},
			},
			// migrate database
			{
				Name:        "migrate",
				Usage:       "doing database migration",
				Subcommands: migrate.CliCommand(cfg.MigrationFolder, cfg.Database.MySQL.DatabaseTcpURI()),
			},
			// start kafka consumer
			{
				Name:   "start-kafka-consumer",
				Usage:  "Start Kafka Consumer",
				Action: startKafkaConsumer,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func startKafkaConsumer(cliCtx *cli.Context) error {
	cliConfig := consumer.CliConfig{
		WithRetry: false,
	}

	setArgs := goSet.NewStringSetWithSize(cliCtx.Args().Len())
	for i := 0; i < cliCtx.Args().Len(); i++ {
		setArgs.Add(cliCtx.Args().Get(i))
	}

	if setArgs.Has(constant.KafkaConsumerWithRetry) {
		cliConfig.WithRetry = true
	}

	consumerTypeCount := 0

	if setArgs.Has(constant.KafkaConsumeTypeForUpdateResultAsync) {
		cliConfig.ConsumerType = constant.KafkaConsumeTypeForUpdateResultAsync
		consumerTypeCount++
	}

	if consumerTypeCount != 1 {
		log.Fatal(fmt.Sprintf("argument must only be one of the following: %s", strings.Join([]string{
			constant.KafkaConsumeTypeForUpdateResultAsync,
		}, constant.SplitByCommaAndSpace)))
		return nil
	}

	consumerSrv := consumer.NewConsumer(cfg)
	newCtx, cancel := context.WithCancel(cliCtx.Context)
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	errCh := make(chan error)
	go func() {
		errCh <- consumerSrv.Consume(newCtx, cliConfig)
	}()
	select {
	case <-sigs:
		logger.Info("startKafkaConsumer context cancel")
	case err := <-errCh:
		if err != nil {
			logger.Errorf("startKafkaConsumer error: %v", err)
			return err
		}
	}
	return nil
}
