package fileawardpoint

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fileawardpoint"
)

type (
	Repo interface {
		FindById(context.Context, int) (*FileAwardPoint, error)
		Save(context.Context, FileAwardPoint) (*FileAwardPoint, error)
	}

	repoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

// NewRepo ...
func NewRepo(db *sql.DB) *repoImpl {
	drv := entsql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client}
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r *repoImpl) FindById(ctx context.Context, id int) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.Query().Where(fileawardpoint.ID(id)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular file award point: %w", err)
	}

	return &FileAwardPoint{*fap}, nil
}

func (r *repoImpl) Save(ctx context.Context, fap FileAwardPoint) (*FileAwardPoint, error) {
	return save(ctx, r.client, fap)
}

// Other Public functions ----------------------------------------------------------------------------------------------

// SaveAll ... public function for using in Test
// - needResult:
//   - TRUE 	=> return list of FileAwardPoint
//   - FALSE => return empty list
func SaveAll(ctx context.Context, client *ent.Client, fapArr []FileAwardPoint, needResult bool) ([]FileAwardPoint, error) {
	// 1. Build bulk
	bulk := make([]*ent.FileAwardPointCreate, len(fapArr))
	for i, fap := range fapArr {
		bulk[i] = mapFileAwardPoint(client, fap)
	}

	// 2. Create by bulk
	fapSavedArr, err := client.FileAwardPoint.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save file award point to DB")
	}

	// 3. Check if you NOT need result => return empty
	if !needResult {
		return []FileAwardPoint{}, nil
	}

	// 4. Map Result & return
	var res []FileAwardPoint
	for _, fapSaved := range fapSavedArr {
		model := &FileAwardPoint{*fapSaved}
		res = append(res, *model)
	}
	return res, nil
}

// Private function ----------------------------------------------------------------------------------------------------

func save(ctx context.Context, client *ent.Client, fap FileAwardPoint) (*FileAwardPoint, error) {
	// 1. Create
	fapSaved, err := mapFileAwardPoint(client, fap).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save file award point to DB")
	}

	// 2. Return
	return &FileAwardPoint{*fapSaved}, nil
}

// mapFileAwardPoint ... must update this function when schema is modified
func mapFileAwardPoint(client *ent.Client, fap FileAwardPoint) *ent.FileAwardPointCreate {
	return client.FileAwardPoint.Create().
		SetMerchantID(fap.MerchantID).
		SetDisplayName(fap.DisplayName).
		SetFileURL(fap.FileURL).
		SetResultFileURL(fap.ResultFileURL).
		SetStatus(fap.Status).
		SetStatsTotalRow(fap.StatsTotalRow).
		SetStatsTotalSuccess(fap.StatsTotalSuccess).
		SetCreatedAt(fap.CreatedAt).
		SetUpdatedAt(fap.UpdatedAt).
		SetCreatedBy(fap.CreatedBy).
		SetUpdatedBy(fap.UpdatedBy)
}
