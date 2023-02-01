package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ServerListen for specifying host & port
type ServerListen struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func (s ServerListen) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// ListenString for listen to 0.0.0.0
func (s ServerListen) ListenString() string {
	return fmt.Sprintf(":%d", s.Port)
}

// ServerConfig for configure HTTP & gRPC host & port
type ServerConfig struct {
	HTTP   ServerListen `mapstructure:"http"`
	GRPC   ServerListen `mapstructure:"grpc"`
	ApiKey string       `mapstructure:"api_key"`
}

// Config for app configuration
type Config struct {
	Server         ServerConfig   `mapstructure:"server"`
	Database       DatabaseConfig `mapstructure:"database"`
	Logger         LoggerConfig   `mapstructure:"logger"`
	JobConfig      JobConfig      `mapstructure:"job"`
	ProviderConfig ProviderConfig `mapstructure:"provider_config"`
}

type ProviderConfig struct {
	FileService FileServiceConfig `mapstructure:"file_service"`
}

const (
	EnvKeyRunProfile = "RUN_PROFILE"
)

const (
	ProfileTest = "TEST"
)

// Load config from config.yml
func Load(paths ...string) Config {
	vip := viper.New()

	vip.SetConfigName("config")
	vip.SetConfigType("yml")
	if len(paths) == 0 {
		vip.AddConfigPath(".") // ROOT
	} else {
		vip.AddConfigPath(paths[0])
	}

	vip.SetEnvPrefix("docker")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// workaround https://github.com/spf13/viper/issues/188#issuecomment-399518663
	// to allow read from environment variables when Unmarshal
	for _, key := range vip.AllKeys() {
		val := vip.Get(key)
		vip.Set(key, val)
	}

	fmt.Println("===== Config file used:", vip.ConfigFileUsed())

	cfg := Config{}
	err = vip.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	checkRunProfile(&cfg)

	return cfg
}

func checkRunProfile(cfg *Config) {
	runProfile := os.Getenv(EnvKeyRunProfile)
	fmt.Println("\n===== Running in Profile =", runProfile)
	if strings.ToUpper(runProfile) == ProfileTest {
		testingDBName := cfg.Database.MySQL.DBName + "_test" // add suffix `_test`
		cfg.Database.MySQL.DBName = testingDBName
		fmt.Println("===== ---> When running in TEST profile, DB is switched to", testingDBName)
	}
	fmt.Println() // only new line
}
