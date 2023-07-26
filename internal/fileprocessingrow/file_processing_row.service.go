package fileprocessingrow

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		SaveExtractedDataFromFile(context.Context, int, []CreateProcessingFileRowJob) error
		GetAllRowsNeedToExecuteByJob(context.Context, int, int16) (map[int32][]*ProcessingFileRow, error)

		UpdateAfterExecutingByJob(context.Context, int, UpdateAfterExecutingByJob) (*ProcessingFileRow, error)

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

// SaveExtractedDataFromFile ...
func (s *ServiceImpl) SaveExtractedDataFromFile(ctx context.Context, fileID int, request []CreateProcessingFileRowJob) error {
	// 1. Clean all old data which relate to fileID
	_ = s.repo.DeleteByFileId(ctx, int64(fileID))

	// 2. Save by batch
	logger.Infof("----- Prepare SaveExtractedDataFromFile with size = %+v", len(request))
	saveListFileFunc := func(subReq []CreateProcessingFileRowJob) error { return createRows(ctx, subReq, s) }
	return utils.BatchExecuting(500, request, saveListFileFunc)
}

func (s *ServiceImpl) GetAllRowsNeedToExecuteByJob(ctx context.Context, fileID int, status int16) (map[int32][]*ProcessingFileRow, error) {
	pfrs, err := s.repo.FindByFileIdAndStatusesForJob(ctx, fileID, status)
	if err != nil {
		return nil, err
	}

	// group task by rowIndex
	var rowContainsFailedTasks []int32
	groupByRow := make(map[int32][]*ProcessingFileRow)
	for _, task := range pfrs {
		groupByRow[task.RowIndex] = append(groupByRow[task.RowIndex], task)
		if task.IsFailedStatus() {
			rowContainsFailedTasks = append(rowContainsFailedTasks, task.RowIndex)
		}
	}

	// remove row which has at least on task failed
	for _, rowIndex := range rowContainsFailedTasks {
		delete(groupByRow, rowIndex)
	}

	return groupByRow, nil
}

func (s *ServiceImpl) UpdateAfterExecutingByJob(ctx context.Context, id int,
	request UpdateAfterExecutingByJob) (*ProcessingFileRow, error) {
	pfr, err := s.repo.UpdateByJob(ctx, id, request.RequestCurl, request.ResponseRaw, request.Status, request.ErrorDisplay, request.ExecutedTime)
	if err != nil {
		logger.Errorf("Failed to update %v, error=%v, request=%+v", Name(), err, request)
		return nil, err
	}

	return pfr, nil
}

// Statistics ...
func (s *ServiceImpl) Statistics(fileID int) (StatisticData, error) {
	statistics, err := s.repo.Statistics(int64(fileID))
	if err != nil {
		logger.Errorf("Error when get Statistics, err = %v", err)
		return StatisticData{}, err
	}

	total := len(statistics)
	totalSuccess := 0
	totalFailed := 0
	totalProcessed := 0
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
		}
	}

	logger.Infof("----- Statistic file %v: total=%v, totalSuccess=%v, totalFailed=%v", fileID, total, totalSuccess, totalFailed)

	statisticData := StatisticData{
		IsFinished:     isFinished(totalSuccess, totalFailed, total),
		ErrorDisplays:  errorDisplays,
		TotalProcessed: totalProcessed,
		TotalSuccess:   totalSuccess,
		TotalFailed:    totalFailed,
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
