package cache

import "github.com/khaled4vokalz/gourl_shortener/internal/config"

func GetCache(conf config.CacheConfig) (Cache, error) {
	if conf.Enabled == true {
		return NewRedisCache(conf)
	} else {
		return NewNoOpCache()
	}
}
