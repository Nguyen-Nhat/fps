package config

// LoyaltyConfig ...
type LoyaltyConfig struct {
	Endpoint string       `mapstructure:"endpoint"`
	APIKey   string       `mapstructure:"x_api_key"`
	Paths    LoyaltyPaths `mapstructure:"paths"`
}

type LoyaltyPaths struct {
	TxnGetList string `mapstructure:"txn_get_list"`
	TxnGrant   string `mapstructure:"txn_grant"`
}
