package fpRowGroup

import (
	"context"
	"fmt"

	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerowgroup"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		FindByFileIdAndStatusIn(context.Context, int64, []int16) ([]*ProcessingFileRowGroup, error)

		SaveAll(context.Context, []ProcessingFileRowGroup, bool) ([]ProcessingFileRowGroup, error)
		UpdateByJob(context.Context, int, string, string, int16, string, int64) (*ProcessingFileRowGroup, error)

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

func (r *repoImpl) FindByFileIdAndStatusIn(ctx context.Context, fileID int64, statuses []int16) ([]*ProcessingFileRowGroup, error) {
	rowGroups, err := r.client.ProcessingFileRowGroup.
		Query().
		Where(
			processingfilerowgroup.FileID(fileID),
			processingfilerowgroup.StatusIn(statuses...),
		).
		All(ctx)

	if err != nil {
		logger.Errorf("fail to get %s by filedID=%d, err = %v", Name(), fileID, err)
		return nil, fmt.Errorf("fail to get %s by status", Name())
	}

	return mapEntArrToProcessingFileRowGroupArr(rowGroups), nil
}

func (r *repoImpl) SaveAll(ctx context.Context, pfrArr []ProcessingFileRowGroup, needResult bool) ([]ProcessingFileRowGroup, error) {
	return SaveAll(ctx, r.client, pfrArr, needResult)
}

func (r *repoImpl) UpdateByJob(ctx context.Context, id int, requestCurl string, responseRaw string,
	status int16, errorDisplay string, executedTime int64) (*ProcessingFileRowGroup, error) {
	fpr, err := r.client.ProcessingFileRowGroup.UpdateOneID(id).
		SetStatus(status).
		SetGroupRequestCurl(requestCurl).
		SetGroupResponseRaw(responseRaw).
		SetErrorDisplay(errorDisplay).
		SetExecutedTime(executedTime).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &ProcessingFileRowGroup{
		ProcessingFileRowGroup: *fpr,
	}, nil
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

func mapEntArrToProcessingFileRowGroupArr(arr ent.ProcessingFileRowGroups) []*ProcessingFileRowGroup {
	var result []*ProcessingFileRowGroup
	for _, v := range arr {
		result = append(result, &ProcessingFileRowGroup{*v})
	}
	return result
}

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
