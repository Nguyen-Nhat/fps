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
		Status:            fap.Status,
		StatsTotalRow:     fap.StatsTotalRow,
		StatsTotalSuccess: fap.StatsTotalSuccess,
		CreatedAt:         fap.CreatedAt,
		CreatedBy:         fap.CreatedBy,
	}
}

func toGetListResponseByEntity(fap *fileawardpoint.FileAwardPoint) FileAwardPoint {
	return FileAwardPoint{
		MerchantId:        fap.MerchantID,
		FileAwardPointId:  fap.ID,
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
