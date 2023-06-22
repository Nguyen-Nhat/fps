package fileprocessing

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/faltservice"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	Service interface {
		CreateFileProcessing(ctx context.Context, req *CreateFileProcessingReqDTO) (*CreateFileProcessingResDTO, error)
		GetFileProcessHistory(ctx context.Context, req *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error)

		FindById(context.Context, int) (*ProcessingFile, error)
		GetListFileByStatuses(context.Context, []int16) ([]*ProcessingFile, error)

		UpdateToFailedStatusWithErrorMessage(context.Context, int, ErrorDisplay, *string) (*ProcessingFile, error)
		UpdateToProcessingStatusWithExtractedData(context.Context, int, int, int) (*ProcessingFile, error)
		UpdateStatusWithStatistics(context.Context, int, int16, int, int, string) (*ProcessingFile, error)
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

// CreateFileProcessing ... Create new file processing. If display name is not provided, it will be extracted from file name
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
			ClientID:       req.ClientID,
			DisplayName:    displayName,
			FileURL:        req.FileURL,
			Status:         StatusInit,
			CreatedBy:      req.CreatedBy,
			FileParameters: req.FileParameters,
		},
	})
	if err != nil {
		logger.Errorf("Cannot create file processing, got: %v", err)
		return nil, err
	}

	// 2.1. Create ProcessingFile in f-alt-server
	if createdProcessingFile != nil {
		go func() {
			_, _ = faltservice.CreateProcessingFile(&faltservice.ProcessingFileParse{
				Status:    StatusInit,
				FpsFileID: createdProcessingFile.ID,
			})
		}()
	}

	return &CreateFileProcessingResDTO{
		ProcessFileID: int32(createdProcessingFile.ID),
	}, err
}

func (s *ServiceImpl) FindById(ctx context.Context, id int) (*ProcessingFile, error) {
	pfs, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return pfs, nil
}

func (s *ServiceImpl) GetListFileByStatuses(ctx context.Context, statuses []int16) ([]*ProcessingFile, error) {
	pfs, err := s.repo.FindByStatuses(ctx, statuses)
	if err != nil {
		return nil, err
	}

	return pfs, nil
}

func (s *ServiceImpl) UpdateToFailedStatusWithErrorMessage(ctx context.Context, id int, errorMessage ErrorDisplay, resultFileURL *string) (*ProcessingFile, error) {
	pf, err := s.repo.UpdateStatusAndErrorDisplay(ctx, id, StatusFailed, errorMessage, resultFileURL)
	if err != nil {
		return nil, err
	}

	return pf, nil
}

func (s *ServiceImpl) UpdateStatusWithStatistics(ctx context.Context, id int, status int16, totalProcessed int, totalSuccess int, resultFileUrl string) (*ProcessingFile, error) {
	pf, err := s.repo.UpdateStatusAndStatsAndResultFileUrl(ctx, id, status, totalProcessed, totalSuccess, resultFileUrl)
	if err != nil {
		return nil, err
	}

	return pf, nil
}

func (s *ServiceImpl) UpdateToProcessingStatusWithExtractedData(ctx context.Context, id int, totalMapping int, totalRow int) (*ProcessingFile, error) {
	fp, err := s.repo.UpdateByExtractedData(ctx, id, StatusProcessing, totalMapping, totalRow)
	if err != nil {
		logger.Errorf("Update %v failed, got err %v", Name(), err)
	}
	return fp, nil
}

func (s *ServiceImpl) GetFileProcessHistory(ctx context.Context, req *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error) {
	var files []*ProcessingFile
	var pagination *response.Pagination

	files, pagination, err := s.repo.FindByClientIdAndPagination(ctx, req)
	if err != nil {
		logger.Infof("Error in FindByClientIdAndPagination")
		return nil, nil, err
	}

	return files, pagination, err
}
