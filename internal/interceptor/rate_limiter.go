package interceptor

import (
	"context"

	"github.com/moremoneymod/auth/internal/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type RateLimiterInterceptor struct {
	rateLimiter *rate_limiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(rateLimiter *rate_limiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{rateLimiter: rateLimiter}
}

func (interceptor *RateLimiterInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !interceptor.rateLimiter.Allow() {
		return nil, grpc.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}
	return handler(ctx, req)
}
