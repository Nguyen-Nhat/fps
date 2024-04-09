package updatestatus

import (
	"context"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	configmapping2 "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
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
	stats, err := job.fprService.Statistics(file.ID)
	if err != nil {
		logger.ErrorT("Cannot statistics for file %v, err=%v", file.ID, err)
		return
	}
	status := file.Status
	resultFileUrl := file.ResultFileURL
	// 2. If all rows are processed
	if stats.IsFinished { // all finished
		logger.InfoT("File is finished executing!!!!")
		// 2.0. Status is Finished
		status = fileprocessing.StatusFinished

		// 2.1. Build and Upload result file to FileService
		isNeedToUploadResultFile := false
		for _, result := range stats.ErrorDisplays {
			if result != "" {
				isNeedToUploadResultFile = true
				break
			}
		}

		if isNeedToUploadResultFile {
			// 2.1.1. Inject Error Display to file
			outputFileType := constant.EmptyString
			switch cfgMapping.OutputFileType {
			case configmapping2.OutputFileTypeCSV:
				outputFileType = utils.CsvContentType
			case configmapping2.OutputFileTypeXLSX:
				outputFileType = utils.XlsxContentType
			default:
				logger.ErrorT("OutputFileType %v is not supported", cfgMapping.OutputFileType)
				return
			}
			fileDataBytes, err := excel.UpdateDataInColumn(file.FileURL, file.ExtFileRequest, cfgMapping.OutputFileType.String(), cfgMapping.DataAtSheet, cfgMapping.ErrorColumnIndex, int(cfgMapping.DataStartAtRow), stats.ErrorDisplays)
			if err != nil {
				logger.ErrorT("Update file with Error Display failed, err=%v", err)
				return
			}

			// 2.1.2. Gen result file name then Upload to file service
			resultFileName := utils.GetResultFileName(file.DisplayName)
			res, err := job.fileService.UploadFileWithBytesData(fileDataBytes, outputFileType, resultFileName)
			if err != nil {
				logger.ErrorT("Upload result file %v failed, err=%v", resultFileName, err)
				return
			}
			resultFileUrl = res
		}
	}

	// 3. Update processing_file: status, result_file_url, total_success
	// todo update processed_row
	_, err = job.fpService.UpdateStatusWithStatistics(ctx, file.ID, status, stats.TotalProcessed, stats.TotalSuccess, resultFileUrl)
	if err != nil {
		logger.ErrorT("Cannot update %v to failed, got error %v", fileprocessing.Name(), err)
		return
	}
}
