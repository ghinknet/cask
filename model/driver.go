package model

type NamespaceInfo interface {
	Key() string
}

type Driver interface {
	NewAdapter(client any, ns NamespaceInfo) (adapter Adapter, ok bool)
}
