package fpRowGroup

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		FindRowGroupForJobExecute(context.Context, int) (map[int][]*ProcessingFileRowGroup, error)

		SaveExtractedRowGroupFromFile(context.Context, int, []CreateRowGroupJob) error
		UpdateAfterExecutingByJob(context.Context, int, UpdateAfterExecutingByJob) (*ProcessingFileRowGroup, error)
	}

	ServiceImpl struct {
		repo Repo
	}
)

var _ Service = &ServiceImpl{}

func NewService(repo Repo) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) FindRowGroupForJobExecute(ctx context.Context, fileID int) (map[int][]*ProcessingFileRowGroup, error) {
	// 1. Get in DB
	rowGroups, err := s.repo.FindByFileIdAndStatusIn(ctx, int64(fileID), []int16{StatusInit, StatusCalledApiFail, StatusCalledApiSuccess})
	if err != nil {
		return nil, err
	}

	// 2. Group by taskIndex
	rowGroupMap := make(map[int][]*ProcessingFileRowGroup)
	for _, rowGroup := range rowGroups {
		taskIndex := int(rowGroup.TaskIndex)
		if rgs, existed := rowGroupMap[taskIndex]; existed {
			rowGroupMap[taskIndex] = append(rgs, rowGroup)
		} else {
			rowGroupMap[taskIndex] = []*ProcessingFileRowGroup{rowGroup}
		}
	}

	return rowGroupMap, nil
}

// SaveExtractedRowGroupFromFile ...
func (s *ServiceImpl) SaveExtractedRowGroupFromFile(ctx context.Context, fileID int, request []CreateRowGroupJob) error {
	// 1. Clean all old data which relate to fileID
	_ = s.repo.DeleteByFileId(ctx, int64(fileID))

	// 2. Save by batch
	logger.Infof("----- Prepare SaveExtractedRowGroupFromFile with size = %+v", len(request))
	saveListFileFunc := func(subReq []CreateRowGroupJob) error { return createRowGroups(ctx, subReq, s) }
	return utils.BatchExecuting(500, request, saveListFileFunc)
}

func (s *ServiceImpl) UpdateAfterExecutingByJob(ctx context.Context, id int,
	request UpdateAfterExecutingByJob) (*ProcessingFileRowGroup, error) {
	pfr, err := s.repo.UpdateByJob(ctx, id, request.RequestCurl, request.ResponseRaw, request.Status, request.ErrorDisplay, request.ExecutedTime)
	if err != nil {
		logger.Errorf("Failed to update %v, error=%v, request=%+v", Name(), err, request)
		return nil, err
	}

	return pfr, nil
}

// private method ------------------------------------------------------------------------------------------------------

func createRowGroups(ctx context.Context, subRequest []CreateRowGroupJob, service *ServiceImpl) error {
	pfrArr := toProcessingFileRowGroupArr(subRequest)
	if _, err := service.repo.SaveAll(ctx, pfrArr, false); err != nil {
		logger.Errorf("error when save all %v, got err %v", Name(), err)
		return err
	}
	return nil
}
