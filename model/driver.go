package model

type Driver interface {
	NewClient(client any) (adapter Adapter, ok bool)
}
