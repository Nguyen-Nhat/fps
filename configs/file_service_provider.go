package config

// FileServiceConfig ...
type FileServiceConfig struct {
	Endpoint string           `mapstructure:"endpoint"`
	Paths    FileServicePaths `mapstructure:"paths"`
}

type FileServicePaths struct {
	UploadDoc string `mapstructure:"upload_doc"`
	Download  string `mapstructure:"download"`
	Delete    string `mapstructure:"delete"`
}

type FAltService struct {
	Endpoint  string `mapstructure:"endpoint"`
	MasterKey string `mapstructure:"master_key"`
	RestKey   string `mapstructure:"rest_key"`
	AppID     string `mapstructure:"app_id"`
	IsEnable  bool   `mapstructure:"is_enable"`
}
