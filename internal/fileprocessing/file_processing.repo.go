package fileprocessing

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfile"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error)
		FindByClientIdAndPagination(context.Context, *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error)
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
func (r *repoImpl) Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error) {
	return save(ctx, r.client, fp)
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
		SetCreatedBy(fp.CreatedBy)
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r *repoImpl) FindByClientIdAndPagination(ctx context.Context, dto *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error) {
	query := r.client.ProcessingFile.Query().Where(processingfile.ClientID(dto.ClientId))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count from db while querying all processing file")
	}
	pagination := response.GetPagination(total, dto.Page, dto.Size)

	fps, err := query.Limit(dto.Size).Offset((dto.Page - 1) * dto.Size).Order(ent.Desc(processingfile.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed querying all processing file: %v", err)
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
