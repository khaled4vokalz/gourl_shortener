package cache

import (
	"time"

	errors "github.com/khaled4vokalz/gourl_shortener/internal/common"
)

func NewNoOpCache() Cache {
	return &NoOpCache{}
}

func (r *NoOpCache) Set(key string, value string, expiration time.Duration) error {
	return nil
}

func (r *NoOpCache) Get(key string) (string, error) {
	return "", errors.NotFound
}
