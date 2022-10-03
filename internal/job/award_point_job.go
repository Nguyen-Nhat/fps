package job

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/robfig/cron/v3"
	"github.com/xo/dburl"
	"sync"
)

var once sync.Once

type AwardPointJob struct {
	fapService fileawardpoint.Service
	cfg        config.AwardPointJobConfig
}

var awardPointJob *AwardPointJob

func InitJob(cfg config.Config) {
	db, err := dburl.Open(cfg.Database.MySQL.DatabaseURI())
	if err != nil {
		logger.Errorf("Fail to open db, got: %v", err)
		return
	}
	logger.Info("Connected to db")
	fapRepo := fileawardpoint.NewRepo(db)
	fapService := fileawardpoint.NewService(fapRepo)

	if awardPointJob == nil {
		once.Do(func() {
			awardPointJob = &AwardPointJob{
				fapService: fapService,
				cfg:        cfg.JobConfig.AwardPointJobConfig,
			}
		})
	}
	awardPointJob.grantPointForEachMemberInFileAwardPoint()
}

func (a AwardPointJob) grantPointForEachMemberInFileAwardPoint() bool {
	c := cron.New()

	id, err := c.AddFunc(a.cfg.Schedule, func() {
		logger.Infof("Running job ...")
	})
	if err != nil {
		logger.Errorf("Init Job failed: %v", err)
	}

	logger.Infof("Init Job success: ID = %v", id)

	c.Start()

	return false
}
