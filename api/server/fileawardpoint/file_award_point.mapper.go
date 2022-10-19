package fileawardpoint

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
)

func toFapDetailResponseByEntity(fap *fileawardpoint.FileAwardPoint) *GetFileAwardPointDetailResponse {
	return &GetFileAwardPointDetailResponse{
		Id:                fap.ID,
		MerchantId:        fap.MerchantID,
		DisplayName:       fap.DisplayName,
		FileUrl:           fap.FileURL,
		ResultFileUrl:     fap.ResultFileURL,
		Status:            mapStatus(fap.Status),
		StatsTotalRow:     fap.StatsTotalRow,
		StatsTotalSuccess: fap.StatsTotalSuccess,
		CreatedAt:         fap.CreatedAt.UnixMilli(),
		CreatedBy:         fap.CreatedBy,
	}
}

func toFapCreateResponseByEntity(fap *fileawardpoint.CreateFileAwardPointResDTO) *CreateFileAwardPointDetailResponse {
	return &CreateFileAwardPointDetailResponse{
		FileAwardPointID: int(fap.FileAwardPointId),
	}
}

func toGetListResponseByEntity(fap *fileawardpoint.FileAwardPoint) FileAwardPoint {
	return FileAwardPoint{
		MerchantId:        fap.MerchantID,
		FileAwardPointId:  fap.ID,
		FileDisplayName:   fap.DisplayName,
		FileUrl:           fap.FileURL,
		ResultFileUrl:     fap.ResultFileURL,
		Status:            mapStatus(fap.Status),
		StatsTotalRow:     fap.StatsTotalRow,
		StatsTotalSuccess: fap.StatsTotalSuccess,
		CreatedAt:         fap.CreatedAt.UnixMilli(),
		CreatedBy:         fap.CreatedBy,
	}
}

func fromFileArrToGetListData(fap []*fileawardpoint.FileAwardPoint, pagination *response.Pagination) *GetListFileAwardPointData {
	var result []FileAwardPoint
	for _, v := range fap {
		result = append(result, toGetListResponseByEntity(v))
	}

	return &GetListFileAwardPointData{
		FileAwardPoints: result,
		Pagination:      *pagination,
	}
}

func mapStatus(statusInDB int16) string {
	switch statusInDB {
	case fileawardpoint.StatusInit:
		return FapStatusInit
	case fileawardpoint.StatusProcessing:
		return FapStatusProcessing
	case fileawardpoint.StatusFailed:
		return FapStatusFailed
	case fileawardpoint.StatusFinished:
		return FapStatusFinished
	default:
		return ""
	}
}
