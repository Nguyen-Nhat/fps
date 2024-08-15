package clientconfig

import (
	"fmt"
	"net/http"

	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// Request DTO =========================================================================================================

type GetClientConfigRequest struct {
	ClientId int32 `json:"clientId"`
}

func bindAndValidateRequestParams(r *http.Request, data *GetClientConfigRequest) error {
	// 1.Bind request params to struct
	params := r.URL.Query()
	if err := utils.BindRequestParamsToStruct(data, params, "json"); err != nil {
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}

	// 2.Validate request params
	if data.ClientId == 0 {
		err := fmt.Errorf("missing required param: clientId")
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}

	return nil
}

// Response DTO ========================================================================================================

type UIConfigImportHistoryTable struct {
	IsShowPreviewProcessFile bool   `json:"isShowPreviewProcessFile"`
	IsShowPreviewResultFile  bool   `json:"isShowPreviewResultFile"`
	IsShowDebug              bool   `json:"isShowDebug"`
	IsShowCreatedBy          bool   `json:"isShowCreatedBy"`
	IsShowReload             bool   `json:"isShowReload"`
	ColorScheme              string `json:"colorScheme"`
}

type UIConfig struct {
	ImportHistoryTable UIConfigImportHistoryTable `json:"importHistoryTable"`
}

type GetClientConfigResponse struct {
	ClientID              int32    `json:"clientId"`
	TenantID              string   `json:"tenantId"`
	MaxFileSize           int32    `json:"maxFileSize"`
	MerchantAttributeName string   `json:"merchantAttributeName"`
	UsingMerchantAttrName bool     `json:"usingMerchantAttrName"`
	InputFileTypes        []string `json:"inputFileTypes"`
	ImportFileTemplateUrl string   `json:"importFileTemplateUrl"`
	UIConfig              UIConfig `json:"uiConfig"`
}
