// Package provides registry of solvers with context support
package solver

import (
	"context"
	"io"
	"slices"
	"sync"
)

// Interface of Puzzle Solver with context support
type PuzzleSolverWithCtx interface {
	PuzzleSolver
	InitCtx(ctx context.Context, reader io.Reader) error
	SolveCtx(ctx context.Context, part int) (string, error)
}

// Registered solver with context support
type RegistryItemWithCtx struct {
	Name        string
	Next        bool
	Constructor func() PuzzleSolverWithCtx
}

// Interface of Puzzle Solver with support for context and stepwise solving
type StepperWithCtx interface {
	PuzzleSolverWithCtx
	Next(ctx context.Context) (string, error)
}

// Registry of solvers with context support
var registryCtx = map[string]RegistryItemWithCtx{}

// Keys of the solvers with context support, sorted
var keysCtx []string

// Mutex guarding access to registry of solvers with context support
var muCtx sync.RWMutex

// Registers a solver and check for supported interfaces
// Keeps the keys ordered
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
	slices.SortFunc(keysCtx, cmpDays)
}

// Lists registered keys of solvers with context support
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

// Factory for solvers with context
func NewWithCtx(name string) (PuzzleSolverWithCtx, bool) {
	solver, ok := registryCtx[name]

	if !ok {
		return nil, false
	}

	return solver.Constructor(), true
}
