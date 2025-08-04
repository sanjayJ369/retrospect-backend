package ratelimiter

import "time"

type RateLimiter interface {
	Allow(key string, duration time.Duration) (bool, error)
}
