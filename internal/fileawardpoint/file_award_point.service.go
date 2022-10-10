package fileawardpoint

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		GetFileAwardPoint(context.Context, *GetFileAwardPointDetailDTO) (*FileAwardPoint, error)
		GetListFileAwardPoint(context.Context, *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error)
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

// Implementation function ---------------------------------------------------------------------------------------------

func (s *ServiceImpl) GetFileAwardPoint(ctx context.Context, req *GetFileAwardPointDetailDTO) (*FileAwardPoint, error) {
	user := middleware.GetUserFromContext(ctx)
	logger.Infof("GetFileAwardPoint by %v", user.Email)
	return s.repo.FindById(ctx, req.Id)
}

func (s *ServiceImpl) GetListFileAwardPoint(ctx context.Context, req *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error) {
	merchantId := req.MerchantId
	var faps []*FileAwardPoint
	var pagination *response.Pagination
	if merchantId != 0 {
		temp, temp1, err := s.repo.FindAndPaginationByMerchantId(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		faps = temp
		pagination = temp1
	} else {
		temp, temp1, err := s.repo.GetAllAndPagination(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		faps = temp
		pagination = temp1
	}

	return faps, pagination, nil
}

func (s *ServiceImpl) CreateFileAwardPoint(ctx context.Context, req *CreateFileAwardPointReqDTO) (*CreateFileAwardPointResDTO, error) {
	user := middleware.GetUserFromContext(ctx)
	logger.Infof("GetFileAwardPoint by %v", user.Email)

	fileAwardRecord, err := s.repo.Save(ctx, FileAwardPoint{
		FileAwardPoint: ent.FileAwardPoint{
			MerchantID:  req.MerchantID,
			DisplayName: req.FileName,
			FileURL:     req.FileUrl,
			Note:        req.Note,
			CreatedBy:   user.Email,
			UpdatedBy:   user.Email,
		},
	})
	if err != nil {
		logger.Errorf("Cannot insert file award point, got %v", err)
		return nil, err
	}
	return &CreateFileAwardPointResDTO{
		FileAwardPointId: int32(fileAwardRecord.ID),
	}, nil
}
