package flatten

import (
	"fmt"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	fpRowGroup "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// validateAndBuildRowGroupData ...
func validateAndBuildRowGroupData(fileID int, configMappingMD configloader.ConfigMappingMD, configMappingsWithData []*configloader.ConfigMappingMD) ([]ErrorRow, []fpRowGroup.CreateRowGroupJob) {
	// config mapping not have config for grouping
	if !configMappingMD.IsSupportGrouping() {
		return nil, nil
	}

	var createRowGroupJobs []fpRowGroup.CreateRowGroupJob
	for _, taskMD := range configMappingMD.Tasks {
		// 1. case no support -> ignore
		if !taskMD.RowGroup.IsSupportGrouping() {
			continue
		}

		// 2. grouping rows and check error/no group
		groupMap, errorRowGroups := groupingRowsBaseOnFileData(configMappingsWithData, taskMD)
		if len(errorRowGroups) > 0 {
			return errorRowGroups, nil // -> return error, do NOT check next task
		}
		if len(groupMap) <= 0 {
			logger.Warnf("--- Task %d support Grouping but cannot find any group in file", taskMD.TaskIndex)
			continue // -> go to next task, this case NOT happens, only check to avoid bug
		}

		// 3. build data for preparing to save to DB
		logger.Infof("--- Find groups: size=%v, data=%+v", len(groupMap), groupMap)
		for groupValue, rowIDs := range groupMap {
			pfr := fpRowGroup.CreateRowGroupJob{
				FileID:       fileID,
				TaskIndex:    taskMD.TaskIndex,
				GroupByValue: groupValue,
				TotalRows:    len(rowIDs),
				RowIndexList: utils.JsonString(rowIDs),
			}
			createRowGroupJobs = append(createRowGroupJobs, pfr)
		}

	}

	return nil, createRowGroupJobs
}

func groupingRowsBaseOnFileData(configMappingsWithData []*configloader.ConfigMappingMD, taskMD configloader.ConfigTaskMD) (map[string][]int, []ErrorRow) {
	var errorRows []ErrorRow
	groupMap := make(map[string][]int) // map[group_value][]rowID

	for _, datum := range configMappingsWithData {
		// 2.1. Get group value
		task := datum.GetConfigTaskMD(taskMD.TaskIndex)
		groupValue, errorRow := getGroupValueInRowOfTask(task)
		if errorRow != nil {
			errorRows = append(errorRows, *errorRow)
			continue // next row, we don't choose `break` because we should show all rows that has error
		}

		// 2.2. Set group value to config task
		datum.RowGroupValue = groupValue

		// 2.3. Group and count rows
		if rowIDs, ok := groupMap[groupValue]; ok {
			rowIDsTmp := append(rowIDs, task.ImportRowIndex)
			groupMap[groupValue] = rowIDsTmp // still update counter, we will show all rows that belongs to this group
			if len(rowIDsTmp) > task.RowGroup.GroupSizeLimit {
				logger.Errorf("file invalid by total item in group %s greater than %d", fmt.Sprintf("[%s]", groupValue), taskMD.RowGroup.GroupSizeLimit)
				errTmp := ErrorRow{task.ImportRowIndex, fmt.Sprintf("Total item in group %s is greater than %d", fmt.Sprintf("[%s]", groupValue), taskMD.RowGroup.GroupSizeLimit)}
				errorRows = append(errorRows, errTmp)
			}
		} else {
			groupMap[groupValue] = []int{task.ImportRowIndex}
		}
	}

	return groupMap, errorRows
}

// getGroupValueInRowOfTask ...
func getGroupValueInRowOfTask(task configloader.ConfigTaskMD) (string, *ErrorRow) {
	columnIndexes := task.RowGroup.GroupByColumns
	// 1. Get group value
	var valueInColumns []string
	for _, columnIndex := range columnIndexes {
		if columnIndex < 0 || columnIndex >= len(task.ImportRowData) {
			logger.Errorf("missing data at row %d, column %d", task.ImportRowIndex, columnIndex)
			return "", &ErrorRow{task.ImportRowIndex, fmt.Sprintf("Missing data at columns %d", columnIndex)}
		}

		value := task.ImportRowData[columnIndex]
		valueInColumns = append(valueInColumns, value)
	}
	groupValuesStr := strings.Join(valueInColumns, "|")
	emptyGroup := strings.Join(make([]string, len(valueInColumns)), "|")
	if groupValuesStr == emptyGroup { // case all group column is empty
		logger.Errorf("missing data at row %d, column %+v", task.ImportRowIndex, columnIndexes)
		return "", &ErrorRow{task.ImportRowIndex, fmt.Sprintf("Missing data at columns %+v", columnIndexes)}
	}

	return groupValuesStr, nil
}
