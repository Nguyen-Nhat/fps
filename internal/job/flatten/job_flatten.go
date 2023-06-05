package flatten

import (
	"context"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

type jobFlatten struct {
	// services
	fpService   fileprocessing.Service
	fprService  fileprocessingrow.Service
	fileService fileservice.IService
	// services config
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

func newJobFlatten(
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
	cfgTaskService configtask.Service,
) *jobFlatten {
	return &jobFlatten{
		fpService:         fpService,
		fprService:        fprService,
		fileService:       fileService,
		cfgMappingService: cfgMappingService,
		cfgTaskService:    cfgTaskService,
	}
}

// Flatten ...
//
//  1. Check file.status = INIT
//
//  2. Get file by URL
//     -> if error 	==> update status = FAILED
//
//  3. Load Config Mapping & Config Task
//     -> if error 	==> update status = FAILED
//
//  4. Validate importing data
//     -> if error 	==> update status = FAILED
//
//  5. Save row data into processing_file_row
//     -> if error 	==> terminate, remaining rows will be executed at next run cycle
//
//  6. Update processing_file: status=Processing, total_mapping, stats_total_row
func (job *jobFlatten) Flatten(ctx context.Context, file fileprocessing.ProcessingFile) {
	logger.Infof("----- Start flattening ProcessingFile with ID = %v \nFile = %+v", file.ID, file)

	// 1. Check status
	if file.Status != fileprocessing.StatusInit {
		logger.ErrorT("Not handle status=%v", file.Status)
		return
	}

	// 2. Get file by URL -> data in first sheet
	sheetData, err := excel.LoadExcelByUrl(file.FileURL)
	if err != nil {
		logger.ErrorT("Cannot get data from fileURL, fileID = %v, url = %v, got error %v", file.ID, file.FileURL, err)
		job.updateFileProcessingToFailed(ctx, file, errFileCannotLoad, nil)
		return
	}

	// 3. Load Config Mapping & Config Task
	cfgLoaderFactory := configloader.NewConfigLoaderFactory(job.cfgMappingService, job.cfgTaskService)
	configMapping, err := cfgLoaderFactory.GetConfigLoader(file)
	if err != nil {
		logger.ErrorT("Cannot load config mapping, fileID = %v", file.ID)
		job.updateFileProcessingToFailed(ctx, file, errConfigMapping, nil)
		return
	}

	// 4. Validate importing data
	configMappingsWithData, errorRows, err := validateImportingData(sheetData, configMapping)
	if err != nil {
		logger.ErrorT("Importing file is invalid, fileID = %v, error = %+v", file.ID, err)
		job.updateFileProcessingToFailed(ctx, file, errFileInvalid, nil)
		return
	}
	if len(errorRows) > 0 {
		// Logging
		logger.ErrorT("Importing file is invalid, fileID = %v, error in %v row(s)", file.ID, len(errorRows))
		logger.ErrorT("Error rows = \n%v\n", utils.JsonString(errorRows))
		// Update file result
		resultFileUrl := job.updateFileResult(configMapping, file.FileURL, errorRows)
		// Update file processing
		job.updateFileProcessingToFailed(ctx, file, errFileInvalid, &resultFileUrl)
		return
	}
	logger.Infof("Config mapping = \n%s\n", utils.JsonString(configMappingsWithData))

	// 5. Save row data into processing_file_row
	if err := job.extractDataAndUpdateFileStatusInDB(ctx, file.ID, configMappingsWithData); err != nil {
		logger.ErrorT("Cannot save extracted data of fileId=%v, got err=%v", file.ID, err)
		return
	}

	// 6. Update processing_file: status=Processing, total_mapping, stats_total_row
	pf, err := job.fpService.UpdateToProcessingStatusWithExtractedData(ctx, file.ID, len(configMappingsWithData), len(configMappingsWithData))
	if err != nil {
		logger.ErrorT("Cannot update %v id=%v,got err=%v", fileprocessing.Name(), file.ID, err)
		return
	}
	logger.Infof("Update %v success, data in DB=%v", fileprocessing.Name(), pf)
}

// ---------------------------------------------------------------------------------------------------------------------

// updateFileProcessingToFailed ...
func (job *jobFlatten) updateFileProcessingToFailed(ctx context.Context, file fileprocessing.ProcessingFile, errMsg fileprocessing.ErrorDisplay, resultFileURL *string) {
	_, updateStatusErr := job.fpService.UpdateToFailedStatusWithErrorMessage(ctx, file.ID, errMsg, resultFileURL)
	if updateStatusErr != nil {
		logger.ErrorT("Cannot update %v to fail, got error %v", fileprocessing.Name(), updateStatusErr)
	}
}

// updateFileResult ...
func (job *jobFlatten) updateFileResult(cfgMapping configloader.ConfigMappingMD, fileURL string, errorRows []ErrorRow) string {
	// 1. Convert errorRows to errorDisplays
	errorDisplays := make(map[int]string)
	for _, errorRow := range errorRows {
		id := errorRow.RowId
		if errMsg, existed := errorDisplays[id]; existed {
			errorDisplays[id] = fmt.Sprintf("%s; %s", errMsg, errorRow.Reason)
		} else {
			errorDisplays[id] = errorRow.Reason
		}
	}

	// 2. Inject error to importing file
	fileDataBytes, err := excel.UpdateDataInColumnOfFile(fileURL, "", cfgMapping.ErrorColumnIndex, cfgMapping.DataStartAtRow, errorDisplays, false)
	if err != nil {
		logger.ErrorT("Update file with Error Display failed, err=%v", err)
		return ""
	}

	// 3. Gen result file name then Upload to file service
	fileName := utils.ExtractFileName(fileURL)
	resultFileName := fileName.FullNameWithSuffix("_result")
	resultFileUrl, err := job.fileService.UploadFileWithBytesData(fileDataBytes, resultFileName)
	if err != nil {
		logger.ErrorT("Upload result file %v failed, err=%v", resultFileName, err)
		return ""
	}

	return resultFileUrl
}

func (job *jobFlatten) extractDataAndUpdateFileStatusInDB(ctx context.Context, fileId int,
	configMappingMDs []configloader.ConfigMappingMD) error {
	// 1. Add extracted data to ProcessingFileRow
	var pfrCreateList []fileprocessingrow.CreateProcessingFileRowJob
	for _, mapping := range configMappingMDs {
		for _, task := range mapping.Tasks {
			pfr := fileprocessingrow.CreateProcessingFileRowJob{
				FileId:      fileId,
				RowIndex:    task.ImportRowIndex,
				RowDataRaw:  utils.JsonString(task.ImportRowData),
				TaskIndex:   task.TaskIndex,
				TaskMapping: utils.JsonString(mapping),
			}
			pfrCreateList = append(pfrCreateList, pfr)
		}
	}

	// 2. Save
	return job.fprService.SaveExtractedDataFromFile(ctx, fileId, pfrCreateList)
}
