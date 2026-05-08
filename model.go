package cask

import "go.gh.ink/cask/model"

type Namespace struct {
	key []string
	raw any
	model.Adapter
}
