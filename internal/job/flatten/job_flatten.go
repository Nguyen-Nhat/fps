package flatten

import (
	"context"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
)

type jobFlatten struct {
	// services
	fpService         fileprocessing.Service
	fprService        fileprocessingrow.Service
	fpRowGroupService fpRowGroup.Service
	fileService       fileservice.IService
	// services config
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

func newJobFlatten(
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fpRowGroupService fpRowGroup.Service,
	fileService fileservice.IService,
	cfgMappingService configmapping.Service,
	cfgTaskService configtask.Service,
) *jobFlatten {
	return &jobFlatten{
		fpService:         fpService,
		fprService:        fprService,
		fpRowGroupService: fpRowGroupService,
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
//  5. 5. validate row group config with file data
//     -> if error 	==> terminate, remaining rows will be executed at next run cycle
//
//  6. Save row and row_group data into processing_file_row
//     -> if error 	==> terminate, remaining rows will be executed at next run cycle
//
//  7. Update processing_file: status=Processing, total_mapping, stats_total_row
func (job *jobFlatten) Flatten(ctx context.Context, file fileprocessing.ProcessingFile) {
	logger.Infof("----- Start flattening ProcessingFile with ID = %v \nFile = %+v", file.ID, file)

	// 1. Check status
	if file.Status != fileprocessing.StatusInit {
		logger.ErrorT("Not handle status=%v", file.Status)
		return
	}

	// 2. Load Config Mapping & Config Task
	cfgLoaderFactory := configloader.NewConfigLoaderFactory(job.cfgMappingService, job.cfgTaskService)
	configMapping, err := cfgLoaderFactory.GetConfigLoader(file)
	if err != nil {
		logger.ErrorT("Cannot load config mapping, fileID = %v", file.ID)
		job.updateFileProcessingToFailed(ctx, file, errConfigMapping, nil)
		return
	}

	// 3. Get file by URL -> data in first sheet
	sheetData, err := excel.LoadExcelByUrl(file.FileURL, configMapping.DataAtSheet)
	if err != nil {
		logger.ErrorT("Cannot get data from fileURL, fileID = %v, url = %v, got error %v", file.ID, file.FileURL, err)
		job.updateFileProcessingToFailed(ctx, file, errFileCannotLoad, nil)
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

	// 5. validate row group config with file data
	// 5.1. Group Row
	errorRows, createRowGroupJobs := validateAndBuildRowGroupData(file.ID, configMapping, configMappingsWithData)
	// 5.2. Check error, logic is same to step (4.2)
	if configMapping.IsSupportGrouping() && (len(errorRows) > 0 || len(createRowGroupJobs) == 0) {
		// Logging
		if len(errorRows) > 0 {
			logger.ErrorT("Importing file is invalid, fileID = %v, error in %v row(s)", file.ID, len(errorRows))
			logger.ErrorT("Error rows = \n%v\n", utils.JsonString(errorRows))
		} else if len(createRowGroupJobs) == 0 {
			logger.ErrorT("Importing file is invalid, fileID = %v, no group", file.ID)
		} else {
			logger.ErrorT("Importing file is invalid, fileID = %v, unknown error", file.ID)
		}
		// Update file result
		resultFileUrl := job.updateFileResult(configMapping, file.FileURL, errorRows)
		// Update file processing
		job.updateFileProcessingToFailed(ctx, file, errFileInvalid, &resultFileUrl)
		return
	}

	// 6. Save row and row_group data into processing_file_row
	if err = job.extractRowAndRowGroupToDB(ctx, file.ID, configMappingsWithData, createRowGroupJobs); err != nil {
		logger.ErrorT("Cannot save extracted data of fileID=%v, got err=%v", file.ID, err)
		return
	}

	// 7. Update processing_file: status=Processing, total_mapping, stats_total_row
	pf, err := job.fpService.UpdateToProcessingStatusWithExtractedData(ctx, file.ID, len(configMapping.Tasks), len(configMappingsWithData))
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
		rowID := errorRow.RowId
		if errMsg, existed := errorDisplays[rowID]; existed {
			errorDisplays[rowID] = fmt.Sprintf("%s; %s", errMsg, errorRow.Reason)
		} else {
			errorDisplays[rowID] = errorRow.Reason
		}
	}

	// 2. Inject error to importing file
	fileDataBytes, err := excel.UpdateDataInColumnOfFile(fileURL, cfgMapping.DataAtSheet, cfgMapping.ErrorColumnIndex, cfgMapping.DataStartAtRow, errorDisplays, false)
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

func (job *jobFlatten) extractRowAndRowGroupToDB(ctx context.Context,
	fileID int,
	configMappingMDs []*configloader.ConfigMappingMD,
	createRowGroupJobs []fpRowGroup.CreateRowGroupJob) error {

	// 1. Save Row Group
	if len(createRowGroupJobs) > 0 {
		err := job.fpRowGroupService.SaveExtractedRowGroupFromFile(ctx, fileID, createRowGroupJobs)
		if err != nil {
			return err
		}
	} else {
		logger.Infof("----- fileID=%d -> No row group need to save to DB", fileID)
	}

	// 2. Add extracted data to ProcessingFileRow
	var pfrCreateList []fileprocessingrow.CreateProcessingFileRowJob
	for _, mapping := range configMappingMDs {
		for _, task := range mapping.Tasks {
			pfr := fileprocessingrow.CreateProcessingFileRowJob{
				FileID:       fileID,
				RowIndex:     task.ImportRowIndex,
				RowDataRaw:   utils.JsonString(task.ImportRowData),
				TaskIndex:    task.TaskIndex,
				TaskMapping:  utils.JsonString(mapping),
				GroupByValue: mapping.RowGroupValue,
			}
			pfrCreateList = append(pfrCreateList, pfr)
		}
	}

	// 3. Save Row-Task
	return job.fprService.SaveExtractedRowTaskFromFile(ctx, fileID, pfrCreateList)
}
