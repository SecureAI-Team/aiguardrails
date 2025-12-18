package usage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter enforces per-key quotas using Redis.
type RateLimiter struct {
	redis *redis.Client
	ns    string
}

// NewRateLimiter constructs a RateLimiter.
func NewRateLimiter(redis *redis.Client, namespace string) *RateLimiter {
	return &RateLimiter{redis: redis, ns: namespace}
}

// Allow increments a counter and checks if within quota per hour.
func (r *RateLimiter) Allow(ctx context.Context, key string, quotaPerHour int64) (bool, int64, error) {
	if quotaPerHour <= 0 {
		return true, 0, nil
	}
	now := time.Now().UTC()
	window := now.Format("2006010215")
	prefix := "quota"
	if r.ns != "" {
		prefix = r.ns + ":" + prefix
	}
	redisKey := fmt.Sprintf("%s:%s:%s", prefix, key, window)
	pipe := r.redis.TxPipeline()
	incr := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, time.Hour*2)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, 0, err
	}
	count := incr.Val()
	if count > quotaPerHour {
		return false, count, nil
	}
	return true, count, nil
}

