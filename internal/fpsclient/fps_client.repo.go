package fpsclient

import (
	"context"
	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fpsclient"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		FindByRequestAndPagination(context.Context, GetListClientDTO) ([]*Client, int, error)
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

func (r *repoImpl) FindByRequestAndPagination(ctx context.Context, dto GetListClientDTO) ([]*Client, int, error) {
	query := r.client.FpsClient.Query()

	// search name like %input%
	if len(dto.Name) > 0 {
		query = query.Where(fpsclient.NameContains(dto.Name))
	}

	total, err := query.Count(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, 0, fmt.Errorf("failed to count from db while querying client")
	}

	clients, err := query.
		Limit(dto.PageSize).
		Offset((dto.Page - 1) * dto.PageSize).
		Order(ent.Desc(fpsclient.FieldID)).
		All(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, 0, fmt.Errorf("failed querying all client")
	}

	return mapEntArrToClientArr(clients), total, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func mapEntArrToClientArr(arr []*ent.FpsClient) []*Client {
	var result []*Client
	for _, v := range arr {
		result = append(result, &Client{*v})
	}
	return result
}
