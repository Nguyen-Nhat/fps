package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MySQLHost     string            `envconfig:"mysql_host" required:"true"`
	MySQLPort     string            `envconfig:"mysql_port" default:"3306"`
	MySQLUsername string            `envconfig:"mysql_username" required:"true"`
	MySQLPassword string            `envconfig:"mysql_password" required:"true"`
	MySQLDatabase string            `envconfig:"mysql_database" required:"true"`
	MySQLOptions  map[string]string `envconfig:"mysql_options" default:"parseTime:true"`
}

func MustParseConfig(prefix string) *Config {
	var cfg Config
	envconfig.MustProcess("idp", &cfg)
	return &cfg
}

func ParseConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process("idp", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) DatabaseURI() string {
	var optParams []string
	for opt, val := range c.MySQLOptions {
		optParams = append(optParams, opt+"="+val)
	}
	optStr := strings.Join(optParams, "&")
	uri := fmt.Sprintf("mysql://%s:%s@%s:%s/%s?%s",
		c.MySQLUsername,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
		optStr)
	return uri
}
