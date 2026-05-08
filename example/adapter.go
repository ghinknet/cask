package example

import (
	"context"

	"go.gh.ink/timex"
)

type Adapter struct{}

func (a Adapter) Get(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (a Adapter) Set(ctx context.Context, value []byte, ttl timex.Duration) error {
	return nil
}

func (a Adapter) Del(ctx context.Context) (bool, error) {
	return false, nil
}

func (a Adapter) Exists(ctx context.Context) (bool, error) {
	return false, nil
}
