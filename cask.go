package cask

import (
	"go.gh.ink/cask/errors"
	"go.gh.ink/cask/internal/state"
	"go.gh.ink/cask/model"
)

func New(client any, key ...string) (n Namespace, err error) {
	// Try to find a suit client
	var ada model.Adapter
	for _, v := range state.Drivers {
		if adapter, ok := v.NewClient(client); ok {
			ada = adapter
		}
	}
	if ada == nil {
		return Namespace{}, errors.ErrDriverNotRegistered
	}

	return Namespace{
		key:     key,
		raw:     client,
		Adapter: ada,
	}, nil
}
