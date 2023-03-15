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
	JWTRefreshLifeTimeHours int    `envconfig:"JWT_LIFETIME_HOURS" default:"48"`
}

// NewConfig creates a new Config.
func NewConfig() *Config {
	var cfg Config
	decimal.MarshalJSONWithoutQuotes = true
	envconfig.MustProcess("", &cfg)
	flag.StringVar(&cfg.ServiceAddress, "a", cfg.ServiceAddress, "grpcServer address")
	flag.StringVar(&cfg.PostgresURI, "d", cfg.PostgresURI, "db address")
	flag.StringVar(&cfg.JWTAuthKey, "j", cfg.JWTAuthKey, "token token key")
	flag.IntVar(&cfg.JWTAuthLifeTimeHours, "l", cfg.JWTAuthLifeTimeHours, "token token token lifetime in hours")
	flag.IntVar(&cfg.JWTRefreshLifeTimeHours, "r", cfg.JWTRefreshLifeTimeHours, "token refresh token lifetime in hours")
	flag.Parse()
	return &cfg
}
