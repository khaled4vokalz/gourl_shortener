package cache

import "github.com/khaled4vokalz/gourl_shortener/internal/config"

func GetCache(conf config.CacheConfig) Cache {
	return NewRedisCache(conf)
}
