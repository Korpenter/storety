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
		userId    string
		hashedKey []byte
		salt      []byte
	}{
		{
			userId:    "test_user1",
			hashedKey: make([]byte, 32),
			salt:      make([]byte, 16),
		},
		{
			userId:    "test_user2",
			hashedKey: make([]byte, 32),
			salt:      make([]byte, 16),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.userId, func(t *testing.T) {
			rand.Read(tt.hashedKey)
			rand.Read(tt.salt)
			err := SaveHashedKeyAndSalt(testFile, tt.userId, tt.hashedKey, tt.salt)
			assert.NoError(t, err, "Error saving hashed key and salt")

			retrievedHashedKey, retrievedSalt, err := GetHashedKeyAndSalt(testFile, tt.userId)
			assert.NoError(t, err, "Error retrieving hashed key and salt")

			assert.Equal(t, tt.hashedKey, retrievedHashedKey, "Retrieved hashed key does not match the original")
			assert.Equal(t, tt.salt, retrievedSalt, "Retrieved salt does not match the original")
		})
	}

	os.Remove(testFile)
}
