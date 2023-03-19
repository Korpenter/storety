// Package config provides configuration management for the Storety client.
package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config is the configuration for the Storety client.
type Config struct {
	ServiceAddress  string `mapstructure:"service_address"`
	JWTAuthToken    string
	JWTRefreshToken string
	CertFile        string `mapstructure:"cert_file"`
	KeyFile         string `mapstructure:"key_file"`
	SaltsFile       string `mapstructure:"salts_file"`
	DBFilePrefix    string `mapstructure:"db_path"`
	EncryptionKey   []byte
}

// NewConfig creates a new Config instance and returns a pointer to it.
// It reads the configuration from the "demo.yaml" file and sets default values if necessary.
func NewConfig() *Config {
	viper.SetConfigFile("cfg.yaml")
	viper.SetDefault("service_address", ":8081")
	viper.SetDefault("jwt_auth_token", nil)
	viper.SetDefault("jwt_refresh_token", nil)
	viper.SetDefault("cert_file", "cert.pem")
	viper.SetDefault("key_file", "key.pem")
	viper.SetDefault("salts_file", "salts.json")
	viper.SetDefault("db_path", "")
	c := &Config{}
	viper.ReadInConfig()
	if err := viper.Unmarshal(c); err != nil {
		log.Fatal(err)
	}
	viper.WriteConfig()
	return c
}

// UpdateTokens updates the tokens in the config.
// It takes the new auth and refresh tokens as parameters and updates the configuration accordingly.
func (c *Config) UpdateTokens(auth, refresh string) error {
	c.JWTAuthToken = auth
	c.JWTRefreshToken = refresh
	return nil
}

// UpdateKey updates the encryption key in the config.
// It takes a password string as a parameter and updates the configuration accordingly.
func (c *Config) UpdateKey(key []byte) {
	c.EncryptionKey = key
}
