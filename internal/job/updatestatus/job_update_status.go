package updatestatus

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
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

	// 1. Statistics success row
	isFinished, totalSuccess, totalFailed, errorDisplays, err := job.fprService.Statistics(file.ID)
	if err != nil {
		logger.ErrorT("Cannot statistics for file %v, err=%v", file.ID, err)
		return
	}

	status := file.Status
	resultFileUrl := file.ResultFileURL
	// 2. If all rows are processed
	if isFinished { // all finished
		logger.InfoT("File is finished executing!!!!")
		// 2.0. Status is Finished
		status = fileprocessing.StatusFinished

		// 2.1. Build and Upload result file to FileService
		if totalFailed > 0 {
			cfgMapping, err := job.cfgMappingService.FindByClientID(ctx, file.ClientID)
			if err != nil {
				logger.ErrorT("Cannot find config mapping by clientID %v, err=%v", file.ClientID, err)
				return
			}

			// 2.1.1. Inject Error Display to file
			fileDataBytes, err := excel.UpdateDataInColumnOfFile(file.FileURL, "", cfgMapping.ErrorColumnIndex, int(cfgMapping.DataStartAtRow), errorDisplays, false)
			if err != nil {
				logger.ErrorT("Update file with Error Display failed, err=%v", err)
				return
			}

			// 2.1.2. Gen result file name then Upload to file service
			fileName := utils.ExtractFileName(file.FileURL)
			resultFileName := fileName.FullNameWithSuffix("_result")
			res, err := job.fileService.UploadFileWithBytesData(fileDataBytes, resultFileName)
			if err != nil {
				logger.ErrorT("Upload result file %v failed, err=%v", resultFileName, err)
				return
			}
			resultFileUrl = res
		}
	}

	// 3. Update processing_file: status, result_file_url, total_success
	// todo update processed_row
	_, err = job.fpService.UpdateStatusWithStatistics(ctx, file.ID, status, totalSuccess, resultFileUrl)
	if err != nil {
		logger.ErrorT("Cannot update %v to failed, got error %v", fileprocessing.Name(), err)
		return
	}
}
