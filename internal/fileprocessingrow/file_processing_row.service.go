package fileprocessingrow

import (
	"context"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		GetAllRowsNeedToExecuteByJob(context.Context, int, int) (map[int32][]*ProcessingFileRow, error)
		GetAllTasksForJobExecuteRowGroup(context.Context, int, int, string) ([]*ProcessingFileRow, error)

		SaveExtractedRowTaskFromFile(context.Context, int, []CreateProcessingFileRowJob) error
		UpdateAfterExecutingByJob(context.Context, int, UpdateAfterExecutingByJob) (*ProcessingFileRow, error)
		UpdateAfterExecutingByJobForListIDs(context.Context, []int, UpdateAfterExecutingByJob) error

		Statistics(int) (StatisticData, error)
	}

	ServiceImpl struct {
		repo Repo
	}
)

var _ Service = &ServiceImpl{}

func NewService(repo Repo) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

// SaveExtractedRowTaskFromFile ...
func (s *ServiceImpl) SaveExtractedRowTaskFromFile(ctx context.Context, fileID int, request []CreateProcessingFileRowJob) error {
	// 1. Clean all old data which relate to fileID
	_ = s.repo.DeleteByFileId(ctx, int64(fileID))

	// 2. Save by batch
	logger.Infof("----- Prepare SaveExtractedRowTaskFromFile with size = %+v", len(request))
	saveListFileFunc := func(subReq []CreateProcessingFileRowJob) error { return createRows(ctx, subReq, s) }
	return utils.BatchExecuting(500, request, saveListFileFunc)
}

func (s *ServiceImpl) GetAllRowsNeedToExecuteByJob(ctx context.Context, fileID int, limit int) (map[int32][]*ProcessingFileRow, error) {
	startAt := time.Now()
	pfrs, err := s.repo.FindRowsByFileIdForJobExecute(ctx, fileID, limit)
	logger.Infof("----- FindByFileIdAndStatusesForJob: limit=%v, totalResult=%v, executed time is %s", limit, len(pfrs), time.Since(startAt))
	if err != nil {
		return nil, err
	}

	// group task by rowIndex
	//var rowContainsFailedTasks []int32
	groupByRow := make(map[int32][]*ProcessingFileRow)
	for _, task := range pfrs {
		groupByRow[task.RowIndex] = append(groupByRow[task.RowIndex], task)
	}

	return groupByRow, nil
}

func (s *ServiceImpl) GetAllTasksForJobExecuteRowGroup(ctx context.Context, fileID int, taskIndex int, groupValue string) ([]*ProcessingFileRow, error) {
	return s.repo.FindByFileIdAndTaskIndexAndGroupValueAndStatus(ctx, int64(fileID), int32(taskIndex), groupValue, []int16{StatusWaitForGrouping, StatusRejected})
}

func (s *ServiceImpl) UpdateAfterExecutingByJob(ctx context.Context, id int,
	request UpdateAfterExecutingByJob) (*ProcessingFileRow, error) {
	// 1. If task failed -> remaining tasks of row are marked to
	if request.Status == StatusFailed {
		// todo ...
		task := request.Task
		affected, err := s.repo.UpdateStatusFromTask(ctx, task.FileID, task.RowIndex, task.TaskIndex)
		logger.Infof("Update remaining tasks to REJECT status with fileID=%v, rowIndex=%v, taskIndexFrom=%v ---> affected=%v, err=%+v",
			task.FileID, task.RowIndex, task.TaskIndex, affected, err)
	}

	// 2. Update task
	pfr, err := s.repo.UpdateByJob(ctx, id, request.TaskMapping, request.RequestCurl, request.ResponseRaw, request.Status, request.ErrorDisplay, request.ExecutedTime)
	if err != nil {
		logger.Errorf("Failed to update %v, error=%v, request=%+v", Name(), err, request)
		return nil, err
	}

	return pfr, nil
}

func (s *ServiceImpl) UpdateAfterExecutingByJobForListIDs(ctx context.Context, ids []int,
	request UpdateAfterExecutingByJob) error {
	err := s.repo.UpdateByJobForListIDs(ctx, ids, request.ResponseRaw, request.Status, request.ErrorDisplay, request.ExecutedTime)
	if err != nil {
		logger.Errorf("Failed to update %v, error=%v, ids=%+v, request=%+v", Name(), err, ids, request)
		return err
	}
	return nil
}

// Statistics ...
func (s *ServiceImpl) Statistics(fileID int) (StatisticData, error) {
	startAt := time.Now()
	statistics, err := s.repo.Statistics(int64(fileID))
	logger.Infof("----- Statistics: executed time is %s", time.Since(startAt))
	if err != nil {
		logger.Errorf("Error when get Statistics, err = %v", err)
		return StatisticData{}, err
	}

	total := len(statistics)
	totalSuccess := 0
	totalFailed := 0
	totalProcessed := 0
	totalWaiting := 0
	errorDisplays := make(map[int]string)
	for _, stats := range statistics {
		if stats.IsSuccessAll() {
			totalSuccess++
		} else if stats.IsContainsFailed() {
			totalFailed++
			errorDisplay := stats.GetErrorDisplay()
			errorDisplays[stats.RowIndex] = errorDisplay
		}

		if stats.IsProcessed() {
			totalProcessed++
		} else if stats.IsWaiting() {
			totalWaiting++
		}
	}

	logger.Infof("----- Statistic file %v: total=%v, totalProcessed=%v, totalSuccess=%v, totalFailed=%v, totalWaiting=%v",
		fileID, total, totalProcessed, totalSuccess, totalFailed, totalWaiting)

	statisticData := StatisticData{
		IsFinished:     isFinished(totalSuccess, totalFailed, total),
		ErrorDisplays:  errorDisplays,
		TotalProcessed: totalProcessed,
		TotalSuccess:   totalSuccess,
		TotalFailed:    totalFailed,
		TotalWaiting:   totalWaiting,
	}
	return statisticData, nil
}

// private method ------------------------------------------------------------------------------------------------------

func isFinished(totalSuccess int, totalFailed int, total int) bool {
	isSuccess := !(totalSuccess == 0 && totalFailed == 0) && total == totalSuccess+totalFailed
	return isSuccess
}

func createRows(ctx context.Context, subRequest []CreateProcessingFileRowJob, service *ServiceImpl) error {
	pfrArr := toProcessingFileRowArr(subRequest)
	if _, err := service.repo.SaveAll(ctx, pfrArr, false); err != nil {
		logger.Errorf("error when save all %v, got err %v", Name(), err)
		return err
	}
	return nil
}
