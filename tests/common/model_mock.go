package common

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"time"
)

var fileAwardPoints = []fileawardpoint.FileAwardPoint{{ent.FileAwardPoint{
	ID:                1,
	MerchantID:        1,
	DisplayName:       "import_file.xlsx",
	FileURL:           "https://a.com",
	ResultFileURL:     "https://a.com",
	Status:            0,
	StatsTotalRow:     100,
	StatsTotalSuccess: 98,
	CreatedAt:         time.Now(),
	UpdatedAt:         time.Now(),
	CreatedBy:         "quy.tm@teko.vn",
	UpdatedBy:         "quy.tm@teko.vn",
},
}, {ent.FileAwardPoint{
	ID:                2,
	MerchantID:        1,
	DisplayName:       "import_file_2.xlsx",
	FileURL:           "https://a.com",
	ResultFileURL:     "https://a.com",
	Status:            0,
	StatsTotalRow:     100,
	StatsTotalSuccess: 98,
	CreatedAt:         time.Now(),
	UpdatedAt:         time.Now(),
	CreatedBy:         "quy.tm@teko.vn",
	UpdatedBy:         "quy.tm@teko.vn",
},
},
}