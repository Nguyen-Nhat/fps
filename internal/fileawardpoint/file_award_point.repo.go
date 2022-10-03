package fileawardpoint

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fileawardpoint"
	"github.com/mitchellh/mapstructure"
)

type (
	Repo interface {
		FindById(ctx context.Context, id int) (*FileAwardPoint, error)
	}

	RepoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &RepoImpl{} // only for mention that RepoImpl implement Repo

// NewRepo ...
func NewRepo(db *sql.DB) *RepoImpl {
	drv := entsql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &RepoImpl{client: client}
}

func (r *RepoImpl) FindById(ctx context.Context, id int) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.Query().Where(fileawardpoint.ID(id)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular file award point: %w", err)
	}
	model := &FileAwardPoint{}
	if err := mapstructure.Decode(fap, model); err != nil {
		return nil, fmt.Errorf("failed decode file award point info from DB")
	}
	return model, nil
}
