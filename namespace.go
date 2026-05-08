package cask

import (
	"strings"

	"go.gh.ink/toolbox/pointer"
)

func (n *Namespace) Namespace(key ...string) *Namespace {
	if n.raw != nil {
		nsCpy := pointer.Copy(n)
		nsCpy.key = append(n.key, key...)
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
