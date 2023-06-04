package handlefileprocessing

import (
	"context"
	"encoding/json"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/fileservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/taskprovider"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
)

type (
	JobHandleProcessingFile interface {
		Execute(context.Context, *fileprocessing.ProcessingFile)
	}

	jobHandleProcessingFileImpl struct {
		fpService   fileprocessing.Service
		fprService  fileprocessingrow.Service
		fileService fileservice.IService
	}
)

const (
	errFileError fileprocessing.ErrorDisplay = "File tải lên sai template"
)

var _ JobHandleProcessingFile = &jobHandleProcessingFileImpl{}

func newJobHandleProcessingFile(
	fpService fileprocessing.Service,
	fprService fileprocessingrow.Service,
	fileService fileservice.IService,
) *jobHandleProcessingFileImpl {
	job := &jobHandleProcessingFileImpl{
		fpService:   fpService,
		fprService:  fprService,
		fileService: fileService,
	}
	return job
}

/**
1. Lấy file ở trạng thái Init, Processing
2. Xử lý file ở trạng thái INIT
3. Xử lý file ở trạng thái Processing
4. Thống kê số row đã thành công, cập nhật lại trạng thái file
*/

func (j *jobHandleProcessingFileImpl) Execute(ctx context.Context, file *fileprocessing.ProcessingFile) {
	// 1. Lấy file ở trạng thái Init, Processing
	logger.Infof("----- Start Executing ProcessingFile with ID = %v, status=%v", file.ID, file.Status)
	if file.Status != fileprocessing.StatusInit && file.Status != fileprocessing.StatusProcessing {
		logger.ErrorT("Not handle status=%v", file.Status)
		return
	}

	// 2. Xử lý file ở trạng thái INIT
	if file.IsInitStatus() {
		logger.Infof("------ In case fileId=%v is at Init status", file.ID)
		// Output:
		//	+ case unknown error occur -> file.status=INIT;
		// 	+ case execute file failed -> file.status=Failed;
		//	+ case handle file success -> file.status=Processing, data is extracted to DB, that is used in next step
		file = j.handleFileInInitStatus(ctx, file)
	}

	// 3. Xử lý file ở trạng thái Processing
	if file.IsProcessingStatus() {
		logger.Infof("------ In case fileId=%v is at Processing status", file.ID)
		j.handleFileInProcessingStatus(ctx, file)
	}

	// 4. Thống kê số row đã thành công, cập nhật lại trạng thái file
	j.statisticAndUpdateFileStatus(ctx, file)
}

// Step 2 ==============================================================================================================

// HandleFileInInitStatus ...
/*
 2. Xử lý file ở trạng thái INIT
    2.1. Tải file từ File Service
    2.1.1. Nếu file không đọc được   -> chuyển trạng thái processing_file.status=Failed
    2.2. Validate sheet Mapping
    2.2.1. Nếu Mapping sai định dạng -> chuyển trạng thái processing_file.status=Failed
    2.3. Validate sheet Data
    2.3.1. Nếu sheet không có data   -> chuyển trạng thái processing_file.status=Failed
    note: không validate định dạng của các row (do không thể)
    2.4. Lưu các row data vào bảng processing_file_row
    2.7. Cập nhật lại processing_file: status=Processing, total_mapping, stats_total_row
*/
func (j *jobHandleProcessingFileImpl) handleFileInInitStatus(ctx context.Context, file *fileprocessing.ProcessingFile) *fileprocessing.ProcessingFile {
	// 2.1. Tải file từ File Service
	// 2.1.1. Nếu file không đọc được   -> chuyển trạng thái processing_file.status=Failed
	sheetDataMap, err := excel.LoadSheetsInExcelByUrl(file.FileURL, []string{sheetImportDataName, sheetMappingName})
	if err != nil {
		logger.ErrorT("Cannot get data from file url %v, got error %v", file.FileURL, err)

		fpUpdated, updateStatusErr := j.fpService.UpdateToFailedStatusWithErrorMessage(ctx, file.ID, errFileError, nil)
		if updateStatusErr != nil {
			logger.ErrorT("Cannot update %v to fail, got error %v", fileprocessing.Name(), err)
			return file // return original file
		}
		return fpUpdated // return file after update status=Failed
	}

	// 2.2. Validate sheet Mapping
	// 2.2.1. Nếu Mapping sai định dạng -> chuyển trạng thái processing_file.status=Failed
	sheetMappingResult, err := excel.ConvertToStruct[
		dto.SheetMappingMetadata,
		dto.MappingRow,
		dto.Converter[dto.SheetMappingMetadata, dto.MappingRow],
	](2, &sheetMappingMetadata, sheetDataMap[sheetMappingName])

	// ---> update status to failed if empty or error
	if err != nil || len(sheetMappingResult.ErrorRows) > 0 {
		if err != nil {
			logger.ErrorT("Cannot convert sheet %v from file url %v, got error %v", sheetMappingName, file.FileURL, err)
		} else {
			var errMsg string
			for _, errorRow := range sheetMappingResult.ErrorRows {
				errMsg += fmt.Sprintf("row=%v, message=%v; ", errorRow.RowId, errorRow.Reason)
			}
			logger.ErrorT("Sheet %v from file url %v has error data, got error %v", sheetMappingName, file.FileURL, errMsg)
		}
		fpUpdated, updateStatusErr := j.fpService.UpdateToFailedStatusWithErrorMessage(ctx, file.ID, errFileError, nil)
		if updateStatusErr != nil {
			logger.ErrorT("Cannot update %v to fail, got error %v", fileprocessing.Name(), err)
			return file // return original file
		}
		return fpUpdated // return file after update status=Failed
	}
	sheetMapping := sheetMappingResult.Data
	logger.Infof("---- Mapping sheet data: %+v", sheetMapping)

	// 2.3. Validate sheet Data
	// 2.3.1. Nếu sheet không có data   -> chuyển trạng thái processing_file.status=Failed
	metadata := getListCellFromMapping(sheetMapping)
	sheetDataResult, err := excel.ConvertToStructByMap(dataIndexStartInDataSheet, metadata, sheetDataMap[sheetImportDataName])
	if err != nil {
		logger.ErrorT("Cannot convert sheet %v from file url %v, got error %v", sheetImportDataName, file.FileURL, err)
		fpUpdated, updateStatusErr := j.fpService.UpdateToFailedStatusWithErrorMessage(ctx, file.ID, errFileError, nil)
		if updateStatusErr != nil {
			logger.ErrorT("Cannot update %v to fail, got error %v", fileprocessing.Name(), err)
			return file // return original file
		}
		return fpUpdated // return file after update status=Failed
	}
	sheetData := sheetDataResult.Data
	logger.Infof("---- ImportData sheet data: %+v", sheetData)

	// 2.4. Lưu các row data đúng định dạng vào bảng processing_file_row
	logger.Infof("Prepare Extract Data to DB and Update %v to Processing", fileprocessing.Name())
	filedUpdated, err := j.extractDataAndUpdateFileStatusInDB(ctx, file.ID, sheetMapping, sheetData)
	if err != nil {
		logger.ErrorT("%v", err)
		return file // return original file
	}
	return filedUpdated // return file after update status=Processing
}

func (j *jobHandleProcessingFileImpl) extractDataAndUpdateFileStatusInDB(ctx context.Context, fileId int,
	sheetMapping []dto.MappingRow, sheetData []map[string]string) (*fileprocessing.ProcessingFile, error) {
	// 1. Add extracted data to ProcessingFileRow
	var pfrCreateList []fileprocessingrow.CreateProcessingFileRowJob
	for rowId, data := range sheetData {
		for _, mapping := range sheetMapping {
			pfr := fileprocessingrow.CreateProcessingFileRowJob{
				FileId:      fileId,
				RowIndex:    rowId,
				RowDataRaw:  toJsonStringNotCareError(data),
				TaskIndex:   mapping.TaskId,
				TaskMapping: toJsonStringNotCareError(mapping),
			}
			pfrCreateList = append(pfrCreateList, pfr)
		}
	}
	if err := j.fprService.SaveExtractedDataFromFile(ctx, fileId, pfrCreateList); err != nil {
		logger.ErrorT("Cannot save extracted data of fileId=%v, got err=%v", fileId, err)
		return nil, err
	}

	// 2. Update ProcessingFile
	pf, err := j.fpService.UpdateToProcessingStatusWithExtractedData(ctx, fileId, len(sheetMapping), len(sheetData))
	if err != nil {
		logger.ErrorT("Cannot update %v id=%v,got err=%v", fileprocessing.Name(), fileId, err)
		return nil, err
	}
	logger.Infof("Update %v success, data in DB=%v", fileprocessing.Name(), pf)
	return pf, nil
}

func toJsonStringNotCareError(input interface{}) string {
	mappingJsonStr, err := json.Marshal(input)
	if err != nil {
		logger.Errorf("Error when convert %v to json string, err %v", input, err)
		return "_no_data_"
	}
	return string(mappingJsonStr)
}

func getListCellFromMapping(sheetMapping []dto.MappingRow) []dto.CellData[string] {
	var metadata []dto.CellData[string]
	for _, mappingData := range sheetMapping {
		requestMap := mappingData.Request
		for _, requestMapping := range requestMap {
			if requestMapping.IsMappingExcel {
				metadata = append(metadata, dto.CellData[string]{ColumnName: requestMapping.MappingKey})
			}
		}
	}
	return metadata
}

// Step 3 ==============================================================================================================

/*
 3. Xử lý file ở trạng thái Processing
    3.1. Lấy toàn bộ các record (task) có processing_file_row.status=Init
    3.2. Group task theo theo row_index, sau đó xử lý từng task theo thứ tự trong Mapping
    3.2.1. Nếu input đầu vào là từ response của task trước và nó không xuất hiện
    -> chuyển trạng thái processing_file_row.status=Failed
    -> chuyển sang group task tiếp theo
    3.2.2. Nếu call api lỗi do timeout hoặc api trả về lỗi
    -> chuyển trạng thái processing_file_row.status=Failed
    -> chuyển sang group task tiếp theo
    3.2.3. Nếu call api thành công
    -> chuyển trạng thái processing_file_row.status=Success
    -> chuyển sang task tiếp theo (nếu còn)
    note: do ít thời gian nên bỏ qua Retry
    3.3. Tải file Result về, bổ sung thêm các row bị lỗi vào, upload file Result lên FileService
*/
func (j *jobHandleProcessingFileImpl) handleFileInProcessingStatus(ctx context.Context, file *fileprocessing.ProcessingFile) {
	taskGroupByRow, _ := j.fprService.GetAllRowsNeedToExecuteByJob(ctx, file.ID, fileprocessingrow.StatusInit)
	if len(taskGroupByRow) == 0 {
		logger.ErrorT("No row need to execute for fileId=%v", file.ID)
		return
	}

	for rowId, tasks := range taskGroupByRow {
		logger.Infof("----- Execute rowId=%v, with %v task(s) in fileId=%v", rowId, len(tasks), file.ID)
		providerClient := taskprovider.NewClient()
		previousResponse := make(map[int32]string) // map[task_index]=<response_string>
		for _, task := range tasks {
			// 1. If success, only get response, then go to next task
			if task.IsSuccessStatus() {
				previousResponse[task.TaskIndex] = task.TaskResponseRaw
				continue
			}

			// 2. Execute task
			logger.Infof("---------- Execute rowId=%v, taskId=%v in fileId=%v", rowId, task.TaskIndex, file.ID)
			requestBody, responseBody, isSuccess, messageRes := providerClient.Execute(task.RowDataRaw, task.TaskMapping, previousResponse)

			// 3. Update task status and save raw data for tracing
			updateRequest := toResponseResult(requestBody, responseBody, messageRes, isSuccess)
			_, err := j.fprService.UpdateAfterExecutingByJob(ctx, task.ID, updateRequest)
			if err != nil {
				logger.ErrorT("Update %v failed ---> ignore remaining tasks", fileprocessingrow.Name())
				break
			}
			if isSuccess { // task success -> put responseBody to previousResponse (map)
				previousResponse[task.TaskIndex] = responseBody
			} else {
				break // task failed  -> break loop, no execute next tasks
			}
		}
	}

}

func toResponseResult(requestBody map[string]string, responseBody string, messageRes string, isSuccess bool) fileprocessingrow.UpdateAfterExecutingByJob {
	// 1. Common value
	reqByte, _ := json.Marshal(requestBody)
	updateRequest := fileprocessingrow.UpdateAfterExecutingByJob{
		RequestRaw:   string(reqByte),
		ResponseRaw:  responseBody,
		Status:       1,
		ErrorDisplay: messageRes,
	}

	// 2. Set status
	if isSuccess {
		updateRequest.Status = fileprocessingrow.StatusSuccess
	} else {
		updateRequest.Status = fileprocessingrow.StatusFailed
	}

	return updateRequest
}

// Step 4 ==============================================================================================================

/*
 4. Thống kê số row đã thành công, cập nhật lại trạng thái file
    4.1. Thống kê total_success
    4.2. Nếu tất cả các row đã được xử lý (có thể thành công/thất bại)
    -> Collect các row bị xử lý lỗi, upload file Result lên FileService
    -> câp nhật processing_file.status=Finished, result_file_url, total_success
*/
func (j *jobHandleProcessingFileImpl) statisticAndUpdateFileStatus(ctx context.Context, file *fileprocessing.ProcessingFile) {
	isFinished, totalSuccess, totalFailed, errorDisplays, err := j.fprService.Statistics(file.ID)
	if err != nil {
		logger.ErrorT("Cannot statistics for file %v, err=%v", file.ID, err)
		return
	}

	status := file.Status
	resultFileUrl := file.ResultFileURL
	if isFinished { // all finished
		logger.InfoT("File is finished executing!!!!")
		// 1. Status is Finished
		status = fileprocessing.StatusFinished

		// 2. Create Result File if has error row
		if totalFailed > 0 {
			// 2.1 Inject Error Display to file
			fileDataBytes, err := excel.UpdateDataInColumnOfFile(file.FileURL, sheetImportDataName, columnErrorName, dataIndexStartInDataSheet, errorDisplays, false)
			if err != nil {
				logger.ErrorT("Update file with Error Display failed, err=%v", err)
				return
			}

			// 2.2. Gen result file name then Upload to file service
			fileName := utils.ExtractFileName(file.FileURL)
			resultFileName := fileName.FullNameWithSuffix("_result")
			res, err := j.fileService.UploadFileWithBytesData(fileDataBytes, resultFileName)
			if err != nil {
				logger.ErrorT("Upload result file %v failed, err=%v", resultFileName, err)
				return
			}
			resultFileUrl = res
		}
	}

	_, err = j.fpService.UpdateStatusWithStatistics(ctx, file.ID, status, totalSuccess, resultFileUrl)
	if err != nil {
		logger.ErrorT("Cannot update %v to failed, got error %v", fileprocessing.Name(), err)
		return
	}
}
