package membertxn

import "context"

type (
	Service interface {
		GetByFileAwardPointId(context.Context, int32) ([]MemberTransaction, error)
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

func (s ServiceImpl) GetByFileAwardPointId(ctx context.Context, fileAwardPointId int32) ([]MemberTransaction, error) {
	return s.repo.FindByFileAwardPointId(ctx, fileAwardPointId)
}
