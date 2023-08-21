package fileprocessingrow

import (
	"context"
	"errors"
	"fmt"

	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		FindRowsByFileIdForJobExecute(context.Context, int, int) ([]*ProcessingFileRow, error)
		FindByFileIdAndTaskIndexAndGroupValueAndStatus(context.Context, int64, int32, string, int16) ([]*ProcessingFileRow, error)

		Save(context.Context, ProcessingFileRow) (*ProcessingFileRow, error)
		SaveAll(context.Context, []ProcessingFileRow, bool) ([]ProcessingFileRow, error)
		UpdateByJob(context.Context, int, string, string, int16, string, int64) (*ProcessingFileRow, error)
		UpdateByJobForListIDs(context.Context, []int, string, int16, string, int64) error

		DeleteByFileId(context.Context, int64) error

		// Custom query ------

		Statistics(int64) ([]CustomStatisticModel, error)
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

// Save ... Implementation function ------------------------------------------------------------------------------------
func (r *repoImpl) Save(ctx context.Context, fp ProcessingFileRow) (*ProcessingFileRow, error) {
	return save(ctx, r.client, fp)
}

func (r *repoImpl) SaveAll(ctx context.Context, pfrArr []ProcessingFileRow, needResult bool) ([]ProcessingFileRow, error) {
	return SaveAll(ctx, r.client, pfrArr, needResult)
}

func (r *repoImpl) FindByID(ctx context.Context, id int) (*ProcessingFileRow, error) {
	fp, err := r.client.ProcessingFileRow.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular %v: %w", Name(), err)
	}

	return &ProcessingFileRow{*fp}, nil
}

func (r *repoImpl) FindByFileIdAndTaskIndexAndGroupValueAndStatus(ctx context.Context,
	fileID int64, taskIndex int32, groupValue string, status int16,
) ([]*ProcessingFileRow, error) {
	pfrs, err := r.client.ProcessingFileRow.
		Query().
		Where(
			processingfilerow.FileID(fileID),
			processingfilerow.TaskIndex(taskIndex),
			processingfilerow.GroupByValue(groupValue),
			processingfilerow.Status(status),
		).
		All(ctx)

	if err != nil {
		logger.Errorf("fail to get %s by filedID=%d, groupValue=%s, err = %v", Name(), fileID, groupValue, err)
		return nil, fmt.Errorf("fail to get %s by groupValue", Name())
	}

	return mapEntArrToProcessingFileArr(pfrs), nil
}

func (r *repoImpl) FindRowsByFileIdForJobExecute(ctx context.Context, fileId int, limit int) ([]*ProcessingFileRow, error) {
	/*
		select *
		from processing_file_row
		where file_id = ?
			and row_index in (
				select distinct row_index from processing_file_row where file_id=? and status=1
			)
			and row_index not in (
				select distinct row_index from processing_file_row where file_id=? and status=3
			);
	*/

	pfs, err := r.client.ProcessingFileRow.Query().
		Where(func(s *sql.Selector) {
			t := sql.Table(processingfilerow.Table)
			eqFileId := sql.EQ(t.C(processingfilerow.FieldFileID), fileId)
			eqStatus := sql.EQ(t.C(processingfilerow.FieldStatus), StatusInit)
			eqExcludeStatus := sql.EQ(t.C(processingfilerow.FieldStatus), StatusFailed)

			// file_id=?
			s.Where(eqFileId)
			// row_id in (select distinct row_id from processing_file_row where file_id=? and status=?)
			s.Where(sql.In(
				s.C(processingfilerow.FieldRowIndex),
				sql.Select(sql.Distinct(t.C(processingfilerow.FieldRowIndex))).
					From(t).
					Where(sql.And(eqFileId, eqStatus)),
			))
			// row_id NOT in (select distinct row_id from processing_file_row where file_id=? and status=?)
			s.Where(sql.NotIn(
				s.C(processingfilerow.FieldRowIndex),
				sql.Select(sql.Distinct(t.C(processingfilerow.FieldRowIndex))).
					From(t).
					Where(sql.And(eqFileId, eqExcludeStatus)),
			))

			s.OrderBy(processingfilerow.FieldRowIndex, processingfilerow.FieldTaskIndex)
			s.Limit(limit)
		}).All(ctx)

	if err != nil {
		logger.Errorf("fail to get %v by status with status %#v", Name(), StatusInit)
		return nil, errors.New("fail to get file processing row by status")
	}

	return mapEntArrToProcessingFileArr(pfs), nil
}

func (r *repoImpl) UpdateByJob(ctx context.Context, id int, requestCurl string, responseRaw string,
	status int16, errorDisplay string, executedTime int64) (*ProcessingFileRow, error) {
	fpr, err := r.client.ProcessingFileRow.UpdateOneID(id).
		SetStatus(status).
		SetTaskRequestCurl(requestCurl).
		SetTaskRequestRaw("").
		SetTaskResponseRaw(responseRaw).
		SetErrorDisplay(errorDisplay).
		SetExecutedTime(executedTime).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &ProcessingFileRow{
		ProcessingFileRow: *fpr,
	}, nil
}

func (r *repoImpl) UpdateByJobForListIDs(ctx context.Context, ids []int, responseRaw string,
	status int16, errorDisplay string, executedTime int64) error {
	_, err := r.client.ProcessingFileRow.Update().
		SetStatus(status).
		SetTaskRequestRaw("").
		SetTaskResponseRaw(responseRaw).
		SetErrorDisplay(errorDisplay).
		SetExecutedTime(executedTime).
		Where(processingfilerow.IDIn(ids...)).
		Save(ctx)

	return err
}

func (r *repoImpl) DeleteByFileId(ctx context.Context, fileId int64) error {
	deletedRowCount, err := r.client.ProcessingFileRow.Delete().Where(processingfilerow.FileID(fileId)).Exec(ctx)
	if err != nil {
		logger.Errorf("Cannot delete records which have fileId=%v, got: %v", fileId, err)
		return err
	} else {
		logger.Warnf("Deleted %v records which have fileId=%v", deletedRowCount, fileId)
		return nil
	}
}

func (r *repoImpl) Statistics(fileId int64) ([]CustomStatisticModel, error) {
	rows, err := r.sqlDB.Query(
		`
			SELECT 
				row_index, GROUP_CONCAT(status), COUNT(*), GROUP_CONCAT(error_display) 
			FROM processing_file_row 
			WHERE file_id = ? 
			GROUP BY row_index
		`,
		fileId)

	if err != nil {
		return []CustomStatisticModel{}, err
	}

	defer rows.Close()

	var statistics []CustomStatisticModel
	for rows.Next() {
		var stats CustomStatisticModel
		if err := rows.Scan(&stats.RowIndex, &stats.Statuses, &stats.Count, &stats.ErrorDisplays); err != nil {
			return []CustomStatisticModel{}, err
		}
		statistics = append(statistics, stats)
	}
	if err = rows.Err(); err != nil {
		return []CustomStatisticModel{}, err
	}

	return statistics, nil
}

// private function ---------------------------------------------------------------------------------------------

func save(ctx context.Context, client *ent.Client, fp ProcessingFileRow) (*ProcessingFileRow, error) {
	// 1. Create
	fpSaved, err := mapProcessingFileRow(client, fp).Save(ctx)

	if err != nil {
		logger.Errorf("Cannot save to file processing, got: %v", err)
		return nil, fmt.Errorf("failed to save file processing to DB")
	}

	// 2. Return
	return &ProcessingFileRow{*fpSaved}, nil
}

func mapProcessingFileRow(client *ent.Client, fpr ProcessingFileRow) *ent.ProcessingFileRowCreate {
	return client.ProcessingFileRow.Create().
		SetFileID(fpr.FileID).
		SetRowIndex(fpr.RowIndex).
		SetRowDataRaw(fpr.RowDataRaw).
		SetTaskIndex(fpr.TaskIndex).
		SetTaskMapping(fpr.TaskMapping).
		SetTaskDependsOn(fpr.TaskDependsOn).
		SetTaskRequestCurl(fpr.TaskRequestCurl).
		SetTaskRequestRaw(fpr.TaskRequestRaw).
		SetTaskResponseRaw(fpr.TaskResponseRaw).
		SetGroupByValue(fpr.GroupByValue).
		SetStatus(fpr.Status).
		SetErrorDisplay(fpr.ErrorDisplay).
		SetExecutedTime(fpr.ExecutedTime).
		SetCreatedAt(fpr.CreatedAt).
		SetUpdatedAt(fpr.UpdatedAt)
}

func mapEntArrToProcessingFileArr(arr ent.ProcessingFileRows) []*ProcessingFileRow {
	var result []*ProcessingFileRow
	for _, v := range arr {
		result = append(result, &ProcessingFileRow{*v})
	}
	return result
}

func SaveAll(ctx context.Context, client *ent.Client, pfrArr []ProcessingFileRow, needResult bool) ([]ProcessingFileRow, error) {
	// 1. Build bulk
	bulk := make([]*ent.ProcessingFileRowCreate, len(pfrArr))
	for i, fap := range pfrArr {
		bulk[i] = mapProcessingFileRow(client, fap)
	}

	// 2. Create by bulk
	fapSavedArr, err := client.ProcessingFileRow.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save %v to DB, got err: %v", Name(), err)
	}

	// 3. Check if you NOT need result => return empty
	if !needResult {
		return []ProcessingFileRow{}, nil
	}

	// 4. Map Result & return
	var res []ProcessingFileRow
	for _, fapSaved := range fapSavedArr {
		model := &ProcessingFileRow{*fapSaved}
		res = append(res, *model)
	}
	return res, nil
}
