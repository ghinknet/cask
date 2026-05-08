package cask

import (
	"strings"

	"go.gh.ink/cask/internal/state"
	"go.gh.ink/toolbox/pointer"
)

func (n *Namespace) Namespace(key ...string) *Namespace {
	if n.raw != nil {
		nsCpy := pointer.Copy(n)
		nsCpy.key = append(n.key, key...)

		// Refresh adapter
		adapter, _ := state.Drivers[nsCpy.adapter].NewAdapter(nsCpy.raw, nsCpy)
		nsCpy.Adapter = adapter

		return nsCpy
	}
	return n
}

func (n *Namespace) Key() string {
	return strings.Join(n.key, ":")
}

func (n *Namespace) Raw() any {
	return n.raw
}
