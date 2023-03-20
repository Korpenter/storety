package tls

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCert(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cert_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	c := &config.Config{
		CertFile: filepath.Join(tmpDir, "cert.pem"),
		KeyFile:  filepath.Join(tmpDir, "key.pem"),
	}

	err = GenerateCert(c)
	assert.NoError(t, err)

	certData, err := os.ReadFile(c.CertFile)
	assert.NoError(t, err)

	block, _ := pem.Decode(certData)
	assert.NotNil(t, block)
	assert.Equal(t, block.Type, "CERTIFICATE")

	_, err = x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)

	keyData, err := os.ReadFile(c.KeyFile)
	assert.NoError(t, err)

	block, _ = pem.Decode(keyData)
	assert.NotNil(t, block)
	assert.Equal(t, block.Type, "PRIVATE KEY")

	_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	assert.NoError(t, err)
}
