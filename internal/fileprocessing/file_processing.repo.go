package fileprocessing

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error)
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
