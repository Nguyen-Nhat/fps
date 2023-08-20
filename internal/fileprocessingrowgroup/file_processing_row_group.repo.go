package fpRowGroup

import (
	"context"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"

	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

type (
	Repo interface {
		SaveAll(context.Context, []ProcessingFileRowGroup, bool) ([]ProcessingFileRowGroup, error)
		DeleteByFileId(context.Context, int64) error
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
	//drvDebug := dialect.Debug(drv)
	//client := ent.NewClient(ent.Driver(drvDebug))
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client, sqlDB: db}
}

func (r *repoImpl) SaveAll(ctx context.Context, pfrArr []ProcessingFileRowGroup, needResult bool) ([]ProcessingFileRowGroup, error) {
	return SaveAll(ctx, r.client, pfrArr, needResult)
}

func (r *repoImpl) DeleteByFileId(ctx context.Context, fileId int64) error {
	deletedRowCount, err := r.client.ProcessingFileRowGroup.Delete().Where(processingfilerowgroup.FileID(fileId)).Exec(ctx)
	if err != nil {
		logger.Errorf("Cannot delete records which have fileId=%v, got: %v", fileId, err)
		return err
	} else {
		logger.Warnf("Deleted %v records which have fileId=%v", deletedRowCount, fileId)
		return nil
	}
}

// private function ---------------------------------------------------------------------------------------------

func SaveAll(ctx context.Context, client *ent.Client, pfrgArr []ProcessingFileRowGroup, needResult bool) ([]ProcessingFileRowGroup, error) {
	// 1. Build bulk
	bulk := make([]*ent.ProcessingFileRowGroupCreate, len(pfrgArr))
	for i, pfrg := range pfrgArr {
		bulk[i] = mapProcessingFileRowGroup(client, pfrg)
	}

	// 2. Create by bulk
	pfrgSavedArr, err := client.ProcessingFileRowGroup.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save %v to DB, got err: %v", Name(), err)
	}

	// 3. Check if you NOT need result => return empty
	if !needResult {
		return []ProcessingFileRowGroup{}, nil
	}

	// 4. Map Result & return
	var res []ProcessingFileRowGroup
	for _, pfrgSaved := range pfrgSavedArr {
		model := &ProcessingFileRowGroup{*pfrgSaved}
		res = append(res, *model)
	}
	return res, nil
}

func mapProcessingFileRowGroup(client *ent.Client, pfrg ProcessingFileRowGroup) *ent.ProcessingFileRowGroupCreate {
	return client.ProcessingFileRowGroup.Create().
		SetFileID(pfrg.FileID).
		SetTaskIndex(pfrg.TaskIndex).
		SetGroupByValue(pfrg.GroupByValue).
		SetTotalRows(pfrg.TotalRows).
		SetRowIndexList(pfrg.RowIndexList).
		SetGroupRequestCurl(pfrg.GroupRequestCurl).
		SetGroupResponseRaw(pfrg.GroupResponseRaw).
		SetStatus(pfrg.Status).
		SetErrorDisplay(pfrg.ErrorDisplay).
		SetExecutedTime(pfrg.ExecutedTime).
		SetCreatedAt(pfrg.CreatedAt).
		SetUpdatedAt(pfrg.UpdatedAt)
}
