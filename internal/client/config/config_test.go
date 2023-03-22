package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.NotNil(t, cfg)
	defer os.Remove("cfg.yaml")
	expectedCfg := &Config{
		ServiceAddress: ":8081",
		CertFile:       "cert.pem",
		KeyFile:        "key.pem",
		SaltsFile:      "salts.json",
		DBFilePrefix:   "",
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestUpdateTokens(t *testing.T) {
	cfg := &Config{}
	cfg.UpdateTokens("testAuthToken", "testRefreshToken")
	assert.Equal(t, "testAuthToken", cfg.JWTAuthToken)
	assert.Equal(t, "testRefreshToken", cfg.JWTRefreshToken)
}

func TestUpdateKey(t *testing.T) {
	cfg := &Config{}
	key := []byte("testEncryptionKey")
	cfg.UpdateKey(key)
	assert.Equal(t, key, cfg.EncryptionKey)
}
