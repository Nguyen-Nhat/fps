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
