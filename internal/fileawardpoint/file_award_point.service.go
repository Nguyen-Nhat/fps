package fileawardpoint

import (
	"context"
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

// Implementation function

func (s *ServiceImpl) GetFileAwardPoint(ctx context.Context, req *GetFileAwardPointDetailDTO) (*FileAwardPoint, error) {
	return s.repo.FindById(ctx, req.Id)
}
