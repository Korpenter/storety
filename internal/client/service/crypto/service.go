// Package service provides client services for interacting with the Storety server.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/samber/do"
	"io"
)

// Crypto is a service that provides encryption and decryption methods.
type Crypto struct {
	cfg *config.Config
}

// NewCrypto creates a new Crypto instance and returns a pointer to it.
// It takes a configuration object as a parameter.
func NewCrypto(i *do.Injector) *Crypto {
	cfg := do.MustInvoke[*config.Config](i)
	return &Crypto{
		cfg: cfg,
	}
}

// EncryptWithAES256 encrypts data with AES256.
// It takes a byte slice of data to be encrypted and returns the encrypted data or an error.
func (c *Crypto) EncryptWithAES256(data []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(c.cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	encBytes := aesgcm.Seal(nonce, nonce, data, nil)
	return encBytes, nil
}

// DecryptWithAES256 decrypts data with AES256.
// It takes a byte slice of encrypted data and returns the decrypted data or an error.
func (c *Crypto) DecryptWithAES256(data []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(c.cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	decBytes, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return decBytes, err
}
