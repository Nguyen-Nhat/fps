package config

import (
	"fmt"
)

type DatabaseConfig struct {
	MySQL MySqlDBConfig `mapstructure:"mysql"`
	Debug DebugDBConfig `mapstructure:"debug"`
}

type MySqlDBConfig struct {
	DBName   string `mapstructure:"db_name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Options  string `mapstructure:"options"`
}

type DebugDBConfig struct {
	Enable   bool   `mapstructure:"enable"`
	PingCron string `mapstructure:"ping_cron"`
}

func (c *MySqlDBConfig) DatabaseURI() string {
	uri := fmt.Sprintf("mysql://%s:%s@%s:%d/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Options)
	return uri
}

func (c *MySqlDBConfig) DatabaseTcpURI() string {
	uri := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Options)
	return uri
}
