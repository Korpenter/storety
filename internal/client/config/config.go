package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	ServiceAddress  string `envconfig:"STORETY_ADDRESS" default:":8081"`
	JWTAuthToken    string `envconfig:"JWT_AUTH_TOKEN" default:"jwt_auth_token"`
	JWTRefreshToken string `envconfig:"JWT_REFRESH_TOKEN" default:"jwt_refresh_token"`
	EncryptionKey   []byte
}

func NewConfig() *Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	fmt.Println(cfg)
	return &cfg
}

func (c *Config) UpdateConfig(auth, refresh, password string) error {
	err := os.Setenv("JWT_AUTH_TOKEN", auth)
	if err != nil {
		return err
	}
	err = os.Setenv("JWT_REFRESH_TOKEN", refresh)
	if err != nil {
		return err
	}

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

	c.JWTAuthToken = auth
	c.JWTRefreshToken = refresh
	c.EncryptionKey = key
	return nil
}
