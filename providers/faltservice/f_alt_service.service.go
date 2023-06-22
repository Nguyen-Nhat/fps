package faltservice

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/kylemcc/parse"
)

// InitParse Init Parse (The core of f-alt-service)
func InitParse(cfg config.Config) {
	if !config.Cfg.ProviderConfig.FAltService.IsEnable {
		logger.Info("===== Integration with f-alt disable")
		return
	}

	parse.Initialize(
		cfg.ProviderConfig.FAltService.AppID,
		cfg.ProviderConfig.FAltService.RestKey,
		cfg.ProviderConfig.FAltService.MasterKey)

	parse.ServerURL(cfg.ProviderConfig.FAltService.Endpoint)

	if _, err := parse.ServerHealthCheck(); err != nil {
		logger.Errorf("===== Connected to f-alt-server failed: %+v", err.Error())
	}
}

// UpdateStatusProcessingFile Update status of model ProcessingFileParse
func UpdateStatusProcessingFile(fpsFileID int, status int16) error {
	if !config.Cfg.ProviderConfig.FAltService.IsEnable {
		logger.Info("===== Integration with f-alt disable")
		return nil
	}

	if fpsFileID != 0 {
		processingFileParse := &ProcessingFileParse{}

		q, err := parse.NewQuery(processingFileParse)
		if err != nil {
			logger.Errorf("===== NewQuery processingFileParse failed: %+v", err.Error())
			return err
		}

		err = q.EqualTo("fpsFileID", fpsFileID).First()
		if err != nil {
			logger.Errorf("===== Query processingFileParse by fpsFileID failed: %+v", err.Error())
			return err
		}

		if u, err := parse.NewUpdate(processingFileParse); err != nil {
			logger.Errorf("===== Update fileParseID=%+v failed: %+v", processingFileParse.Id, err.Error())
			return err
		} else {
			u = u.Set("status", status)
			err = u.Execute()
			if err != nil {
				logger.Errorf("===== Update fileParseID=%+v, status=%+v failed: %+v", processingFileParse.Id, status, err.Error())
				return err
			}
		}
	}

	return nil
}

// CreateProcessingFile Create record ProcessingFileParse
func CreateProcessingFile(processingFileParse *ProcessingFileParse) (*ProcessingFileParse, error) {
	if !config.Cfg.ProviderConfig.FAltService.IsEnable {
		logger.Info("===== Integration with f-alt disable")
		return nil, nil
	}

	err := parse.Create(processingFileParse, false)
	if err != nil {
		logger.Errorf("===== Save to Parse failed, got: %+v", err.Error())
		return processingFileParse, err
	}

	return processingFileParse, nil
}
