package entity

import (
	"context"
	"time"
)

type RateLimiterRepository interface {
	AddKey(context.Context, string, time.Duration) error
	Exists(context.Context, string) (int64, error)
	Increment(context.Context, string) (int64, error)
	Expire(context.Context, string, time.Duration) error
	Delete(context.Context, string) error
	AddHash(context.Context, string, int, time.Duration) error
	Find(context.Context, string) (bool, int, time.Duration, error)
}
