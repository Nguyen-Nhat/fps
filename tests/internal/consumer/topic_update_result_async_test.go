package consumer

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/consumer/topics"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
	"go.tekoapis.com/tekone/app/aggregator/file_processing_service/api"
	"go.tekoapis.com/tekone/library/teka"
)

type topicUpdateResultAsyncSuite struct {
	suite.Suite
	worker    *topics.Worker
	ctx       context.Context
	db        *sql.DB
	entClient *ent.Client
	fixedTime time.Time
}

func TestTopicUpdateResultAsync(t *testing.T) {
	ts := &topicUpdateResultAsyncSuite{}
	ts.ctx = context.Background()
	ts.fixedTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	ts.db, ts.entClient = common.PrepareDatabaseSqlite(ts.ctx, ts.T(), true)
	ts.tearDown()
	ts.worker = topics.NewWorker(topics.WorkerAdjust{
		FileProcessingRowRepository: fileprocessingrow.NewRepo(ts.db),
	})

	defer func() {
		ts.tearDown()
	}()

	suite.Run(t, ts)
}

func (ts *topicUpdateResultAsyncSuite) tearDown() {
	common.TruncateAllTables(ts.ctx, ts.db, ts.entClient)
}

func (ts *topicUpdateResultAsyncSuite) assert(payload api.TaskResultEvent, wantFile string) {
	defer func() {
		ts.tearDown()
	}()

	data, err := payload.Marshal()
	assert.Nil(ts.T(), err)

	err = ts.worker.TopicUpsertResultAsync(context.Background(), &teka.Message{
		ConsumerMessage: &sarama.ConsumerMessage{
			Value:     data,
			Timestamp: ts.fixedTime.Add(time.Second * 10),
		},
	})
	assert.Nil(ts.T(), err)

	processingFileRows, err := ts.entClient.ProcessingFileRow.
		Query().
		All(ts.ctx)
	assert.Nil(ts.T(), err)

	for idx := range processingFileRows {
		processingFileRows[idx].UpdatedAt = ts.fixedTime
	}

	goldie.New(ts.T()).AssertJson(ts.T(), wantFile, processingFileRows)
}

func (ts *topicUpdateResultAsyncSuite) Test200_HappyCase_ThenReturnSuccess() {
	common.MockProcessingFileRow(ts.ctx, ts.entClient, []fileprocessingrow.ProcessingFileRow{
		{ProcessingFileRow: ent.ProcessingFileRow{
			FileID:    2,
			RowIndex:  1,
			TaskIndex: 2,
			Status:    fileprocessingrow.StatusWaitForAsync,
			CreatedAt: ts.fixedTime,
			UpdatedAt: ts.fixedTime,
		}},
	})
	payload := api.TaskResultEvent{
		RefId: 1,
		Code:  0,
		ResultData: &api.ResultData{
			ProcessingResult: "result",
			Message:          "",
			FinishedAt:       "2021-01-01T00:00:00Z",
		},
	}
	ts.assert(payload, "happy_case_update_success")
}

func (ts *topicUpdateResultAsyncSuite) TestUpdateWithMessageError_ThenReturnSuccess() {
	common.MockProcessingFileRow(ts.ctx, ts.entClient, []fileprocessingrow.ProcessingFileRow{
		{ProcessingFileRow: ent.ProcessingFileRow{
			FileID:       2,
			RowIndex:     1,
			TaskIndex:    2,
			Status:       fileprocessingrow.StatusWaitForAsync,
			ErrorDisplay: "",
			ExecutedTime: 0,
			CreatedAt:    ts.fixedTime,
			UpdatedAt:    ts.fixedTime,
		}},
	})

	payload := api.TaskResultEvent{
		RefId: 1,
		Code:  1,
		ResultData: &api.ResultData{
			ProcessingResult: "result",
			Message:          "error exist",
			FinishedAt:       "2021-01-01T00:00:00Z",
		},
	}
	ts.assert(payload, "update_with_message_error")
}

func (ts *topicUpdateResultAsyncSuite) TestAlreadyUpdated_ThenReturnSuccess() {
	updatedDataStr := "updatedData"
	common.MockProcessingFileRow(ts.ctx, ts.entClient, []fileprocessingrow.ProcessingFileRow{
		{ProcessingFileRow: ent.ProcessingFileRow{
			FileID:               2,
			RowIndex:             1,
			TaskIndex:            1,
			Status:               fileprocessingrow.StatusSuccess,
			ResultAsync:          &updatedDataStr,
			ReceiveResultAsyncAt: &ts.fixedTime,
			CreatedAt:            ts.fixedTime,
			UpdatedAt:            ts.fixedTime,
		}},
		{ProcessingFileRow: ent.ProcessingFileRow{
			FileID:    2,
			RowIndex:  1,
			TaskIndex: 2,
			Status:    fileprocessingrow.StatusSuccess,
			CreatedAt: ts.fixedTime,
			UpdatedAt: ts.fixedTime,
		}},
	})

	payload := api.TaskResultEvent{
		RefId: 1,
		Code:  0,
		ResultData: &api.ResultData{
			ProcessingResult: "result",
			Message:          "",
			FinishedAt:       "2021-01-01T00:00:00Z",
		},
	}
	ts.assert(payload, "already_updated")
}
