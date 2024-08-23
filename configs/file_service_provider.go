package config

// FileServiceConfig ...
type FileServiceConfig struct {
	Endpoint              string           `mapstructure:"endpoint"`
	ExternalEndpointRegex string           `mapstructure:"internal_endpoint_regex"`
	Paths                 FileServicePaths `mapstructure:"paths"`
}

type FileServicePaths struct {
	UploadDoc string `mapstructure:"upload_doc"`
	Download  string `mapstructure:"download"`
	Delete    string `mapstructure:"delete"`
}

type FAltService struct {
	Endpoint         string `mapstructure:"endpoint"`
	MasterKey        string `mapstructure:"master_key"`
	RestKey          string `mapstructure:"rest_key"`
	SessionExpiredIn int    `mapstructure:"session_expired_in"`
	AppID            string `mapstructure:"app_id"`
	IsEnable         bool   `mapstructure:"is_enable"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
}
