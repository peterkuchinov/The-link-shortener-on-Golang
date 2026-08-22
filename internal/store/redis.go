package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(redisURL string) (*RedisStore, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	Client := redis.NewClient(opts)

	if err := Client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisStore{client: Client}, nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func (s *RedisStore) buildKey(code string) string {
	return fmt.Sprintf("link:%s", code)
}

func (s *RedisStore) Get(ctx context.Context, code string) (string, error) {
	return s.client.Get(ctx, s.buildKey(code)).Result()
}

func (s *RedisStore) Set(ctx context.Context, code string, url string, ttl time.Duration) error {
	return s.client.Set(ctx, s.buildKey(code), url, ttl).Err()
}

func (s *RedisStore) Delete(ctx context.Context, code string) error {
	return s.client.Del(ctx, s.buildKey(code)).Err()
}
