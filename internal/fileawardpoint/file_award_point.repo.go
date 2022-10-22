package fileawardpoint

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"math"

	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		FindById(context.Context, int) (*FileAwardPoint, error)
		FindByStatuses(context.Context, []int16) ([]*FileAwardPoint, error)
		FindAndPaginationByMerchantId(context.Context, *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error)
		GetAllAndPagination(context.Context, *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error)
		Save(context.Context, FileAwardPoint) (*FileAwardPoint, error)
		UpdateStatusOne(context.Context, int, int16) (*FileAwardPoint, error)
		UpdateTotalRowOne(context.Context, int, int) (*FileAwardPoint, error)
		UpdateResultFileUrlOne(context.Context, int, string) (*FileAwardPoint, error)
		UpdateStatsTotalSuccessOne(context.Context, int, int) (*FileAwardPoint, error)
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

func (r *repoImpl) FindAndPaginationByMerchantId(ctx context.Context, dto *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error) {
	query := r.client.FileAwardPoint.Query().Where(fileawardpoint.MerchantID(int64(dto.MerchantId)))
	pagination, err := getPagination(ctx, query, dto.Page, dto.Size)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count from db while querying file award with merchantId")
	}

	faps, err := query.Limit(dto.Size).Offset((dto.Page - 1) * dto.Size).Order(ent.Desc(fileawardpoint.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed querying list file award point with merchantId: %w", err)
	}

	return mapEntArrToFileAwardPointArr(faps), pagination, nil
}

func (r *repoImpl) FindByStatuses(ctx context.Context, statuses []int16) ([]*FileAwardPoint, error) {
	faps, err := r.client.FileAwardPoint.Query().Where(fileawardpoint.StatusIn(statuses...)).All(ctx)

	if err != nil {
		logger.Errorf("fail to get file award point by status with status %#v", statuses)
		return nil, errors.New("fail to get file award point by status")
	}

	return mapEntArrToFileAwardPointArr(faps), nil
}

func (r *repoImpl) GetAllAndPagination(ctx context.Context, dto *GetListFileAwardPointDTO) ([]*FileAwardPoint, *response.Pagination, error) {
	query := r.client.FileAwardPoint.Query()
	pagination, err := getPagination(ctx, query, dto.Page, dto.Size)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count from db while querying all file award point")
	}

	faps, err := query.Limit(dto.Size).Offset((dto.Page - 1) * dto.Size).Order(ent.Desc(fileawardpoint.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed querying all file award point: %w", err)
	}

	return mapEntArrToFileAwardPointArr(faps), pagination, nil
}

func (r *repoImpl) Save(ctx context.Context, fap FileAwardPoint) (*FileAwardPoint, error) {
	return save(ctx, r.client, fap)
}

func (r *repoImpl) UpdateStatusOne(ctx context.Context, id int, status int16) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.UpdateOneID(id).SetStatus(status).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &FileAwardPoint{
		FileAwardPoint: *fap,
	}, nil
}
func (r *repoImpl) UpdateResultFileUrlOne(ctx context.Context, id int, resultFileUrl string) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.UpdateOneID(id).SetResultFileURL(resultFileUrl).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &FileAwardPoint{
		FileAwardPoint: *fap,
	}, nil
}
func (r *repoImpl) UpdateTotalRowOne(ctx context.Context, id int, totalRow int) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.UpdateOneID(id).SetStatsTotalRow(int32(totalRow)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &FileAwardPoint{
		FileAwardPoint: *fap,
	}, nil
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
		logger.Errorf("Cannot save to file award point, got: %v", err)
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
		SetNote(fap.Note).
		SetStatsTotalRow(fap.StatsTotalRow).
		SetStatsTotalSuccess(fap.StatsTotalSuccess).
		SetCreatedBy(fap.CreatedBy).
		SetUpdatedBy(fap.UpdatedBy)
}

func mapEntArrToFileAwardPointArr(arr ent.FileAwardPoints) []*FileAwardPoint {
	var result []*FileAwardPoint
	for _, v := range arr {
		result = append(result, &FileAwardPoint{*v})
	}
	return result
}

func getPagination(ctx context.Context, query *ent.FileAwardPointQuery, page int, size int) (*response.Pagination, error) {
	total, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &response.Pagination{
		CurrentPage: page,
		PageSize:    size,
		TotalItems:  total,
		TotalPage:   int(math.Ceil(float64(total) / float64(size))),
	}, nil
}

func (r *repoImpl) UpdateStatsTotalSuccessOne(ctx context.Context, id int, totalSuccess int) (*FileAwardPoint, error) {
	fap, err := r.client.FileAwardPoint.UpdateOneID(id).SetStatsTotalSuccess(int32(totalSuccess)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &FileAwardPoint{
		FileAwardPoint: *fap,
	}, nil
}
