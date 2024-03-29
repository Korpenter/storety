// Package config provides the configuration for the Storety server.
package config

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/shopspring/decimal"
)

// Config is the configuration for the Storety server.
type Config struct {
	ServiceAddress          string `envconfig:"RUN_ADDRESS" default:":8081"`
	PostgresURI             string `envconfig:"DATABASE_URI" default:""`
	JWTAuthKey              string `envconfig:"JWT_AUTH_KEY" default:"defaultAuthKey"`
	JWTAuthLifeTimeHours    int    `envconfig:"JWT_LIFETIME_HOURS" default:"24"`
	JWTRefreshLifeTimeHours int    `envconfig:"JWT_REFRESH_LIFETIME_HOURS" default:"48"`
	CertFile                string `envconfig:"TLS_CERT_FILE" default:"cert.pem" json:"cert_file"`
	KeyFile                 string `envconfig:"TLS_KEY_FILE" default:"key.pem" json:"key_file"`
}

// NewConfig creates a new Config instance and returns a pointer to it.
// It reads configuration values from environment variables and command-line flags.
func NewConfig() *Config {
	var cfg Config
	decimal.MarshalJSONWithoutQuotes = true
	envconfig.MustProcess("", &cfg)
	flag.StringVar(&cfg.ServiceAddress, "a", cfg.ServiceAddress, "grpcServer address")
	flag.StringVar(&cfg.PostgresURI, "d", cfg.PostgresURI, "db address")
	flag.StringVar(&cfg.JWTAuthKey, "j", cfg.JWTAuthKey, "token token key")
	flag.IntVar(&cfg.JWTAuthLifeTimeHours, "l", cfg.JWTAuthLifeTimeHours, "token token token lifetime in hours")
	flag.IntVar(&cfg.JWTRefreshLifeTimeHours, "r", cfg.JWTRefreshLifeTimeHours, "token refresh token lifetime in hours")
	flag.StringVar(&cfg.CertFile, "c", cfg.CertFile, "tls cert file path")
	flag.StringVar(&cfg.KeyFile, "k", cfg.KeyFile, "tls key file path")
	flag.Parse()
	return &cfg
}
