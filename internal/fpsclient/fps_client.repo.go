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
		FindById(context.Context, int32) (*Client, error)
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

func (r *repoImpl) FindById(ctx context.Context, id int32) (*Client, error) {
	client, err := r.client.FpsClient.Get(ctx, int(id))
	if err != nil {
		logger.Errorf(err.Error())
		return nil, fmt.Errorf("failed to get client by id")
	}
	return &Client{*client}, nil
}

func SaveAll(ctx context.Context, client *ent.Client, fpsClients []Client) ([]Client, error) {
	// 1. Build bulk
	bulk := make([]*ent.FpsClientCreate, len(fpsClients))
	for idx, fpsClient := range fpsClients {
		bulk[idx] = mapFpsClient(client, fpsClient)
	}

	// 2. Create by bulk
	fcSavedArr, err := client.FpsClient.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Map Result & return
	var res []Client
	for _, fcSaved := range fcSavedArr {
		model := &Client{*fcSaved}
		res = append(res, *model)
	}
	return res, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func mapEntArrToClientArr(arr []*ent.FpsClient) []*Client {
	var result []*Client
	for _, v := range arr {
		result = append(result, &Client{*v})
	}
	return result
}

func mapFpsClient(client *ent.Client, cm Client) *ent.FpsClientCreate {
	return client.FpsClient.Create().
		SetClientID(cm.ClientID).
		SetName(cm.Name).
		SetDescription(cm.Description).
		SetImportFileTemplateURL(cm.ImportFileTemplateURL).
		SetCreatedBy(cm.CreatedBy)
}
