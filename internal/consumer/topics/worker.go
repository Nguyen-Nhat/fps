package topics

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
)

type Worker struct {
	cfg                         config.Config
	fileProcessingRowRepository fileprocessingrow.Repo
}

type WorkerAdjust struct {
	Cfg                         config.Config
	FileProcessingRowRepository fileprocessingrow.Repo
}

func NewWorker(adjust WorkerAdjust) *Worker {
	return &Worker{
		cfg:                         adjust.Cfg,
		fileProcessingRowRepository: adjust.FileProcessingRowRepository,
	}
}
