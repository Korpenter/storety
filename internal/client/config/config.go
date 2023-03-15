package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServiceAddress  string `mapstructure:"service_address"`
	JWTAuthToken    string `mapstructure:"jwt_auth_token"`
	JWTRefreshToken string `mapstructure:"jwt_refresh_token"`
	EncryptionKey   []byte
}

func NewConfig() *Config {
	viper.SetConfigFile("demo.yaml")
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

func (c *Config) UpdateTokens(auth, refresh string) error {
	c.JWTAuthToken = auth
	c.JWTRefreshToken = refresh
	viper.Set("jwt_auth_token", auth)
	viper.Set("jwt_refresh_token", refresh)
	viper.WriteConfig()
	return nil
}

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
