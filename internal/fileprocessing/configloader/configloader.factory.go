package configloader

import (
	"errors"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

// CfgLoaderFactory ....................................................................................................
type CfgLoaderFactory struct {
	// configCache ... map[clientID]ConfigMappingMD -> simulate an In-Memory Cache
	configCache map[int32]ConfigMappingMD

	// services
	cfgMappingService configmapping.Service
	cfgTaskService    configtask.Service
}

func NewConfigLoaderFactory(
	cfgMappingService configmapping.Service,
	cfgTaskService configtask.Service,
) *CfgLoaderFactory {
	// configCache ...
	return &CfgLoaderFactory{
		configCache: make(map[int32]ConfigMappingMD),
		// service
		cfgMappingService: cfgMappingService,
		cfgTaskService:    cfgTaskService,
	}
}

func (factory *CfgLoaderFactory) GetConfigLoader(file fileprocessing.ProcessingFile) (ConfigMappingMD, error) {
	// 1. Check config in Cache
	configInCache, isOk := factory.configCache[file.ClientID]
	if isOk {
		return configInCache, nil
	}

	// 2. If not existed, load config then put to cache
	// 2.1. Get Config Loader
	cfgLoader, isOk := factory.factory(file)
	if !isOk {
		logger.Error("failed to get config loader")
		return ConfigMappingMD{}, errors.New("failed to get config loader")
	}
	// 2.2. Load config, then put to cache
	config, err := cfgLoader.Load(file)
	if err != nil {
		return ConfigMappingMD{}, err
	} else {
		factory.configCache[file.ClientID] = configInCache
		return config, nil
	}
}

// Private method ------------------------------------------------------------------------------------------------------

// factory ... return (configLoader, isOk)
func (factory *CfgLoaderFactory) factory(file fileprocessing.ProcessingFile) (ConfigLoader, bool) {
	// Currently, we only support 1 config loader
	// We will support more Implementation of ConfigLoader
	return &databaseConfigLoaderV1{
		factory.cfgMappingService, factory.cfgTaskService,
	}, true
}
