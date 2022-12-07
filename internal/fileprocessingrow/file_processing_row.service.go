package fileprocessingrow

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		SaveExtractedDataFromFile(context.Context, int, []CreateProcessingFileRowJob) error
		GetAllRowsNeedToExecuteByJob(context.Context, int, int16) (map[int32][]*ProcessingFileRow, error)

		UpdateAfterExecutingByJob(context.Context, int, UpdateAfterExecutingByJob) (*ProcessingFileRow, error)

		Statistics(int) (bool, int, int, []string, error)
	}

	ServiceImpl struct {
		repo Repo
	}
)

var _ Service = &ServiceImpl{}

// SaveExtractedDataFromFile ...
func (s *ServiceImpl) SaveExtractedDataFromFile(ctx context.Context, fileId int, request []CreateProcessingFileRowJob) error {
	// 1. Clean all old data which relate to fileId
	_ = s.repo.DeleteByFileId(ctx, int64(fileId))

	// 2. Save data
	pfrArr := toProcessingFileRowArr(request)
	if _, err := s.repo.SaveAll(ctx, pfrArr, false); err != nil {
		logger.Errorf("error when save all %v, got err %v", Name(), err)
		return err
	}
	logger.Infof("Saved %v extracted data from fileId=%v", len(request), fileId)
	return nil
}

func (s *ServiceImpl) GetAllRowsNeedToExecuteByJob(ctx context.Context, fileId int, status int16) (map[int32][]*ProcessingFileRow, error) {
	pfrs, err := s.repo.FindByFileIdAndStatusesForJob(ctx, fileId, status)
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
	logger.Infof("Prepare update %v with request=%+v", Name(), request)
	pfr, err := s.repo.UpdateByJob(ctx, id, request.RequestRaw, request.ResponseRaw, request.Status, request.ErrorDisplay)
	if err != nil {
		logger.Errorf("Failed to update %v, error=%v", Name(), err)
		return nil, err
	}

	return pfr, nil
}

// Statistics ... return (isFinished, totalSuccess, totalFailed, errorDisplays, error)
func (s *ServiceImpl) Statistics(fileId int) (bool, int, int, []string, error) {
	statistics, err := s.repo.Statistics(int64(fileId))
	if err != nil {
		logger.Errorf("Error when get Statistics, err = %v", err)
		return false, 0, 0, nil, err
	}

	total := len(statistics)
	totalSuccess := 0
	totalFailed := 0
	var errorDisplays []string
	for _, stats := range statistics {
		if stats.IsSuccessAll() {
			totalSuccess++
			errorDisplays = append(errorDisplays, "")
		} else if stats.IsContainsFailed() {
			totalFailed++
			errorDisplay := stats.GetErrorDisplay()
			errorDisplays = append(errorDisplays, errorDisplay)
		}
	}

	logger.Infof("----- Statistic file %v: total=%v, totalSuccess=%v, totalFailed=%v", fileId, total, totalSuccess, totalFailed)

	isSuccess := total == totalSuccess+totalFailed

	return isSuccess, totalSuccess, totalFailed, errorDisplays, nil
}

func NewService(repo Repo) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}
