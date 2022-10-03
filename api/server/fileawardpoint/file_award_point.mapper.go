package fileawardpoint

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"

func toFapDetailResponseByEntity(fap *fileawardpoint.FileAwardPoint) *GetFileAwardPointDetailResponse {
	return &GetFileAwardPointDetailResponse{
		Id:                fap.Id,
		MerchantId:        fap.MerchantId,
		DisplayName:       fap.DisplayName,
		FileUrl:           fap.ResultFileUrl,
		ResultFileUrl:     fap.ResultFileUrl,
		Status:            fap.Status,
		StatsTotalRow:     fap.StatsTotalRow,
		StatsTotalSuccess: fap.StatsTotalSuccess,
		CreatedAt:         fap.CreatedAt,
		CreatedBy:         fap.CreatedBy,
	}
}
