package example

import (
	"context"

	"go.gh.ink/timex"
)

type Adapter struct{}

func (a Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (a Adapter) Set(ctx context.Context, key string, value []byte, ttl timex.Duration) error {
	return nil
}

func (a Adapter) Del(ctx context.Context, keys ...string) (int64, error) {
	return 0, nil
}

func (a Adapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	return 0, nil
}
