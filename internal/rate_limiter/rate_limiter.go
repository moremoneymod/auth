package rate_limiter

import (
	"context"
	"fmt"
	"time"
)

type TokenBucketLimiter struct {
	tokenBucket chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{make(chan struct{}, limit)}

	for i := 0; i < limit; i++ {
		limiter.tokenBucket <- struct{}{}
	}

	replenishmentInterval := period / time.Duration(limit)
	go limiter.startPeriodReplenishment(ctx, replenishmentInterval)

	return limiter
}

func (l *TokenBucketLimiter) startPeriodReplenishment(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.tokenBucket <- struct{}{}
		case <-ctx.Done():
			return
		}
	}
}

func (l *TokenBucketLimiter) Allow() bool {
	select {
	case <-l.tokenBucket:
		fmt.Println(1)
		return true
	default:
		return false
	}
}
