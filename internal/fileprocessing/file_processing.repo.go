package fileprocessing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/samber/lo"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfile"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/faltservice"
)

type (
	Repo interface {
		Save(ctx context.Context, fp ProcessingFile) (*ProcessingFile, error)
		BulkInsert(ctx context.Context, fps []ProcessingFile) error
		FindByClientIdAndPagination(context.Context, *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error)
		FindByID(context.Context, int) (*ProcessingFile, error)
		FindByStatuses(context.Context, []int16) ([]*ProcessingFile, error)
		UpdateStatusOne(context.Context, int, int16) (*ProcessingFile, error)
		UpdateStatusAndErrorDisplay(context.Context, int, int16, ErrorDisplay, *string) (*ProcessingFile, error)
		UpdateStatusAndStatsAndResultFileUrl(context.Context, int, int16, int, int, string) (*ProcessingFile, error)
		UpdateByExtractedData(ctx context.Context, id int, status int16, totalMapping int, statsTotalRow int) (*ProcessingFile, error)
		Delete(ctx context.Context, clientIds []int32) error
		PingDB(context.Context, int)
	}

	repoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

// NewRepo ...
func NewRepo(db *sql.DB) Repo {
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
		logger.Errorf("fail to get file award point by status with status %#v, err = %v", statuses, err)
		return nil, errors.New("fail to get file award point by status")
	}

	return mapEntArrToProcessingFileArr(pfs), nil
}

func (r *repoImpl) UpdateStatusOne(ctx context.Context, id int, status int16) (*ProcessingFile, error) {
	fap, err := r.client.ProcessingFile.UpdateOneID(id).SetStatus(status).Save(ctx)
	if err != nil {
		return nil, err
	}

	// Update status of ProcessingFile in f-alt-server
	go func() {
		_ = faltservice.UpdateStatusProcessingFile(fap.ID, status)
	}()

	return &ProcessingFile{
		ProcessingFile: *fap,
	}, nil
}

func (r *repoImpl) UpdateStatusAndErrorDisplay(ctx context.Context, id int, status int16, errorDisplay ErrorDisplay, resultFileURL *string) (*ProcessingFile, error) {
	updateOps := r.client.ProcessingFile.UpdateOneID(id).
		SetStatus(status).
		SetErrorDisplay(string(errorDisplay))

	if resultFileURL != nil && len(*resultFileURL) > 0 {
		updateOps.SetResultFileURL(*resultFileURL)
	}

	fap, err := updateOps.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Update status of ProcessingFile in f-alt-server
	go func() {
		_ = faltservice.UpdateStatusProcessingFile(fap.ID, status)
	}()

	return &ProcessingFile{
		ProcessingFile: *fap,
	}, nil
}

func (r *repoImpl) UpdateStatusAndStatsAndResultFileUrl(ctx context.Context, id int, status int16, totalProcessed int, totalSuccess int,
	resultFileUrl string) (*ProcessingFile, error) {
	fap, err := r.client.ProcessingFile.UpdateOneID(id).
		SetStatus(status).
		SetStatsTotalProcessed(int32(totalProcessed)).
		SetStatsTotalSuccess(int32(totalSuccess)).
		SetResultFileURL(resultFileUrl).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Update status of ProcessingFile in f-alt-server
	go func() {
		_ = faltservice.UpdateStatusProcessingFile(fap.ID, status)
	}()

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

	// Update status of ProcessingFile in f-alt-server
	go func() {
		_ = faltservice.UpdateStatusProcessingFile(pf.ID, status)
	}()

	if err != nil {
		return nil, err
	}
	return &ProcessingFile{
		ProcessingFile: *pf,
	}, nil
}

func (r *repoImpl) PingDB(ctx context.Context, id int) {
	before := time.Now()
	_, _ = r.client.ProcessingFile.Query().Where(processingfile.ID(id)).Only(ctx)
	after := time.Now() // will remove
	sub := after.Sub(before)
	subStr := ""
	if sub > 10*time.Second {
		subStr = "\t------> too long (>10s)"
	} else if sub > 5*time.Second {
		subStr = "\t---> too long (>5s)"
	}
	logger.Debugf("---------------> Ping DB ...... execute_time = %v%s", sub, subStr)
}

func (r *repoImpl) BulkInsert(ctx context.Context, fps []ProcessingFile) error {
	// 1. Build bulk
	processingFileBulk := make([]*ent.ProcessingFileCreate, len(fps))
	for i, fp := range fps {
		processingFileBulk[i] = mapProcessingFile(r.client, fp)
	}

	// 2. Create by bulk
	processingFileBulkChunk := lo.Chunk(processingFileBulk, constant.DefaultInsertBatchSize)
	for _, bulk := range processingFileBulkChunk {
		_, err := r.client.ProcessingFile.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to save %s to DB, error: %+v", Name(), err)
		}
	}

	return nil
}

func (r *repoImpl) Delete(ctx context.Context, clientIds []int32) error {
	_, err := r.client.ProcessingFile.Delete().Where(processingfile.ClientIDIn(clientIds...)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete %s from DB, error: %+v", Name(), err)
	}
	return nil
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
		SetExtFileRequest(fp.ExtFileRequest).
		SetFileURL(fp.FileURL).
		SetResultFileURL(fp.ResultFileURL).
		SetStatus(fp.Status).
		SetFileParameters(fp.FileParameters).
		SetSellerID(fp.SellerID).
		SetTotalMapping(fp.TotalMapping).
		SetStatsTotalRow(fp.StatsTotalRow).
		SetStatsTotalProcessed(fp.StatsTotalProcessed).
		SetStatsTotalSuccess(fp.StatsTotalSuccess).
		SetErrorDisplay(fp.ErrorDisplay).
		SetCreatedBy(fp.CreatedBy).
		SetCreatedAt(fp.CreatedAt).
		SetUpdatedAt(fp.UpdatedAt).
		SetMerchantID(fp.MerchantID).
		SetTenantID(fp.TenantID)
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r *repoImpl) FindByClientIdAndPagination(ctx context.Context, dto *GetFileProcessHistoryDTO) ([]*ProcessingFile, *response.Pagination, error) {
	query := r.client.ProcessingFile.Query().Where(processingfile.ClientID(dto.ClientID))

	if dto.CreatedBy != constant.EmptyString {
		query = query.Where(processingfile.CreatedBy(dto.CreatedBy))
	}
	if dto.SellerId > 0 {
		query = query.Where(processingfile.SellerID(dto.SellerId))
	}
	if len(dto.CreatedByEmails) > 0 {
		query = query.Where(processingfile.CreatedByIn(dto.CreatedByEmails...))
	}
	if len(dto.ProcessFileIds) > 0 {
		query = query.Where(processingfile.IDIn(dto.ProcessFileIds...))
	}
	if len(dto.SearchFileName) > 0 {
		query = query.Where(processingfile.DisplayNameContains(dto.SearchFileName))
	}
	if len(dto.MerchantId) > 0 {
		query = query.Where(processingfile.MerchantID(dto.MerchantId))
	}
	if len(dto.TenantId) > 0 {
		query = query.Where(processingfile.TenantID(dto.TenantId))
	}

	total, err := query.Count(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, nil, fmt.Errorf("failed to count from db while querying all processing file")
	}
	pagination := response.GetPagination(total, dto.Page, dto.PageSize)

	fps, err := query.Limit(dto.PageSize).Offset((dto.Page - 1) * dto.PageSize).Order(ent.Desc(processingfile.FieldCreatedAt)).All(ctx)
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
		return nil, fmt.Errorf("failed to save %s to DB, error: %+v", Name(), err)
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
