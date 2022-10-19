package config

// JobConfig ...
type JobConfig struct {
	AwardPointJobConfig      AwardPointJobConfig      `mapstructure:"award_point"`
	UpdateStatusFAPJobConfig UpdateStatusFAPJobConfig `mapstructure:"update_status_fap"`
}

// AwardPointJobConfig ...
type AwardPointJobConfig struct {
	Schedule string `mapstructure:"schedule"`
}

// UpdateStatusFAPJobConfig ...
type UpdateStatusFAPJobConfig struct {
	Schedule        string `mapstructure:"schedule"`
	MaxCheckingTime int    `mapstructure:"max_checking_time"`
}
