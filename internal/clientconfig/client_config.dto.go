package clientconfig

type UIConfigImportHistoryTable struct {
	IsShowPreviewProcessFile bool   `json:"isShowPreviewProcessFile"`
	IsShowPreviewResultFile  bool   `json:"isShowPreviewResultFile"`
	IsShowDebug              bool   `json:"isShowDebug"`
	IsShowCreatedBy          bool   `json:"isShowCreatedBy"`
	IsShowReload             bool   `json:"isShowReload"`
	ColorScheme              string `json:"colorScheme"`
}

type UIConfigDTO struct {
	ImportHistoryTable UIConfigImportHistoryTable `json:"importHistoryTable"`
}

type GetClientConfigResDTO struct {
	ClientID              int32
	TenantID              string
	MaxFileSize           int32
	MerchantAttributeName string
	UsingMerchantAttrName bool
	InputFileTypes        []string
	ImportFileTemplateUrl string
	UIConfig              UIConfigDTO
}

func GetDefaultUiConfig() *UIConfigDTO {
	return &UIConfigDTO{
		ImportHistoryTable: UIConfigImportHistoryTable{
			IsShowPreviewProcessFile: true,
			IsShowPreviewResultFile:  true,
			IsShowDebug:              true,
			IsShowCreatedBy:          true,
			IsShowReload:             true,
		},
	}
}
