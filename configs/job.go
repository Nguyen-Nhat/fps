package config

// JobConfig ...
type JobConfig struct {
	AwardPointJobConfig AwardPointJobConfig `mapstructure:"award_point"`
}

// AwardPointJobConfig ...
type AwardPointJobConfig struct {
	Schedule string `mapstructure:"schedule"`
}
