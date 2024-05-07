package topics

import (
	"context"
	"encoding/json"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"go.tekoapis.com/tekone/app/aggregator/file_processing_service/api"
	"go.tekoapis.com/tekone/library/teka"
)

const payloadCodeSuccess = 0

func (w *Worker) TopicUpsertResultAsync(ctx context.Context, msg *teka.Message) error {
	if msg == nil {
		logger.Info("TopicUpsertResultAsync | topic message null")
		return nil
	}
	var payload api.TaskResultEvent
	if err := payload.Unmarshal(msg.Value); err != nil {
		logger.Errorf("TopicUpsertResultAsync | Marshaling topic payload with msg %v got error value %v", msg, err)
		return err
	}

	logger.Infof("TopicUpsertResultAsync | topic payload %v", payload)
	task, err := w.fileProcessingRowRepository.FindByID(ctx, int(payload.RefId))
	if err != nil {
		logger.Errorf("TopicUpsertResultAsync | FindByID got error value %v", err)
		return err
	}
	if task == nil {
		logger.Errorf("TopicUpsertResultAsync | Task not found with id %v", payload.RefId)
		return nil
	}

	resultAsync, err := json.Marshal(payload.ResultData)
	if err != nil {
		logger.Errorf("TopicUpsertResultAsync | Marshaling result data got error value %v", err)
		return err
	}
	result := string(resultAsync)
	task.ResultAsync = &result
	if payload.Code == payloadCodeSuccess {
		task.Status = fileprocessingrow.StatusSuccess
	} else {
		task.Status = fileprocessingrow.StatusFailed
	}
	task.ReceiveResultAsyncAt = &msg.Timestamp

	err = w.fileProcessingRowRepository.Update(ctx, task)
	if err != nil {
		logger.Errorf("TopicUpsertResultAsync | Update result async got error value %v", err)
		return err
	}
	return nil
}
