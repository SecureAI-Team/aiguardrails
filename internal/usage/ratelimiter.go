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
}

// NewRateLimiter constructs a RateLimiter.
func NewRateLimiter(redis *redis.Client) *RateLimiter {
	return &RateLimiter{redis: redis}
}

// Allow increments a counter and checks if within quota per hour.
func (r *RateLimiter) Allow(ctx context.Context, key string, quotaPerHour int64) (bool, int64, error) {
	if quotaPerHour <= 0 {
		return true, 0, nil
	}
	now := time.Now().UTC()
	window := now.Format("2006010215")
	redisKey := fmt.Sprintf("quota:%s:%s", key, window)
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

