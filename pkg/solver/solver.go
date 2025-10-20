package solver

import (
	"context"
	"errors"
	"io"
	"sort"
	"sync"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrTimeout      = errors.New("solver timeout")
	ErrUnknownPart  = errors.New("unknown part")
)

type PuzzleSolver interface {
	Init(reader io.Reader) error
	Solve(part int) (string, error)
}

type PuzzleSolverWithContext interface {
	Init(ctx context.Context, reader io.Reader) error
	Solve(ctx context.Context, part int) (string, error)
}

type SolverAdapter struct {
	base PuzzleSolver
}

func NewSolverAdapter(base PuzzleSolver) *SolverAdapter {
	return &SolverAdapter{base: base}
}
func (s SolverAdapter) Solve(ctx context.Context, part int) (string, error) {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		r, err := s.base.Solve(part)
		resCh <- r
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return "", ErrTimeout
	case r := <-resCh:
		return r, <-errCh
	}
}

func (s SolverAdapter) Init(ctx context.Context, reader io.Reader) error {
	errCh := make(chan error)

	go func() {
		err := s.base.Init(reader)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ErrTimeout
	case err := <-errCh:
		return err
	}
}

type Stepper interface {
	Next() (string, error)
}

type RegistryItem struct {
	Name        string
	Next        bool
	Constructor func() PuzzleSolver
}
type RegistryItemPublic struct {
	Name string `json:"name"`
	Next bool   `json:"next"`
} //@name RegistryItem

var registry = map[string]RegistryItem{}
var keys []string
var mu sync.RWMutex

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
	sort.Strings(keys)
}

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

func New(name string) (PuzzleSolver, bool) {
	solver, ok := registry[name]

	if !ok {
		return nil, false
	}

	return solver.Constructor(), true
}
