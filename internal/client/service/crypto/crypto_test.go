package crypto

import (
	"bytes"
	"crypto/rand"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecryptWithAES256(t *testing.T) {
	var err error
	cfg := &config.Config{
		EncryptionKey: make([]byte, 32),
	}
	_, err = rand.Read(cfg.EncryptionKey)
	assert.NoError(t, err, "Error generating encryption key")
	cryptoSvc := NewCrypto(cfg)

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte(""),
		},
		{
			name: "small data",
			data: []byte("hello world"),
		},
		{
			name: "large data",
			data: make([]byte, 1000),
		},
	}
	var encryptedData, decryptedData []byte

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedData, err = cryptoSvc.EncryptWithAES256(tt.data)
			assert.NoError(t, err, "Error encrypting data")

			decryptedData, err = cryptoSvc.DecryptWithAES256(encryptedData)
			assert.NoError(t, err, "Error decrypting data")

			assert.True(t, bytes.Equal(tt.data, decryptedData), "Original data and decrypted data do not match")
		})
	}
}
