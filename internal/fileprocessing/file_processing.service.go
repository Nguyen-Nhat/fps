package fileprocessing

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		CreateFileProcessing(ctx context.Context, req *CreateFileProcessingReqDTO) (*CreateFileProcessingResDTO, error)
	}

	ServiceImpl struct {
		repo Repo
	}
)

var _ Service = &ServiceImpl{}

// CreateFileProcessing: Create new file processing. If display name is not provided, it will be extract from file name
func (s *ServiceImpl) CreateFileProcessing(ctx context.Context, req *CreateFileProcessingReqDTO) (*CreateFileProcessingResDTO, error) {

	// 1. Preprocess data
	// Get file name from file URL in case display name was not provided
	displayName := req.DisplayName
	if displayName == "" {
		logger.Warnf("Not receive display name from request. Extract from file URL %s", req.FileURL)
		displayName = utils.ExtractFileName(req.FileURL).FullName
	}

	// 2. Create new file processing
	createdProcessingFile, err := s.repo.Save(ctx, ProcessingFile{
		ProcessingFile: ent.ProcessingFile{
			ClientID:    req.ClientID,
			DisplayName: displayName,
			FileURL:     req.FileURL,
			Status:      StatusInit,
			CreatedBy:   req.CreatedBy,
		},
	})
	if err != nil {
		logger.Errorf("Cannot create file processing, got: %v", err)
		return nil, err
	}

	return &CreateFileProcessingResDTO{
		ProcessFileID: int32(createdProcessingFile.ID),
	}, err
}

func NewService(repo Repo) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}
