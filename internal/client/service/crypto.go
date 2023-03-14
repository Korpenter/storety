package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/Mldlr/storety/internal/client/config"
	"io"
)

type Crypto struct {
	cfg *config.Config
}

func NewCrypto(cfg *config.Config) *Crypto {
	return &Crypto{
		cfg: cfg,
	}
}

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
