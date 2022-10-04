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
		FindById(context.Context, int) (*FileAwardPoint, error)
		Save(context.Context, FileAwardPoint) (*FileAwardPoint, error)
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

// Implementation function ---------------------------------------------------------------------------------------------

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

func (r *RepoImpl) Save(ctx context.Context, fap FileAwardPoint) (*FileAwardPoint, error) {
	return save(ctx, r.client, fap)
}

// Other Public functions ----------------------------------------------------------------------------------------------

// SaveAll ... public function for using in Test
func SaveAll(ctx context.Context, client *ent.Client, faps []FileAwardPoint) ([]FileAwardPoint, error) {
	// 1. Build bulk
	bulk := make([]*ent.FileAwardPointCreate, len(faps))
	for i, fap := range faps {
		bulk[i] = mapFileAwardPoint(client, fap)
	}

	// 2. Create by bulk
	fapsSaved, err := client.FileAwardPoint.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save file award point to DB")
	}

	// 3. Map
	var res []FileAwardPoint
	for _, fapSaved := range fapsSaved {
		model := &FileAwardPoint{}
		if err := mapstructure.Decode(fapSaved, model); err != nil {
			return nil, fmt.Errorf("failed decode file award point info from DB")
		}
		res = append(res, *model)
	}

	// 4. Return
	return res, nil
}

// Private function ----------------------------------------------------------------------------------------------------

func save(ctx context.Context, client *ent.Client, fap FileAwardPoint) (*FileAwardPoint, error) {
	// 1. Create
	fapSaved, err := mapFileAwardPoint(client, fap).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save file award point to DB")
	}

	// 2. Map
	model := &FileAwardPoint{}
	if err := mapstructure.Decode(fapSaved, model); err != nil {
		return nil, fmt.Errorf("failed decode file award point info from DB")
	}

	// 3. Return
	return model, nil
}

// mapFileAwardPoint ... must update this function when schema is modified
func mapFileAwardPoint(client *ent.Client, fap FileAwardPoint) *ent.FileAwardPointCreate {
	return client.FileAwardPoint.Create().
		SetMerchantID(fap.MerchantId).
		SetDisplayName(fap.DisplayName).
		SetFileURL(fap.FileUrl).
		SetResultFileURL(fap.ResultFileUrl).
		SetStatus(fap.Status).
		SetStatsTotalRow(fap.StatsTotalRow).
		SetStatsTotalSuccess(fap.StatsTotalSuccess).
		SetCreatedAt(fap.CreatedAt).
		SetUpdatedAt(fap.UpdatedAt).
		SetCreatedBy(fap.CreatedBy).
		SetUpdatedBy(fap.UpdatedBy)
}
