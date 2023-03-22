package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/constants"
	"os"
)

// SaveAuthData saves the hashed key and salt for the given user ID
func SaveAuthData(filename, userId string, hashedKey, salt []byte, authToken, refreshToken string) error {
	data, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	keysAndSalts := make(map[string]models.AuthData)
	if len(data) > 0 {
		err = json.Unmarshal(data, &keysAndSalts)
		if err != nil {
			return err
		}
	}

	keysAndSalts[userId] = models.AuthData{
		HashedKey:    base64.StdEncoding.EncodeToString(hashedKey),
		Salt:         base64.StdEncoding.EncodeToString(salt),
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}

	newData, err := json.MarshalIndent(keysAndSalts, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, newData, 0600)
	if err != nil {
		return err
	}
	return nil
}

// GetAuthData returns the hashed key and salt for the given user ID
func GetAuthData(filename, userId string) ([]byte, []byte, string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, "", "", err
	}
	keysAndSalts := make(map[string]models.AuthData)
	err = json.Unmarshal(data, &keysAndSalts)
	if err != nil {
		return nil, nil, "", "", err
	}
	userData, ok := keysAndSalts[userId]
	if !ok {
		return nil, nil, "", "", constants.ErrUserNotFound
	}
	hashedKey, err := base64.StdEncoding.DecodeString(userData.HashedKey)
	if err != nil {
		return nil, nil, "", "", err
	}
	salt, err := base64.StdEncoding.DecodeString(userData.Salt)
	if err != nil {
		return nil, nil, "", "", err
	}

	return hashedKey, salt, userData.AuthToken, userData.RefreshToken, nil
}
