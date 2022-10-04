package common

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"time"
)

var fileAwardPoints = []fileawardpoint.FileAwardPoint{
	{
		Id:                1,
		MerchantId:        1,
		DisplayName:       "import_file.xlsx",
		FileUrl:           "https://a.com",
		ResultFileUrl:     "https://a.com",
		Status:            0,
		StatsTotalRow:     100,
		StatsTotalSuccess: 98,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         "quy.tm@teko.vn",
		UpdatedBy:         "quy.tm@teko.vn",
	},
	{
		Id:                2,
		MerchantId:        1,
		DisplayName:       "import_file_2.xlsx",
		FileUrl:           "https://a.com",
		ResultFileUrl:     "https://a.com",
		Status:            0,
		StatsTotalRow:     100,
		StatsTotalSuccess: 98,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CreatedBy:         "quy.tm@teko.vn",
		UpdatedBy:         "quy.tm@teko.vn",
	},
}
