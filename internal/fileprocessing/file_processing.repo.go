package fileprocessing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfile"

	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error)
		FindByClientIdAndPagination(context.Context, *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error)
		FindByID(context.Context, int) (*ProcessingFile, error)
		FindByStatuses(context.Context, []int16) ([]*ProcessingFile, error)
		UpdateStatusOne(context.Context, int, int16) (*ProcessingFile, error)
		UpdateStatusAndErrorDisplay(context.Context, int, int16, ErrorDisplay) (*ProcessingFile, error)
		UpdateStatusAndStatsAndResultFileUrl(context.Context, int, int16, int, string) (*ProcessingFile, error)
		UpdateByExtractedData(ctx context.Context, id int, status int16, totalMapping int, statsTotalRow int) (*ProcessingFile, error)
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

// Save ... Implementation function ------------------------------------------------------------------------------------
func (r *repoImpl) Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error) {
	return save(ctx, r.client, fp)
}

func (r *repoImpl) FindByID(ctx context.Context, id int) (*ProcessingFile, error) {
	fp, err := r.client.ProcessingFile.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular file award point: %w", err)
	}

	return &ProcessingFile{*fp}, nil
}

func (r *repoImpl) FindByStatuses(ctx context.Context, statuses []int16) ([]*ProcessingFile, error) {
	pfs, err := r.client.ProcessingFile.Query().Where(processingfile.StatusIn(statuses...)).All(ctx)

	if err != nil {
		logger.Errorf("fail to get file award point by status with status %#v", statuses)
		return nil, errors.New("fail to get file award point by status")
	}

	return mapEntArrToProcessingFileArr(pfs), nil
}

func (r *repoImpl) UpdateStatusOne(ctx context.Context, id int, status int16) (*ProcessingFile, error) {
	fap, err := r.client.ProcessingFile.UpdateOneID(id).SetStatus(status).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &ProcessingFile{
		ProcessingFile: *fap,
	}, nil
}

func (r *repoImpl) UpdateStatusAndErrorDisplay(ctx context.Context, id int, status int16, errorDisplay ErrorDisplay) (*ProcessingFile, error) {
	fap, err := r.client.ProcessingFile.UpdateOneID(id).
		SetStatus(status).
		SetErrorDisplay(string(errorDisplay)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &ProcessingFile{
		ProcessingFile: *fap,
	}, nil
}

func (r *repoImpl) UpdateStatusAndStatsAndResultFileUrl(ctx context.Context, id int, status int16, totalSuccess int,
	resultFileUrl string) (*ProcessingFile, error) {
	fap, err := r.client.ProcessingFile.UpdateOneID(id).
		SetStatus(status).
		SetStatsTotalSuccess(int32(totalSuccess)).
		SetResultFileURL(resultFileUrl).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &ProcessingFile{
		ProcessingFile: *fap,
	}, nil
}

func (r *repoImpl) UpdateByExtractedData(ctx context.Context, id int, status int16, totalMapping int, statsTotalRow int) (*ProcessingFile, error) {
	pf, err := r.client.ProcessingFile.UpdateOneID(id).
		SetStatus(status).
		SetTotalMapping(int32(totalMapping)).
		SetStatsTotalRow(int32(statsTotalRow)).
		Save(ctx)

	if err != nil {
		return nil, err
	}
	return &ProcessingFile{
		ProcessingFile: *pf,
	}, nil
}

// private function ---------------------------------------------------------------------------------------------

func save(ctx context.Context, client *ent.Client, fp ProcessingFile) (*ProcessingFile, error) {
	// 1. Create
	fpSaved, err := mapProcessingFile(client, fp).Save(ctx)

	if err != nil {
		logger.Errorf("Cannot save to file processing, got: %v", err)
		return nil, fmt.Errorf("failed to save file processing to DB")
	}

	// 2. Return
	return &ProcessingFile{*fpSaved}, nil
}

func mapProcessingFile(client *ent.Client, fp ProcessingFile) *ent.ProcessingFileCreate {
	return client.ProcessingFile.Create().
		SetClientID(fp.ClientID).
		SetDisplayName(fp.DisplayName).
		SetFileURL(fp.FileURL).
		SetResultFileURL(fp.ResultFileURL).
		SetStatus(fp.Status).
		SetTotalMapping(fp.TotalMapping).
		SetStatsTotalRow(fp.StatsTotalRow).
		SetStatsTotalSuccess(fp.StatsTotalSuccess).
		SetErrorDisplay(fp.ErrorDisplay).
		SetCreatedBy(fp.CreatedBy)
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r *repoImpl) FindByClientIdAndPagination(ctx context.Context, dto *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error) {
	query := r.client.ProcessingFile.Query().Where(processingfile.ClientID(dto.ClientId))

	total, err := query.Count(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, nil, fmt.Errorf("failed to count from db while querying all processing file")
	}
	pagination := response.GetPagination(total, dto.Page, dto.Size)

	fps, err := query.Limit(dto.Size).Offset((dto.Page - 1) * dto.Size).Order(ent.Desc(processingfile.FieldCreatedAt)).All(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, nil, fmt.Errorf("failed querying all processing file")
	}

	return mapEntArrToProcessingFileArr(fps), pagination, nil
}

func mapEntArrToProcessingFileArr(arr []*ent.ProcessingFile) []*ProcessingFile {
	var result []*ProcessingFile
	for _, v := range arr {
		result = append(result, &ProcessingFile{*v})
	}
	return result
}

func SaveAll(ctx context.Context, client *ent.Client, fpArr []ProcessingFile, needResult bool) ([]ProcessingFile, error) {
	// 1. Build bulk
	bulk := make([]*ent.ProcessingFileCreate, len(fpArr))
	for i, fp := range fpArr {
		bulk[i] = mapProcessingFile(client, fp)
	}

	// 2. Create by bulk
	fpSavedArr, err := client.ProcessingFile.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save file award point to DB")
	}

	// 3. Check if you NOT need result => return empty
	if !needResult {
		return []ProcessingFile{}, nil
	}

	// 4. Map Result & return
	var res []ProcessingFile
	for _, fpSaved := range fpSavedArr {
		model := &ProcessingFile{*fpSaved}
		res = append(res, *model)
	}
	return res, nil
}
