package fpRowGroup

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		SaveExtractedRowGroupFromFile(context.Context, int, []CreateRowGroupJob) error
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

// SaveExtractedRowGroupFromFile ...
func (s *ServiceImpl) SaveExtractedRowGroupFromFile(ctx context.Context, fileID int, request []CreateRowGroupJob) error {
	// 1. Clean all old data which relate to fileID
	_ = s.repo.DeleteByFileId(ctx, int64(fileID))

	// 2. Save by batch
	logger.Infof("----- Prepare SaveExtractedRowGroupFromFile with size = %+v", len(request))
	saveListFileFunc := func(subReq []CreateRowGroupJob) error { return createRowGroups(ctx, subReq, s) }
	return utils.BatchExecuting(500, request, saveListFileFunc)
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
