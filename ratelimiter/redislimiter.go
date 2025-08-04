package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(redisURL string) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	return &RedisRateLimiter{
		client: client,
	}
}

func (r *RedisRateLimiter) Allow(key string, duration time.Duration) (bool, error) {
	ok, err := r.client.SetNX(context.Background(), key, 1, duration).Result()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}
