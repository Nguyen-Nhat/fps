package fileawardpoint

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"

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
