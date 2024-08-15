package clientconfig

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fpsclient"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		GetClientConfigById(context.Context, int32) (*GetClientConfigResDTO, error)
	}

	ServiceImpl struct {
		fpsClientRepo     fpsclient.Repo
		configMappingRepo configmapping.Repo
	}
)

var _ Service = &ServiceImpl{}

func NewService(db *sql.DB) *ServiceImpl {
	return &ServiceImpl{
		fpsClientRepo:     fpsclient.NewRepo(db),
		configMappingRepo: configmapping.NewRepo(db),
	}
}

func (s *ServiceImpl) GetClientConfigById(ctx context.Context, id int32) (*GetClientConfigResDTO, error) {
	fpsClient, err := s.fpsClientRepo.FindById(ctx, id)
	if err != nil {
		logger.Infof("Error in GetClientConfigById fpsClientRepo.FindById err %+v", err)
		return nil, err
	}

	configMapping, err := s.configMappingRepo.FindByClientID(ctx, id)
	if err != nil {
		logger.Infof("Error in GetClientConfigById configMappingRepo.FindByClientID err %+v", err)
		return nil, err
	}

	uiConfig := &UIConfigDTO{}

	// if not config, use default config, otherwise parse config from db
	if configMapping.UIConfig == constant.EmptyString {
		uiConfig = GetDefaultUiConfig()
	} else if err = json.Unmarshal([]byte(configMapping.UIConfig), uiConfig); err != nil {
		logger.Infof("Error in GetClientConfigById, unmashal ui config err %+v", err)
		return nil, err
	}

	return &GetClientConfigResDTO{
		ClientID:              fpsClient.ClientID,
		TenantID:              configMapping.TenantID,
		MaxFileSize:           configMapping.MaxFileSize,
		MerchantAttributeName: configMapping.MerchantAttributeName,
		UsingMerchantAttrName: configMapping.UsingMerchantAttrName,
		InputFileTypes:        strings.Split(configMapping.InputFileType, constant.SplitByComma),
		ImportFileTemplateUrl: fpsClient.ImportFileTemplateURL,
		UIConfig:              *uiConfig,
	}, nil
}
