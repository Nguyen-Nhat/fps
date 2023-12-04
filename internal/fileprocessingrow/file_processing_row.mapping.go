package fileprocessingrow

import (
	"sort"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

func toProcessingFileRow(request CreateProcessingFileRowJob) ProcessingFileRow {
	return ProcessingFileRow{
		ProcessingFileRow: ent.ProcessingFileRow{
			FileID:       int64(request.FileID),
			RowIndex:     int32(request.RowIndex),
			RowDataRaw:   request.RowDataRaw,
			TaskIndex:    int32(request.TaskIndex),
			TaskMapping:  request.TaskMapping,
			GroupByValue: request.GroupByValue,
			Status:       StatusInit,
			ExecutedTime: -1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
}

func toArrGetListFileRowsItem(taskMap map[int][]*ProcessingFileRow, fileID int) []GetListFileRowsItem {
	rowIndexes := make([]int, 0, len(taskMap))
	for rowIndex := range taskMap {
		rowIndexes = append(rowIndexes, rowIndex)
	}
	sort.Ints(rowIndexes) // order rowIndex

	var result []GetListFileRowsItem
	for _, rowIndex := range rowIndexes {
		tasksInRow := taskMap[rowIndex]
		task := tasksInRow[0]
		tmp := GetListFileRowsItem{
			FileID:       fileID,
			RowIndex:     rowIndex,
			RowDataRaw:   task.RowDataRaw,
			ExecutedTime: -1,
			Tasks:        converter.Map(tasksInRow, toTaskInRowItem),
		}
		result = append(result, tmp)
	}
	return result
}

// ---------------------------------------------------------------------------------------------------------------------

func toTaskInRowItem(taskInRow *ProcessingFileRow) TaskInRowItem {
	taskName := getTaskName(taskInRow)

	return TaskInRowItem{
		TaskIndex:       int(taskInRow.TaskIndex),
		TaskRequestCurl: taskInRow.TaskRequestCurl,
		TaskResponseRaw: taskInRow.TaskResponseRaw,
		TaskName:        taskName,
		GroupByValue:    taskInRow.GroupByValue,
		Status:          taskInRow.Status,
		ErrorDisplay:    taskInRow.ErrorDisplay,
		ExecutedTime:    int(taskInRow.ExecutedTime),
		CreatedAt:       taskInRow.CreatedAt.UnixMilli(),
		UpdatedAt:       taskInRow.UpdatedAt.UnixMilli(),
	}
}

func getTaskName(taskInRow *ProcessingFileRow) string {
	// 1. Load Data and Mapping
	configMapping, err := converter.StringJsonToStruct("config mapping", taskInRow.TaskMapping, configloader.ConfigMappingMD{})
	if err != nil {
		return ""
	}

	// 2. Get task
	for _, taskMD := range configMapping.Tasks {
		if taskMD.TaskIndex == int(taskInRow.TaskIndex) {
			return taskMD.TaskName // return when find task
		}
	}

	// 3. return no name
	return "no name"
}
