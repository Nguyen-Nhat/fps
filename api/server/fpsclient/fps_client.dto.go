package fpsclient

import (
	"fmt"
	"net/http"

	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// GetListClientData ...
type GetListClientData struct {
	request.PageRequest
	Name string `json:"name"`
}

// GetListClientResponse ...
type GetListClientResponse struct {
	Clients    []ClientDTO            `json:"clients"`
	Pagination response.PaginationNew `json:"pagination"`
}

// ClientDTO ...
type ClientDTO struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SampleFileURL string `json:"sampleFileURL"`
	CreatedAt     int64  `json:"createdAt"`
	CreatedBy     string `json:"createdBy"`
}

// bindAndValidateRequestParams ...
func bindAndValidateRequestParams(r *http.Request) (*GetListClientData, error) {
	// 1. Bind request params to struct
	data := &GetListClientData{}
	params := r.URL.Query()
	if err := utils.BindRequestParamsToStruct(data, params, "json"); err != nil {
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return nil, error2.ErrInvalidRequestWithError(err)
	}
	data.InitDefaultValue()

	// 2. Validate request params
	if len(data.Name) > 255 {
		err := fmt.Errorf("name out of range")
		logger.Errorf("GetListClientData is invalid ... error %+v", err)
		return nil, error2.ErrInvalidRequestWithError(err)
	}

	// 3. Validate pagination
	if err := data.ValidatePagination(); err != nil {
		logger.Errorf("GetListClientData is invalid pagination ... error %+v", err)
		return nil, error2.ErrInvalidRequestWithError(err)
	}

	return data, nil
}
