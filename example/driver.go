package example

import (
	"go.gh.ink/cask/model"
)

type Driver struct{}

func (d Driver) NewAdapter(client any, ns model.NamespaceInfo) (adapter model.Adapter, ok bool) {
	if _, ok = client.(bool); ok {
		return Adapter{}, true
	}
	return nil, false
}
