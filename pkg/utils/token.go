package utils

import (
	"crypto/rand"

	"github.com/jxskiss/base62"
)

func GenerateSecureToken(length int) (string, error) {
	var b = make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base62.EncodeToString(b), nil
}
