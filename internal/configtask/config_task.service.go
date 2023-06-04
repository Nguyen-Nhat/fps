package configtask

import "context"

type (
	Service interface {
		FindByConfigMappingID(context.Context, int32) ([]*ConfigTask, error)
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

func (s *serviceImpl) FindByConfigMappingID(ctx context.Context, clientID int32) ([]*ConfigTask, error) {
	return s.repo.FindByConfigMappingID(ctx, clientID)
}
