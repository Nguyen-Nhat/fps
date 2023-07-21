package faltservice

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/kylemcc/parse"
	"time"
)

type Session struct {
	ParseSession parse.Session
	ExpiredIn    time.Time
}

var InstanceSession *Session

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

// Authenticate to f-alt-server (Apply Singleton pattern)
func Authenticate() (*Session, error) {
	// If nil or expired (now + 3 days > timeExpiredIn)
	if InstanceSession == nil || InstanceSession.ExpiredIn.Before(time.Now().AddDate(0, 0, 3)) {
		login, err := parse.Login(
			config.Cfg.ProviderConfig.FAltService.Username,
			config.Cfg.ProviderConfig.FAltService.Password,
			nil)
		if err != nil {
			logger.Errorf("===== Authenticate to Parse failed, got: %+v", err.Error())
			return nil, err
		}

		InstanceSession = &Session{
			ParseSession: login,
			ExpiredIn:    time.Now().AddDate(0, config.Cfg.ProviderConfig.FAltService.SessionExpiredIn, 0),
		}
	}

	return InstanceSession, nil
}

// UpdateStatusProcessingFile Update status of model ProcessingFileParse
func UpdateStatusProcessingFile(fpsFileID int, status int16) error {
	if !config.Cfg.ProviderConfig.FAltService.IsEnable {
		logger.Info("===== Integration with f-alt disable")
		return nil
	}

	auth, err := Authenticate()
	if err != nil {
		logger.Errorf("===== Authenticate to parse failed: %+v", err.Error())
		return err
	}

	defer HandleParseErr(err)

	if fpsFileID != 0 {
		processingFileParse := &ProcessingFileParse{}

		q, err := auth.ParseSession.NewQuery(processingFileParse)
		if err != nil {
			logger.Errorf("===== NewQuery processingFileParse failed: %+v", err.Error())
			return err
		}

		err = q.EqualTo("fpsFileID", fpsFileID).First()
		if err != nil {
			logger.Errorf("===== Query processingFileParse by fpsFileID failed: %+v", err.Error())
			return err
		}

		if u, err := auth.ParseSession.NewUpdate(processingFileParse); err != nil {
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

	auth, err := Authenticate()
	if err != nil {
		logger.Errorf("===== Authenticate to parse failed: %+v", err.Error())
		return nil, err
	}

	defer HandleParseErr(err)

	err = auth.ParseSession.Create(processingFileParse)
	if err != nil {
		logger.Errorf("===== Save to Parse failed, got: %+v", err.Error())
		return processingFileParse, err
	}

	return processingFileParse, nil
}

// HandleParseErr Handle err of parse
func HandleParseErr(err error) {
	parseErr, _ := err.(parse.ParseError)

	// ParseServerCode = 209 - InvalidSessionToken
	if parseErr != nil && parseErr.Code() == constant.ParseInvalidSessionTokenCode {
		// The device’s session token is no longer valid.
		// The application should ask the service to log in to f-alt-server again.
		InstanceSession = nil
	}
}
