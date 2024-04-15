package updatestatus

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/filewriter"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

func (job *jobUpdateStatus) buildResultFileAndUpload(ctx context.Context,
	file fileprocessing.ProcessingFile, cfgMapping *configmapping.ConfigMapping,
	errorDisplays map[int]string) (string, error) {

	// 1. Inject Error Display to file
	var fileDataBytes *bytes.Buffer
	var resultFileConfigs []configloader.ResultFileConfigMD
	var err error

	// 1. Init file writer
	dataStartAtRow := int(cfgMapping.DataStartAtRow)
	fw, err := filewriter.NewFileWriter(file.FileURL, cfgMapping.DataAtSheet, dataStartAtRow, file.ExtFileRequest, cfgMapping.OutputFileType)
	if err != nil {
		return "", err
	}

	// 2. Set data to filed
	_ = json.Unmarshal([]byte(cfgMapping.ResultFileConfig), &resultFileConfigs)
	if len(resultFileConfigs) == 0 { // 2.1. case normal -> write data to error column
		if err = fw.UpdateDataInColumnOfFile(cfgMapping.ErrorColumnIndex, errorDisplays); err != nil {
			logger.ErrorT("Update file with Error Display failed, err=%v", err)
			return "", err
		}
	} else { // 2.2. special case -> write data based on ResultFileConfigMD
		// 2.2.1. Get list task from DB
		taskIDs := toTaskIDs(resultFileConfigs)
		resultAsync, er := job.fprService.GetResultAsync(ctx, file.ID, taskIDs...)
		if er != nil {
			return "", er
		}

		// 2.2.2. Explore each column and write data to that column
		for _, resultFileConfig := range resultFileConfigs {
			// Collect data for column
			columnData := getColumnData(resultAsync, resultFileConfig)
			// Write to column
			if err = fw.UpdateDataInColumnOfFile(resultFileConfig.ColumnKey, columnData); err != nil {
				logger.ErrorT("Update file with Error Display failed, err=%v", err)
				return "", err
			}
		}
	}

	// 3. Get file bytes & return
	fileDataBytes, err = fw.GetFileBytes()
	if err != nil {
		return "", err
	}

	// 4. Gen result file name then Upload to file service
	resultFileName := utils.GetResultFileName(file.DisplayName)
	return job.fileService.UploadFileWithBytesData(fileDataBytes, fw.OutputFileContentType(), resultFileName)
}

func getColumnData(results []fileprocessingrow.ResultAsyncDAO, resultFileConfig configloader.ResultFileConfigMD) map[int]string {
	columnData := make(map[int]string)
	for _, result := range results {
		if result.TaskIndex != resultFileConfig.ValueInTaskID {
			continue
		}

		valueNeedToFill := gjson.Get(result.ResultAsync, resultFileConfig.ValuePath).String()
		columnData[result.RowIndex] = valueNeedToFill
	}
	return columnData
}

func toTaskIDs(resultFileConfigs []configloader.ResultFileConfigMD) []int32 {
	idMap := make(map[int32]int32)
	var taskIDs []int32
	for _, cfg := range resultFileConfigs {
		taskID := cfg.ValueInTaskID
		if _, existed := idMap[taskID]; !existed {
			idMap[taskID] = taskID
			taskIDs = append(taskIDs, taskID)
		}
	}
	return taskIDs
}
