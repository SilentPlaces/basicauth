package strings

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generate random base64 string with lengthInBytes size
func GenerateRandomString(lengthInBytes int) (string, error) {
	tokenBytes := make([]byte, lengthInBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}
