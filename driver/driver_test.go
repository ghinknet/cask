package driver

import (
	"testing"

	"go.gh.ink/cask/internal/state"
	"go.gh.ink/cask/model"
)

type stubDriver struct{}

func (d stubDriver) NewAdapter(client any, ns model.NamespaceInfo) (model.Adapter, bool) {
	return nil, true
}

func TestRegisterAndList(t *testing.T) {
	oldDrivers := state.Drivers
	state.Drivers = map[string]model.Driver{}
	t.Cleanup(func() { state.Drivers = oldDrivers })

	Register("alpha", stubDriver{})
	Register("beta", stubDriver{})

	list := List()
	if len(list) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(list))
	}

	found := map[string]bool{}
	for _, name := range list {
		found[name] = true
	}
	if !found["alpha"] || !found["beta"] {
		t.Fatalf("expected list to contain 'alpha' and 'beta', got %v", list)
	}

	Register("alpha", stubDriver{})
	if len(state.Drivers) != 2 {
		t.Fatalf("expected register to overwrite without changing count")
	}
}

