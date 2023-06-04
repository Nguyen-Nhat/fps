package configmapping

import (
	"context"
	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
)

type (
	Repo interface {
		FindByClientID(context.Context, int32) (*ConfigMapping, error)
	}

	repoImpl struct {
		client *ent.Client
		sqlDB  *dbsql.DB
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

// NewRepo ...
func NewRepo(db *dbsql.DB) Repo {
	drv := sql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client, sqlDB: db}
}

func (r *repoImpl) FindByClientID(ctx context.Context, clientID int32) (*ConfigMapping, error) {
	cm, err := r.client.ConfigMapping.Query().Where(configmapping.ClientID(clientID)).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular config mapping: %w", err)
	}

	return &ConfigMapping{*cm}, nil
}
