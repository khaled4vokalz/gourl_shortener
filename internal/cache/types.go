package cache

import "time"

type NoOpCache struct{}

type Cache interface {
	Set(key string, value string, ttlSeconds time.Duration) error
	Get(key string) (string, error)
}
