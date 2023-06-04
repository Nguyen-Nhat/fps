package config

// JobConfig ...
type JobConfig struct {
	FileProcessingConfig FileProcessingConfig `mapstructure:"file_processing"`
	FlattenConfig        SchedulerConfig      `mapstructure:"flatten"`
}

// FileProcessingConfig ...
type FileProcessingConfig struct {
	Schedule string `mapstructure:"schedule"`
}

// SchedulerConfig ...
type SchedulerConfig struct {
	Schedule string `mapstructure:"schedule"`
}
