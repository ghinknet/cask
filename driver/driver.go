package driver

import (
	"maps"
	"slices"

	"go.gh.ink/cask/internal/state"
	"go.gh.ink/cask/model"
)

func Register(name string, driver model.Driver) {
	state.Drivers[name] = driver
}

func List() []string {
	return slices.Collect(maps.Keys(state.Drivers))
}
