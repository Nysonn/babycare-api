package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService wraps a Redis client for JSON-serialised get/set/delete operations.
type CacheService struct {
	client *redis.Client
}

// NewCacheService parses redisURL, creates a client, and verifies connectivity.
// Returns an error if the connection cannot be established.
func NewCacheService(redisURL string) (*CacheService, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("cache: parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("cache: ping redis: %w", err)
	}

	return &CacheService{client: client}, nil
}

// Get returns the raw string stored at key, or an error if the key does not exist.
func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Set marshals value to JSON and stores it under key with the given TTL.
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: marshal value: %w", err)
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes one or more keys from the cache.
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	return s.client.Del(ctx, keys...).Err()
}
