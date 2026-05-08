package cask

import (
	"context"
	"testing"

	"go.gh.ink/cask/internal/state"
	"go.gh.ink/cask/model"
	"go.gh.ink/timex"
)

type testAdapter struct {
	key string
}

func (a *testAdapter) Get(ctx context.Context) ([]byte, error) { return []byte(a.key), nil }
func (a *testAdapter) Set(ctx context.Context, value []byte, ttl timex.Duration) error {
	return nil
}
func (a *testAdapter) Del(ctx context.Context) (bool, error) { return true, nil }
func (a *testAdapter) Exists(ctx context.Context) (bool, error) { return true, nil }

type testDriver struct {
	calls      int
	lastKey    string
	lastClient any
}

func (d *testDriver) NewAdapter(client any, ns model.NamespaceInfo) (model.Adapter, bool) {
	d.calls++
	d.lastKey = ns.Key()
	d.lastClient = client
	return &testAdapter{key: ns.Key()}, true
}

func TestNamespaceRefreshesAdapterOnChain(t *testing.T) {
	oldDrivers := state.Drivers
	state.Drivers = map[string]model.Driver{}
	t.Cleanup(func() { state.Drivers = oldDrivers })

	drv := &testDriver{}
	state.Drivers["test"] = drv

	ns := &Namespace{
		key:     []string{"root"},
		raw:     "client",
		adapter: "test",
		Adapter: &testAdapter{key: "root"},
	}

	ns2 := ns.Namespace("child")
	if ns2 == ns {
		t.Fatalf("expected a new Namespace instance")
	}
	if ns.Key() != "root" {
		t.Fatalf("expected original key to remain 'root', got %q", ns.Key())
	}
	if ns2.Key() != "root:child" {
		t.Fatalf("expected namespaced key 'root:child', got %q", ns2.Key())
	}
	if drv.calls != 1 {
		t.Fatalf("expected driver to be called once, got %d", drv.calls)
	}
	if drv.lastKey != "root:child" {
		t.Fatalf("expected driver to receive key 'root:child', got %q", drv.lastKey)
	}
	if drv.lastClient != "client" {
		t.Fatalf("expected driver to receive original client")
	}

	adapter, ok := ns2.Adapter.(*testAdapter)
	if !ok {
		t.Fatalf("expected adapter type *testAdapter")
	}
	if adapter.key != "root:child" {
		t.Fatalf("expected adapter key 'root:child', got %q", adapter.key)
	}
}

func TestNamespaceReturnsSelfWhenRawNil(t *testing.T) {
	ns := &Namespace{key: []string{"root"}}
	ns2 := ns.Namespace("child")
	if ns2 != ns {
		t.Fatalf("expected same Namespace instance when raw is nil")
	}
}

func TestKeyAndRaw(t *testing.T) {
	ns := &Namespace{key: []string{"a", "b"}, raw: 123}
	if ns.Key() != "a:b" {
		t.Fatalf("expected key 'a:b', got %q", ns.Key())
	}
	if ns.Raw() != 123 {
		t.Fatalf("expected raw to return original value")
	}

	empty := &Namespace{}
	if empty.Key() != "" {
		t.Fatalf("expected empty key to be empty string, got %q", empty.Key())
	}
}

func TestNamespaceChainDoesNotMutateParent(t *testing.T) {
	oldDrivers := state.Drivers
	state.Drivers = map[string]model.Driver{}
	t.Cleanup(func() { state.Drivers = oldDrivers })

	drv := &testDriver{}
	state.Drivers["test"] = drv

	ns := &Namespace{
		key:     []string{"root"},
		raw:     "client",
		adapter: "test",
		Adapter: &testAdapter{key: "root"},
	}

	ns1 := ns.Namespace("a")
	ns2 := ns1.Namespace("b")

	if ns1.Key() != "root:a" {
		t.Fatalf("expected ns1 key 'root:a', got %q", ns1.Key())
	}
	if ns2.Key() != "root:a:b" {
		t.Fatalf("expected ns2 key 'root:a:b', got %q", ns2.Key())
	}
	if ns.Key() != "root" {
		t.Fatalf("expected original key to remain 'root', got %q", ns.Key())
	}
}

func TestNamespaceDoesNotMutateParentWhenSliceHasCapacity(t *testing.T) {
	oldDrivers := state.Drivers
	state.Drivers = map[string]model.Driver{}
	t.Cleanup(func() { state.Drivers = oldDrivers })

	drv := &testDriver{}
	state.Drivers["test"] = drv

	base := make([]string, 1, 4)
	base[0] = "root"

	ns := &Namespace{
		key:     base,
		raw:     "client",
		adapter: "test",
		Adapter: &testAdapter{key: "root"},
	}

	_ = ns.Namespace("child")

	if len(ns.key) != 1 || ns.key[0] != "root" {
		t.Fatalf("expected parent key to remain [root], got %v", ns.key)
	}
}
