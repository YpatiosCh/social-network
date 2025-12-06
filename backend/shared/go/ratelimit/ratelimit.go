package ratelimit

import (
	"context"
	"errors"
	"time"
)

var ErrBadValue = errors.New("bad value received")
var ErrNotANumber = errors.New("value stored not a number")
var ErrStorageProblem = errors.New("problem with storage mechanism")

// what the rate limiter will be using to save its ratelimit entries
type storage interface {
	IncrEx(ctx context.Context, key string, duration time.Duration) (currentCount int, err error)
}

type ratelimiter struct {
	globalPrefix string //will be used as a prefix on all keys, when the regular save function is called
	storage      storage
}

func NewRateLimiter(globalPrefix string, storage storage) (ratelimiter, error) {
	rateLimiter := ratelimiter{
		globalPrefix: globalPrefix,
		storage:      storage,
	}
	return rateLimiter, nil
}

// Allow checks if the action identified by `key` is allowed under the rate limit defined by `limit` and `duration`. It will prefix the key with the global prefix.
func (rl *ratelimiter) Allow(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	return rl.AllowRawKey(ctx, rl.globalPrefix+key, limit, duration)
}

// AllowRawKey checks if the action identified by `key` is allowed under the rate limit defined by `limit` and `duration`.
func (rl *ratelimiter) AllowRawKey(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	storageLimit, err := rl.storage.IncrEx(ctx, key, duration)
	if err != nil {
		return false, errors.Join(ErrStorageProblem, err)
	}
	if storageLimit > limit {
		return false, nil
	}
	return true, nil
}
