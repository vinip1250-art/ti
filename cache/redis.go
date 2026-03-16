package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultExpiration      = 24 * time.Hour * 7 // 7 days
	IndexerComandoTorrents = "indexer:comando_torrents"
)

type Redis struct {
	client            *redis.Client
	defaultExpiration time.Duration
}

func NewRedis() *Redis {
	// REDIS_URL takes priority (Render, Railway, Heroku, etc. provide this format)
	// e.g. redis://:password@host:6379 or rediss://:password@host:6379
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err == nil {
			return &Redis{
				client:            redis.NewClient(opts),
				defaultExpiration: DefaultExpiration,
			}
		}
	}

	// Fallback: individual REDIS_HOST + REDIS_PASSWORD env vars
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	// redisPassword can be empty when the server has no authentication enabled
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:6379", redisHost),
			Password: redisPassword,
		}),
		defaultExpiration: DefaultExpiration,
	}
}

func (r *Redis) SetDefaultExpiration(expiration time.Duration) {
	r.defaultExpiration = expiration
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, r.defaultExpiration).Err()
}

func (r *Redis) SetWithExpiration(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
