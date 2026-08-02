package helper

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken membuat token acak hex sepanjang n byte.
func GenerateRandomToken(n int) (string, error) {
	if n <= 0 {
		n = 32
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
