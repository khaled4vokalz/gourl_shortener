package cache

import (
	"context"
	"fmt"
	"time"

	errors "github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(conf config.CacheConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		DB:   conf.Database,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}

	return &RedisCache{Client: client}, nil
}

func (r *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	item, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.GetLogger().Debug(fmt.Sprintf("No url found in cache for key '%s'", key))
		return "", errors.NotFound
	} else {
		return item, err
	}
}
