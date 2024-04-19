package updatestatus

import (
	"context"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
)

type jobUpdateStatus struct {
	// services
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	fileService fileservice.IService
	// services config
	cfgMappingService configmapping.Service
}

func newJobUpdateStatus(
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
) *jobUpdateStatus {
	return &jobUpdateStatus{
		fpService:         fpService,
		fprService:        fprService,
		fileService:       fileService,
		cfgMappingService: cfgMappingService,
	}
}

// UpdateStatus ...
//  1. Statistics success row
//  2. If all rows are processed
//     2.1. Build and Upload result file to FileService
//  3. Update processing_file: status, result_file_url, total_success
func (job *jobUpdateStatus) UpdateStatus(ctx context.Context, file fileprocessing.ProcessingFile) {
	logger.Infof("----- Start Update Status ProcessingFile with ID = %v", file.ID)
	cfgMapping, err := job.cfgMappingService.FindByClientID(ctx, file.ClientID)
	if err != nil {
		logger.ErrorT("Cannot find config mapping by clientID %v, err=%v", file.ClientID, err)
		return
	}
	// Check time out
	if file.CreatedAt.Add(time.Second * time.Duration(cfgMapping.Timeout)).Before(time.Now()) {
		logger.ErrorT("File is timeout")
		err = job.fprService.ForceTimeout(ctx, file.ID)
		if err != nil {
			logger.ErrorT("Cannot force timeout for file %v, err=%v", file.ID, err)
			return
		}
	}

	// 1. Statistics success row
	stats, err := job.fprService.Statistics(ctx, file.ID)
	if err != nil {
		logger.ErrorT("Cannot statistics for file %v, err=%v", file.ID, err)
		return
	}
	status := file.Status
	var resultFileURL = file.ResultFileURL
	// 2. If all rows are processed
	if stats.IsFinished { // all finished
		logger.InfoT("File is finished executing!!!!")
		// 2.1. Status is Finished
		status = fileprocessing.StatusFinished

		// 2.2. Update result file
		resultFileURL, err = job.buildResultFileAndUpload(ctx, file, cfgMapping, stats.ErrorDisplays)
		if err != nil {
			logger.ErrorT("Update file with Error Display failed, err=%v", err)
			return
		}
	}

	// 3. Update processing_file: status, result_file_url, total_success
	// todo update processed_row
	_, err = job.fpService.UpdateStatusWithStatistics(ctx, file.ID, status, stats.TotalProcessed, stats.TotalSuccess, resultFileURL)
	if err != nil {
		logger.ErrorT("Cannot update %v to failed, got error %v", fileprocessing.Name(), err)
		return
	}
}
