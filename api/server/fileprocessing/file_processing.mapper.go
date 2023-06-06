package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
)

func fromInternalToGetFileHistoryData(fp []*fileprocessing.ProcessingFile, pagination *response.Pagination) *GetFileProcessHistoryData {
	result := make([]ProcessingHistoryFile, 0)
	for _, v := range fp {
		result = append(result, toProcessHistoryFileFromEntity(v))
	}

	return &GetFileProcessHistoryData{
		ProcessingFiles: result,
		Pagination:      *pagination,
	}
}

func toProcessHistoryFileFromEntity(pf *fileprocessing.ProcessingFile) ProcessingHistoryFile {
	return ProcessingHistoryFile{
		ClientId:            pf.ClientID,
		ProcessingFileId:    pf.ID,
		FileDisplayName:     pf.DisplayName,
		FileUrl:             pf.FileURL,
		ResultFileUrl:       pf.ResultFileURL,
		Status:              mapStatus(pf.Status),
		StatsTotalRow:       pf.StatsTotalRow,
		StatsTotalProcessed: pf.StatsTotalProcessed,
		StatsTotalSuccess:   pf.StatsTotalSuccess,
		ErrorDisplay:        pf.ErrorDisplay,
		CreatedAt:           pf.CreatedAt.UnixMilli(),
		CreatedBy:           pf.CreatedBy,
	}
}

func mapStatus(statusInDB int16) string {
	switch statusInDB {
	case fileprocessing.StatusInit:
		return FpStatusInit
	case fileprocessing.StatusProcessing:
		return FpStatusProcessing
	case fileprocessing.StatusFailed:
		return FpStatusFailed
	case fileprocessing.StatusFinished:
		return FpStatusFinished
	default:
		return ""
	}
}
