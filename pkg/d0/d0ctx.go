package d0

import (
	"advent2024/pkg/solver"
	"context"
	"fmt"
	"io"
	"strconv"
)

// PuzzleStruct with Context
type PuzzleStructWithCtx struct {
	PuzzleStruct
}

// Registers day wih the registry
func init() {
	solver.RegisterWithCtx(day, func() solver.PuzzleSolverWithCtx {
		return NewSolverWithCtx()
	})
}

// Constructor
func NewSolverWithCtx() *PuzzleStructWithCtx {
	return &PuzzleStructWithCtx{}
}

// Initializes the PuzzleStruct with input
func (p *PuzzleStructWithCtx) InitCtx(ctx context.Context, reader io.Reader) error {
	return p.PuzzleStruct.Init(reader)
}

// Solves the puzzle
// Accepts part as parameter
// Returns string containing the solution of the puzzle
func (p *PuzzleStructWithCtx) SolveCtx(ctx context.Context, part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		select {
		case <-ctx.Done():
			return "", solver.ErrTimeout
		default:
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		select {
		case <-ctx.Done():
			return "", solver.ErrTimeout
		default:
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
