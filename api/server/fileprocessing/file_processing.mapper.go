package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"github.com/tidwall/gjson"
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
	totalErrorPlus := getTotalErrorPlus(pf.FileParameters)

	return ProcessingHistoryFile{
		ClientId:            pf.ClientID,
		ProcessingFileId:    pf.ID,
		FileDisplayName:     pf.DisplayName,
		FileUrl:             pf.FileURL,
		ResultFileUrl:       pf.ResultFileURL,
		Status:              mapStatus(pf.Status),
		SellerID:            pf.SellerID,
		StatsTotalRow:       pf.StatsTotalRow + totalErrorPlus,
		StatsTotalProcessed: pf.StatsTotalProcessed + totalErrorPlus,
		StatsTotalSuccess:   pf.StatsTotalSuccess,
		ErrorDisplay:        pf.ErrorDisplay,
		CreatedAt:           pf.CreatedAt.UnixMilli(),
		CreatedBy:           pf.CreatedBy,
	}
}

// totalErrorPlusField ... is created when
// - BFF pre-processed file and there are some rows are error
// - BFF want to give the number of error rows to FPS to store
// - When FPS return Import History, FPS will add this number to result statistic, include fields: statsTotalRow and statsTotalProcessed
// Related ticket: https://jira.teko.vn/browse/SRE-3428
const totalErrorPlusField = "totalErrorPlus"

// getTotalErrorPlus ... use totalErrorPlusField
func getTotalErrorPlus(fileParameters string) int32 {
	if len(fileParameters) == 0 {
		return 0
	}

	res := gjson.Get(fileParameters, totalErrorPlusField)
	totalErrorPlus := int32(res.Int())
	if totalErrorPlus <= 0 { // ignore negative value
		return 0
	} else {
		return totalErrorPlus
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
