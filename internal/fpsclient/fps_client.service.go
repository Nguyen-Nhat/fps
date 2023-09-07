package fpsclient

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		GetListClients(context.Context, GetListClientDTO) ([]*Client, response.PaginationNew, error)
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

func (s *ServiceImpl) GetListClients(ctx context.Context, req GetListClientDTO) ([]*Client, response.PaginationNew, error) {
	clients, total, err := s.repo.FindByRequestAndPagination(ctx, req)
	if err != nil {
		logger.Infof("Error in GetListClients, err %+v", err)
		return nil, response.PaginationNew{}, err
	}

	pagination := response.GetPaginationNew(total, req.PageRequest)

	return clients, pagination, err
}
