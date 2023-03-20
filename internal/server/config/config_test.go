package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	oldServiceAddress := os.Getenv("RUN_ADDRESS")
	oldPostgresURI := os.Getenv("DATABASE_URI")
	oldJWTAuthKey := os.Getenv("JWT_AUTH_KEY")
	oldJWTAuthLifeTimeHours := os.Getenv("JWT_LIFETIME_HOURS")
	oldJWTRefreshLifeTimeHours := os.Getenv("JWT_REFRESH_LIFETIME_HOURS")
	oldCertFile := os.Getenv("TLS_CERT_FILE")
	oldKeyFile := os.Getenv("TLS_KEY_FILE")

	os.Setenv("RUN_ADDRESS", ":1234")
	os.Setenv("DATABASE_URI", "test_db")
	os.Setenv("JWT_AUTH_KEY", "testKey")
	os.Setenv("JWT_LIFETIME_HOURS", "5")
	os.Setenv("JWT_REFRESH_LIFETIME_HOURS", "10")
	os.Setenv("TLS_CERT_FILE", "test_cert.pem")
	os.Setenv("TLS_KEY_FILE", "test_key.pem")

	cfg := NewConfig()

	assert.Equal(t, ":1234", cfg.ServiceAddress)
	assert.Equal(t, "test_db", cfg.PostgresURI)
	assert.Equal(t, "testKey", cfg.JWTAuthKey)
	assert.Equal(t, 5, cfg.JWTAuthLifeTimeHours)
	assert.Equal(t, 10, cfg.JWTRefreshLifeTimeHours)
	assert.Equal(t, "test_cert.pem", cfg.CertFile)
	assert.Equal(t, "test_key.pem", cfg.KeyFile)

	os.Setenv("RUN_ADDRESS", oldServiceAddress)
	os.Setenv("DATABASE_URI", oldPostgresURI)
	os.Setenv("JWT_AUTH_KEY", oldJWTAuthKey)
	os.Setenv("JWT_LIFETIME_HOURS", oldJWTAuthLifeTimeHours)
	os.Setenv("JWT_REFRESH_LIFETIME_HOURS", oldJWTRefreshLifeTimeHours)
	os.Setenv("TLS_CERT_FILE", oldCertFile)
	os.Setenv("TLS_KEY_FILE", oldKeyFile)
}
