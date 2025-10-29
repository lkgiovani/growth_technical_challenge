package cache

import (
	"context"
	"encoding/json"
	"fmt"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, keys ...string) error
	DelPattern(ctx context.Context, pattern string) error
}

func GetJSON(ctx context.Context, cache Cache, key string, dest interface{}) error {
	data, err := cache.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func SetJSON(ctx context.Context, cache Cache, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return cache.Set(ctx, key, data)
}
