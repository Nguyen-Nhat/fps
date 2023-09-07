package filerow

import (
	"net/http"

	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// GetListFileRowData ...
type GetListFileRowData struct {
	request.PageRequest
}

// GetListFileRowResponse ...
type GetListFileRowResponse struct {
	Rows       []fileprocessingrow.GetListFileRowsItem `json:"rows"`
	Pagination response.PaginationNew                  `json:"pagination"`
}

// bindAndValidateRequestParams ...
func bindAndValidateRequestParams(r *http.Request) (*GetListFileRowData, error) {
	// 1. Bind request params to struct
	data := &GetListFileRowData{}
	params := r.URL.Query()
	if err := utils.BindRequestParamsToStruct(data, params, "json"); err != nil {
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return nil, error2.ErrInvalidRequestWithError(err)
	}
	data.InitDefaultValue()

	// 2. Validate request other data
	// ...

	// 3. Validate pagination
	if err := data.ValidatePagination(); err != nil {
		logger.Errorf("GetListClientData is invalid pagination ... error %+v", err)
		return nil, error2.ErrInvalidRequestWithError(err)
	}

	return data, nil
}
