package ratelimiter

import "time"

type StubRateLimiter struct{}

func NewStubRateLimiter() *StubRateLimiter {
	return &StubRateLimiter{}
}

func (s *StubRateLimiter) Allow(key string, duration time.Duration) (bool, error) {
	return true, nil
}
