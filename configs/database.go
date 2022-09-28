package config

import (
	"fmt"
)

type DatabaseConfig struct {
	MySQL MySqlDBConfig `mapstructure:"mysql"`
}

type MySqlDBConfig struct {
	DBName   string `mapstructure:"db_name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Options  string `mapstructure:"pptions"`
}

func (c *MySqlDBConfig) DatabaseURI() string {
	uri := fmt.Sprintf("mysql://%s:%s@%s:%s/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Options)
	return uri
}
