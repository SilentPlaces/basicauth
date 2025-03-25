package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

// HashToSHA1 Generate SHA-1 from input string and return encoded string
func HashToSHA1(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
