// Package provides registry of solvers
package solver

import (
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Package Errors
// Errors returned by the solver can be tested againts these errors
// using errors.Is
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrTimeout      = errors.New("solver timeout")
	ErrUnknownPart  = errors.New("unknown part")
)

// Interface of Puzzle Solver
type PuzzleSolver interface {
	Init(reader io.Reader) error
	Solve(part int) (string, error)
}

// Registered solver
type RegistryItem struct {
	Name        string
	Next        bool
	Constructor func() PuzzleSolver
}

// Registered solver for export purposes
type RegistryItemPublic struct {
	Name string `json:"name"`
	Next bool   `json:"next"`
} //@name RegistryItem

// Interface of Puzzle Solver supporting stepwise solving
type Stepper interface {
	PuzzleSolver
	Next() (string, error)
}

// Registry of solvers
var registry = map[string]RegistryItem{}

// Keys of the solvers, sorted
var keys []string

// Mutex guarding access to registry
var mu sync.RWMutex

// Registers a solver and check for supported interfaces
// Keeps the keys ordered
func Register(name string, constructor func() PuzzleSolver) {
	mu.Lock()
	defer mu.Unlock()

	item := RegistryItem{Name: name, Constructor: constructor}

	var ps PuzzleSolver

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

	registry[name] = item

	// sort the keys
	keys = append(keys, name)
	slices.SortFunc(keys, cmpDays)
}

// Lists registered keys
func ListRegistryItems() []RegistryItemPublic {
	mu.RLock()
	defer mu.RUnlock()

	items := make([]RegistryItemPublic, 0, len(registry))

	for _, k := range keys {
		v := registry[k]
		items = append(items, RegistryItemPublic{Name: v.Name, Next: v.Next})
	}

	return items
}

// Factory for solvers
func New(name string) (PuzzleSolver, bool) {
	solver, ok := registry[name]

	if !ok {
		return nil, false
	}

	return solver.Constructor(), true
}

func cmpDays(this, other string) int {
	tStr := strings.TrimPrefix(this, "d")
	oStr := strings.TrimPrefix(other, "d")

	t, _ := strconv.Atoi(tStr)
	o, _ := strconv.Atoi(oStr)

	if t == o {
		return 0
	} else if t > o {
		return 1
	}

	return -1
}
