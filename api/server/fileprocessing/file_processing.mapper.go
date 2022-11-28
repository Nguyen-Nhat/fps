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

func toProcessHistoryFileFromEntity(fap *fileprocessing.ProcessingFile) ProcessingHistoryFile {
	return ProcessingHistoryFile{
		ClientId:          fap.ClientID,
		ProcessingFileId:  fap.ID,
		FileDisplayName:   fap.DisplayName,
		FileUrl:           fap.FileURL,
		ResultFileUrl:     fap.ResultFileURL,
		Status:            fap.Status,
		StatsTotalRow:     fap.StatsTotalRow,
		StatsTotalSuccess: fap.StatsTotalSuccess,
		CreatedAt:         fap.CreatedAt,
		CreatedBy:         fap.CreatedBy,
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
