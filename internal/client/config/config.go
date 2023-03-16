// Package config provides configuration management for the Storety client.
package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config is the configuration for the Storety client.
type Config struct {
	ServiceAddress  string `mapstructure:"service_address"`
	JWTAuthToken    string `mapstructure:"jwt_auth_token"`
	JWTRefreshToken string `mapstructure:"jwt_refresh_token"`
	EncryptionKey   []byte
}

// NewConfig creates a new Config instance and returns a pointer to it.
// It reads the configuration from the "demo.yaml" file and sets default values if necessary.
func NewConfig() *Config {
	viper.SetConfigFile("cfg.yaml")
	viper.SetDefault("service_address", ":8081")
	viper.SetDefault("jwt_auth_token", nil)
	viper.SetDefault("jwt_refresh_token", nil)
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
	viper.Set("jwt_auth_token", auth)
	viper.Set("jwt_refresh_token", refresh)
	viper.WriteConfig()
	return nil
}

// UpdateKey updates the encryption key in the config.
// It takes a password string as a parameter and updates the configuration accordingly.
func (c *Config) UpdateKey(password string) error {
	key := []byte(password)
	if len(key) < 32 {
		for {
			key = append(key, key[0])
			if len(key) == 32 {
				break
			}
		}
	} else if len(key) > 32 {
		key = key[:32]
	}
	c.EncryptionKey = key
	return nil
}
