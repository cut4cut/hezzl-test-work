package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	_defaultAddr        = "redis:6379"
	_defaultPassword    = ""
	_defaultDB          = 0
	_defaultExpiratione = time.Duration(60 * 1e9) // 60s
)

type RedisClient struct {
	*redis.Client
	ExpirationeSecond time.Duration
}

func New(options ...int) (*RedisClient, error) {
	addr := _defaultAddr
	expir := _defaultExpiratione

	switch len(options) {
	case 1:
		addr = fmt.Sprintf("redis:%d", options[0])
	case 2:
		addr = fmt.Sprintf("redis:%d", options[0])
		expir = time.Duration(options[1] * 1e9)
	default:
		addr = _defaultAddr
		expir = _defaultExpiratione
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: _defaultPassword,
		DB:       _defaultDB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis client can't connect, error: %w", err)
	}

	return &RedisClient{client, expir}, nil
}

func (c *RedisClient) SetValue(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, p, c.ExpirationeSecond).Err()
}

func (c *RedisClient) GetValue(ctx context.Context, key string, dest interface{}) error {
	p, err := c.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(p, dest)
}
