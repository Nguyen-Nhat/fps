package fileawardpoint

import (
	"context"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		GetFileAwardPoint(context.Context, *GetFileAwardPointDetailDTO) (*FileAwardPoint, error)
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
