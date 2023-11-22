package stubs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type StubRedisClient struct{}

func (c *StubRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringResult("", nil)
}

func (c *StubRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("OK", nil)
}
