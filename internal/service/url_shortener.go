package service

import (
	"crypto/sha256"

	"github.com/jxskiss/base62"
)

func GenerateShortenedURL(original string, bytes_to_take int8) string {
	hash := sha256.Sum256([]byte(original))
	encoded := base62.EncodeToString(hash[:bytes_to_take])
	return encoded
}
