package config

// JobConfig ...
type JobConfig struct {
	FileProcessingConfig FileProcessingConfig `mapstructure:"file_processing"`
	FlattenConfig        SchedulerConfig      `mapstructure:"flatten"`
	ExecuteTaskConfig    SchedulerConfig      `mapstructure:"execute_task"`
	UpdateStatusConfig   SchedulerConfig      `mapstructure:"update_status"`
}

// FileProcessingConfig ...
type FileProcessingConfig struct {
	Schedule string `mapstructure:"schedule"`
}

// SchedulerConfig ...
type SchedulerConfig struct {
	Schedule string `mapstructure:"schedule"`
}
