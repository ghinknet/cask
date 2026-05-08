package model

import (
	"context"

	"go.gh.ink/timex"
)

type BaseStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl timex.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
}

type Expirer interface {
	Expire(ctx context.Context, key string, ttl timex.Duration) error
	TTL(ctx context.Context, key string) (timex.Duration, error)
}

type Adapter interface {
	BaseStore
}
