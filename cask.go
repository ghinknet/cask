package cask

import (
	"go.gh.ink/cask/errors"
	"go.gh.ink/cask/internal/state"
)

func New(client any, key ...string) (n *Namespace, err error) {
	ns := &Namespace{
		key: key,
		raw: client,
	}

	// Try to find a suit client
	for _, v := range state.Drivers {
		if adapter, ok := v.NewClient(client, ns); ok {
			ns.Adapter = adapter
		}
	}
	if ns.Adapter == nil {
		return new(Namespace), errors.ErrDriverNotRegistered
	}

	return ns, nil
}
