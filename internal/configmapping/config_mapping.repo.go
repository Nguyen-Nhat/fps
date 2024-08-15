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

func mapConfigMapping(client *ent.Client, cm ConfigMapping) *ent.ConfigMappingCreate {
	return client.ConfigMapping.Create().
		SetClientID(cm.ClientID).
		SetTenantID(cm.TenantID).
		SetMaxFileSize(cm.MaxFileSize).
		SetDataAtSheet(cm.DataAtSheet).
		SetDataStartAtRow(cm.DataStartAtRow).
		SetRequireColumnIndex(cm.RequireColumnIndex).
		SetErrorColumnIndex(cm.ErrorColumnIndex).
		SetMerchantAttributeName(cm.MerchantAttributeName).
		SetUsingMerchantAttrName(cm.UsingMerchantAttrName).
		SetInputFileType(cm.InputFileType).
		SetUIConfig(cm.UIConfig).
		SetCreatedBy(cm.CreatedBy)
}

func SaveAll(ctx context.Context, client *ent.Client, cmArr []ConfigMapping) ([]ConfigMapping, error) {
	// 1. Build bulk
	bulk := make([]*ent.ConfigMappingCreate, len(cmArr))
	for i, cm := range cmArr {
		bulk[i] = mapConfigMapping(client, cm)
	}

	// 2. Create by bulk
	cmSavedArr, err := client.ConfigMapping.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Map Result & return
	var res []ConfigMapping
	for _, cmSaved := range cmSavedArr {
		model := &ConfigMapping{*cmSaved}
		res = append(res, *model)
	}
	return res, nil
}
