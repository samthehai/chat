package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) *RedisClient {
	client := redis.NewClient(&redis.Options{Addr: redisURL})

	return &RedisClient{client: client}
}

func (c *RedisClient) LPush(key string, values []byte) error {
	if err := c.client.LPush(key, values).Err(); err != nil {
		return fmt.Errorf("redis lpush: %w", err)
	}

	return nil
}

func (c *RedisClient) SAdd(key string, values []byte) error {
	if err := c.client.SAdd(key, values).Err(); err != nil {
		return fmt.Errorf("redis sadd: %w", err)
	}

	return nil
}

func (c *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	cmd := c.client.LRange(key, start, stop)

	res, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange: %w", err)
	}

	return res, nil
}

func (c *RedisClient) SMembers(key string) ([]string, error) {
	cmd := c.client.SMembers(key)

	res, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange: %w", err)
	}

	return res, nil
}
