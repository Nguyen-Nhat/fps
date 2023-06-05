package configmapping

import (
	"context"
)

type (
	Service interface {
		FindByClientID(context.Context, int32) (*ConfigMapping, error)
	}

	serviceImpl struct {
		repo Repo
	}
)

var _ Service = &serviceImpl{}

func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) FindByClientID(ctx context.Context, clientID int32) (*ConfigMapping, error) {
	return s.repo.FindByClientID(ctx, clientID)
}
