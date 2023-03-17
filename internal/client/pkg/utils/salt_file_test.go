package utils

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testFile = "test_hashed_keys_and_salts.json"

func TestSaveAndGetHashedKeyAndSalt(t *testing.T) {
	testCases := []struct {
		userId       string
		hashedKey    []byte
		salt         []byte
		authToken    string
		refreshToken string
	}{
		{
			userId:       "test_user1",
			hashedKey:    make([]byte, 32),
			salt:         make([]byte, 16),
			authToken:    "authToken",
			refreshToken: "refreshToken",
		},
		{
			userId:       "test_user2",
			hashedKey:    make([]byte, 32),
			salt:         make([]byte, 16),
			authToken:    "authToken2",
			refreshToken: "refreshToken2",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.userId, func(t *testing.T) {
			rand.Read(tt.hashedKey)
			rand.Read(tt.salt)
			err := SaveAuthData(testFile, tt.userId, tt.hashedKey, tt.salt, tt.authToken, tt.refreshToken)
			assert.NoError(t, err, "Error saving hashed key and salt")

			retrievedHashedKey, retrievedSalt, authToken, refreshToken, err := GetAuthData(testFile, tt.userId)
			assert.NoError(t, err, "Error retrieving hashed key and salt")

			assert.Equal(t, tt.hashedKey, retrievedHashedKey, "Retrieved hashed key does not match the original")
			assert.Equal(t, tt.salt, retrievedSalt, "Retrieved salt does not match the original")
			assert.Equal(t, tt.authToken, authToken, "Retrieved authToken does not match the original")
			assert.Equal(t, tt.refreshToken, refreshToken, "Retrieved refreshToken does not match the original")
		})
	}

	os.Remove(testFile)
}
