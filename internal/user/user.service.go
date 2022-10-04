package user

import (
	"context"
)

type Service interface {
	CreateUser(context.Context, *CreateUserRequestDTO) (*User, error)
}

type ServiceImpl struct {
	repo Repo
}

func NewService(repo Repo) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateUser(ctx context.Context, user *CreateUserRequestDTO) (*User, error) {
	u := toUserFromCreateDTO(user)
	return s.repo.CreateUser(ctx, u)
}
