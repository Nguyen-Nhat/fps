package consumer

import (
	"context"
	"time"

	"github.com/xo/dburl"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/consumer/topics"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
	"go.tekoapis.com/tekone/library/teka"
	"go.tekoapis.com/tekone/library/teka/events"
)

// Consumer handle event form Kafka
type Consumer struct {
	cfg config.Config

	fileProcessingRowRepository fileprocessingrow.Repo
}

type CliConfig struct {
	WithRetry    bool
	ConsumerType string
}

// NewConsumer Create new consumer
func NewConsumer(cfg config.Config) *Consumer {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		panic(err)
	}
	logger.Infof("Connected to db %v", cfg.Database.MySQL.DBName)

	return &Consumer{
		cfg:                         cfg,
		fileProcessingRowRepository: fileprocessingrow.NewRepo(db),
	}
}

func (c Consumer) Consume(ctx context.Context, cliConfig CliConfig) error {
	errCh := make(chan error)
	switch cliConfig.ConsumerType {
	case constant.KafkaConsumeTypeForUpdateResultAsync:
		go func() {
			topicName := c.cfg.Kafka.UpdateResultAsyncTopic
			consumerName := converter.TopicName2ConsumerName(c.cfg.Kafka.ConsumerPrefixName, topicName)
			subscriber := teka.MustNewSubscriber(
				events.Endpoint_FPS_CONSUMER,
				&events.Any{},
				teka.KAFKA,
				c.cfg.Kafka.ConnectionHost,
				teka.WithLogger(true),
				teka.WithPublisherName(topicName),
				teka.WithSubscriberName(consumerName),
				teka.WithRetry(false, nil),
				teka.WithDQL(false),
			)

			errCh <- subscriber.Consume(ctx, c.ProcessUpdateResultAsync)
		}()
	}
	return <-errCh
}

func (c Consumer) ProcessUpdateResultAsync(msg *teka.Message) (err error) {
	ctx := context.Background()
	logger.Infof("ProcessUpdateResultAsync start at %v", time.Now())
	defer func() {
		logger.Infof("ProcessUpdateResultAsync end at %v", time.Now())
		r := recover()
		if err != nil || r != nil {
			logger.Errorf("ProcessUpdateResultAsync error: %v, recover: %v", err, r)
			return
		}
	}()
	worker := topics.NewWorker(topics.WorkerAdjust{
		Cfg:                         c.cfg,
		FileProcessingRowRepository: c.fileProcessingRowRepository,
	})
	return worker.TopicUpsertResultAsync(ctx, msg)
}
