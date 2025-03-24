package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(lengthInBytes int) (string, error) {
	tokenBytes := make([]byte, lengthInBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}
