package cask

import (
	"strings"
)

func (n Namespace) Namespace(key ...string) Namespace {
	if n.raw != nil {
		n.key = append(n.key, key...)
	}
	return n
}

func (n Namespace) Key() string {
	return strings.Join(n.key, ":")
}
