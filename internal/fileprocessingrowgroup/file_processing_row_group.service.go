package fpRowGroup

type (
	Service interface {
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

// private method ------------------------------------------------------------------------------------------------------
