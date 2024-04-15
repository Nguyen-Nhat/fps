package fileprocessingrow

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	commonQuery "git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/query"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

type (
	Repo interface {
		FindRowsByFileIdForJobExecute(context.Context, int, int) ([]*ProcessingFileRow, error)
		FindByFileIdAndTaskIndexAndGroupValueAndStatus(context.Context, int64, int32, string, []int16) ([]*ProcessingFileRow, error)
		FindRowIdsByFileIdAndFilter(context.Context, int64, GetListFileRowsRequest) ([]int, int, error)
		FindRowsByIDsAndOffsetLimit(context.Context, int64, []int) ([]*ProcessingFileRow, error)
		FindByID(ctx context.Context, id int) (*ProcessingFileRow, error)
		FindByFileIDAndTaskIndexesAndResultAsyncNotEmpty(ctx context.Context, fileID int64, taskIndexes ...int32) ([]ResultAsyncDAO, error)

		Save(context.Context, ProcessingFileRow) (*ProcessingFileRow, error)
		SaveAll(context.Context, []ProcessingFileRow, bool) ([]ProcessingFileRow, error)
		UpdateByJob(context.Context, int, string, string, string, int16, string, int64) (*ProcessingFileRow, error)
		UpdateByJobForListIDs(context.Context, []int, string, int16, string, int64) error
		UpdateStatusFromTask(context.Context, int64, int32, int32) (int, error)
		Update(ctx context.Context, data *ProcessingFileRow) error
		ForceTimeout(ctx context.Context, fileId int64) error

		DeleteByFileId(context.Context, int64) error

		// Custom query ------

		Statistics(context.Context, int64) ([]CustomStatisticModel, error)
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
	fileID int64, taskIndex int32, groupValue string, statuses []int16,
) ([]*ProcessingFileRow, error) {
	pfrs, err := r.client.ProcessingFileRow.
		Query().
		Where(
			processingfilerow.FileID(fileID),
			processingfilerow.TaskIndex(taskIndex),
			processingfilerow.GroupByValue(groupValue),
			processingfilerow.StatusIn(statuses...),
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

func (r *repoImpl) FindRowIdsByFileIdAndFilter(ctx context.Context, fileID int64, req GetListFileRowsRequest) ([]int, int, error) {
	query := r.client.ProcessingFileRow.Query().
		Where(processingfilerow.FileID(fileID))

	// query by filter
	// todo ...

	allRowIDs, err := query.Unique(true).Select(processingfilerow.FieldRowIndex).Ints(ctx)
	if err != nil {
		return nil, 0, err
	} else if len(allRowIDs) == 0 {
		return []int{}, 0, nil
	}

	total := len(allRowIDs)

	rowIDs, err := query.
		Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		Order(ent.Asc(processingfilerow.FieldRowIndex)).
		GroupBy(processingfilerow.FieldRowIndex).
		Ints(ctx)
	if err != nil {
		return nil, 0, err
	}

	return rowIDs, total, nil
}

func (r *repoImpl) FindRowsByIDsAndOffsetLimit(ctx context.Context, fileID int64, rowIDs []int) ([]*ProcessingFileRow, error) {
	rowTasks, err := r.client.ProcessingFileRow.Query().
		Where(
			processingfilerow.FileID(fileID),
			processingfilerow.RowIndexIn(converter.IntArrToInt32Arr(rowIDs)...),
		).
		Order(
			ent.Asc(processingfilerow.FieldRowIndex),
			ent.Asc(processingfilerow.FieldTaskIndex),
		).
		All(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return mapEntArrToProcessingFileArr(rowTasks), nil
}

func (r *repoImpl) FindByFileIDAndTaskIndexesAndResultAsyncNotEmpty(ctx context.Context,
	fileID int64, taskIndexes ...int32) ([]ResultAsyncDAO, error) {

	query := `
			SELECT row_index, task_index, result_async
			FROM processing_file_row 
			WHERE file_id = ? AND task_index IN (?) AND result_async <> ''
			ORDER BY row_index, task_index;
	`

	convertFunc := func(rows *dbsql.Rows, dao *ResultAsyncDAO) error {
		return rows.Scan(&dao.RowIndex, &dao.TaskIndex, &dao.ResultAsync)
	}

	taskIndexesStr := strings.Join(converter.Map(taskIndexes, func(t int32) string {
		return fmt.Sprintf("%d", t)
	}), ",")

	return commonQuery.RunRawQuery(ctx, r.sqlDB, query, convertFunc, fileID, taskIndexesStr)
}

func (r *repoImpl) UpdateByJob(ctx context.Context, id int,
	taskMapping string,
	requestCurl string, responseRaw string,
	status int16, errorDisplay string, executedTime int64) (*ProcessingFileRow, error) {
	fpr, err := r.client.ProcessingFileRow.UpdateOneID(id).
		SetStatus(status).
		SetTaskMapping(taskMapping).
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

func (r *repoImpl) UpdateStatusFromTask(ctx context.Context, fileID int64, rowIndex int32, fromTaskIndex int32) (int, error) {
	return r.client.ProcessingFileRow.Update().
		SetStatus(StatusRejected).
		Where(
			processingfilerow.FileID(fileID),
			processingfilerow.RowIndex(rowIndex),
			processingfilerow.TaskIndexGT(fromTaskIndex),
		).
		Save(ctx)
}

func (r *repoImpl) ForceTimeout(ctx context.Context, fileId int64) error {
	_, err := r.client.ProcessingFileRow.Update().
		Where(processingfilerow.FileID(fileId)).
		Where(processingfilerow.StatusNotIn(StatusSuccess, StatusFailed, StatusRejected)).
		SetStatus(StatusFailed).
		SetErrorDisplay(constant.Timeout).
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

func (r *repoImpl) Statistics(ctx context.Context, fileID int64) ([]CustomStatisticModel, error) {
	rawQuery := `
			SELECT 
				row_index, GROUP_CONCAT(status), COUNT(*), IFNULL(GROUP_CONCAT(IF(error_display='', null, error_display)),'') 
			FROM processing_file_row 
			WHERE file_id = ? 
			GROUP BY row_index
		`

	convertFunc := func(rows *dbsql.Rows, dao *CustomStatisticModel) error {
		return rows.Scan(&dao.RowIndex, &dao.Statuses, &dao.Count, &dao.ErrorDisplays)
	}

	return commonQuery.RunRawQuery(ctx, r.sqlDB, rawQuery, convertFunc, fileID)
}

func (r *repoImpl) Update(ctx context.Context, data *ProcessingFileRow) error {
	query := r.client.ProcessingFileRow.Update()

	query = query.Where(processingfilerow.ID(data.ID))

	if data.ReceiveResultAsyncAt != nil {
		query.Where(processingfilerow.Or(
			processingfilerow.ReceiveResultAsyncAtIsNil(),
			processingfilerow.ReceiveResultAsyncAtLT(*data.ReceiveResultAsyncAt),
		))
	}

	if data.ResultAsync != nil {
		query = query.SetNillableResultAsync(data.ResultAsync)
	}
	if data.ReceiveResultAsyncAt != nil {
		query = query.SetNillableReceiveResultAsyncAt(data.ReceiveResultAsyncAt)
	}
	if data.Status != 0 {
		query = query.SetStatus(data.Status)
	}
	return query.Exec(ctx)
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
