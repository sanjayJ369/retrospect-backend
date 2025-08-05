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
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic("invalid Redis URL: " + err.Error())
	}
	client := redis.NewClient(opt)
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
