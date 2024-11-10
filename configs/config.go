package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
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
	Server          ServerConfig   `mapstructure:"server"`
	Database        DatabaseConfig `mapstructure:"database"`
	Logger          LoggerConfig   `mapstructure:"logger"`
	JobConfig       JobConfig      `mapstructure:"job"`
	ProviderConfig  ProviderConfig `mapstructure:"provider_config"`
	FlagSupHost     string         `mapstructure:"flag_sup_host"`
	MigrationFolder string         `mapstructure:"migration_folder"`
	Kafka           KafkaConfig    `mapstructure:"kafka"`
	ExtraConfig     ExtraConfig    `mapstructure:"extra_config"`
	MessageFolder   string         `mapstructure:"message_folder"`
}

type ProviderConfig struct {
	FileService FileServiceConfig `mapstructure:"file_service"`
	FAltService FAltService       `mapstructure:"f_alt_service"`
}

type KafkaConfig struct {
	ConnectionHost     string `mapstructure:"connection_host"`
	ConsumerPrefixName string `mapstructure:"consumer_prefix_name"`

	// Topic
	UpdateResultAsyncTopic string `mapstructure:"update_result_async_topic"`
}

type ExtraConfig struct {
	Epic1139EnableSellers    string  `mapstructure:"epic1139_enable_sellers"` // list of seller id, split by comma
	Epic1139EnableSellersObj []int32 `mapstructure:"-"`
}

const (
	EnvKeyRunProfile = "RUN_PROFILE"
)

const (
	ProfileTest = "TEST"
)

var Cfg = Config{}

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

	cfg := loadDefaultConfig()
	err = vip.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if err := cfg.parse(); err != nil {
		panic(err)
	}

	checkRunProfile(&cfg)

	return cfg
}

func loadDefaultConfig() Config {
	return Config{
		ProviderConfig: ProviderConfig{
			FileService: FileServiceConfig{
				ExternalEndpointRegex: `https:\/\/files(\.dev|\.stag|\.prod)?\.tekoapis\.(net|com)`,
			},
		},
	}
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

func (c *Config) parse() error {
	// parse extra config
	if c.ExtraConfig.Epic1139EnableSellers != constant.EmptyString {
		sellers := strings.Split(c.ExtraConfig.Epic1139EnableSellers, constant.SplitByComma)
		c.ExtraConfig.Epic1139EnableSellersObj = make([]int32, len(sellers))
		for idx, sellerStr := range sellers {
			sellerId, err := strconv.Atoi(sellerStr)
			if err != nil {
				return err
			}
			c.ExtraConfig.Epic1139EnableSellersObj[idx] = int32(sellerId)
		}
	}

	// parse another here

	return nil
}
