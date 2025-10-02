package solver

import (
	"io"
	"strings"
)

type PuzzleSolver interface {
	Init(reader io.Reader) error
	Solve(part int) (string, error)
}

var registry = map[string]func() PuzzleSolver{}

func Register(name string, constructor func() PuzzleSolver) {
	registry[name] = constructor
}

func ListRegister() []string {
	registered_keys := make([]string, len(registry))

	i := 0
	for k := range registry {
		registered_keys[i] = strings.Clone(k)
		i++
	}

	return registered_keys
}

func New(name string) (PuzzleSolver, bool) {
	constructor, ok := registry[name]

	if !ok {
		return nil, false
	}

	return constructor(), true
}
