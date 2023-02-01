package config

// JobConfig ...
type JobConfig struct {
	FileProcessingConfig FileProcessingConfig `mapstructure:"file_processing"`
}

// FileProcessingConfig ...
type FileProcessingConfig struct {
	Schedule string `mapstructure:"schedule"`
}
