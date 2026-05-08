package model

import (
	"context"

	"go.gh.ink/timex"
)

type BaseStore interface {
	Get(ctx context.Context) ([]byte, error)
	Set(ctx context.Context, value []byte, ttl timex.Duration) error
	Del(ctx context.Context) (bool, error)
	Exists(ctx context.Context) (bool, error)
}

type Expirer interface {
	Expire(ctx context.Context, ttl timex.Duration) error
	TTL(ctx context.Context) (timex.Duration, error)
}

type Adapter interface {
	BaseStore
}
