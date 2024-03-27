package database

import (
	"fmt"
	"net/url"
)

// DBConfig used to set config for database.
type DBConfig interface {
	String() string
	DSN() string
}

// Config used to set base config for all database types.
type Config struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
}

// DSN returns the Domain Source Name.
func (c Config) DSN() string {
	options := c.Options
	if options != "" {
		if options[0] != '?' {
			options = "?" + options
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		options)
}

// MySQLConfig used to set config for MySQL.
type MySQLConfig struct {
	Config `mapstructure:",squash"`
}

// String returns MySQL connection URI.
func (c MySQLConfig) String() string {
	return fmt.Sprintf("mysql://%s", c.DSN())
}

// MySQLDefaultConfig returns default config for mysql, usually use on development.
func MySQLDefaultConfig() MySQLConfig {
	return MySQLConfig{Config{
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "sample",
		Username: "default",
		Password: "secret",
		Options:  "?parseTime=true",
	}}
}
