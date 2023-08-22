package configtask

import (
	"context"
	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configtask"
)

type (
	Repo interface {
		FindByConfigMappingID(context.Context, int32) ([]*ConfigTask, error)
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

func (r *repoImpl) FindByConfigMappingID(ctx context.Context, configMappingID int32) ([]*ConfigTask, error) {
	cts, err := r.client.ConfigTask.Query().
		Where(configtask.ConfigMappingID(configMappingID)).
		Order(ent.Asc(configtask.FieldTaskIndex)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying config task: %w", err)
	}

	return mapEntArrToEntityArr(cts), nil
}

// Private methods -----------------------------------------------------------------------------------------------------

func mapEntArrToEntityArr(arr ent.ConfigTasks) []*ConfigTask {
	var result []*ConfigTask
	for _, v := range arr {
		result = append(result, &ConfigTask{*v})
	}
	return result
}
