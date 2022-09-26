// YOU CAN EDIT YOUR CUSTOM CONFIG HERE

package config

import (
	"go.tekoapis.com/kitchen/database"
)

type Config struct {
	Base `mapstructure:",squash"`
	// Custom here
	MySQL           database.MySQLConfig `json:"mysql" mapstructure:"mysql"`
	MigrationFolder string               `json:"migration_folder" mapstructure:"migration_folder"`
}

func loadDefaultConfig() *Config {
	return &Config{
		Base:            *defaultBaseConfig(),
		MySQL:           database.MySQLDefaultConfig(),
		MigrationFolder: "file://migrations",
	}
}
