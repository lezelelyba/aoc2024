package solver

import (
	"context"
	"io"
	"sort"
	"sync"
)

type PuzzleSolverWithCtx interface {
	PuzzleSolver
	InitCtx(ctx context.Context, reader io.Reader) error
	SolveCtx(ctx context.Context, part int) (string, error)
}

type RegistryItemWithCtx struct {
	Name        string
	Next        bool
	Constructor func() PuzzleSolverWithCtx
}

type StepperWithCtx interface {
	PuzzleSolverWithCtx
	Next(ctx context.Context) (string, error)
}

var registryCtx = map[string]RegistryItemWithCtx{}
var keysCtx []string
var muCtx sync.RWMutex

func RegisterWithCtx(name string, constructor func() PuzzleSolverWithCtx) {
	muCtx.Lock()
	defer muCtx.Unlock()

	item := RegistryItemWithCtx{Name: name, Constructor: constructor}

	var ps PuzzleSolverWithCtx

	// recovers from panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				ps = nil
			}
		}()
		ps = constructor()
	}()

	if _, ok := ps.(Stepper); ok {
		item.Next = true
	}

	registryCtx[name] = item

	// sort the keys
	keysCtx = append(keysCtx, name)
	sort.Strings(keysCtx)
}

func ListRegistryItemsWithCtx() []RegistryItemPublic {
	muCtx.RLock()
	defer muCtx.RUnlock()

	items := make([]RegistryItemPublic, 0, len(registryCtx))

	for _, k := range keysCtx {
		v := registryCtx[k]
		items = append(items, RegistryItemPublic{Name: v.Name, Next: v.Next})
	}

	return items
}

func NewWithCtx(name string) (PuzzleSolverWithCtx, bool) {
	solver, ok := registryCtx[name]

	if !ok {
		return nil, false
	}

	return solver.Constructor(), true
}
