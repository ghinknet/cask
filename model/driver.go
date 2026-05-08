package model

type NamespaceInfo interface {
	Key() string
}

type Driver interface {
	NewClient(client any, ns NamespaceInfo) (adapter Adapter, ok bool)
}
