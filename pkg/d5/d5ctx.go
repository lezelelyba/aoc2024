package d5

import (
	"advent2024/pkg/solver"
	"context"
	"fmt"
	"io"
	"slices"
	"strconv"
)

type PuzzleStructWithCtx struct {
	PuzzleStruct
}

func init() {
	solver.RegisterWithCtx(day, func() solver.PuzzleSolverWithCtx {
		return NewSolverWithCtx()
	})
}

func NewSolverWithCtx() *PuzzleStructWithCtx {
	return &PuzzleStructWithCtx{}
}

func (p *PuzzleStructWithCtx) InitCtx(ctx context.Context, reader io.Reader) error {
	return p.PuzzleStruct.Init(reader)
}

func (p *PuzzleStructWithCtx) SolveCtx(ctx context.Context, part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for i, update := range p.updates {
			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			if slices.IsSortedFunc(update, p.sortFunc()) {
				sum += update[len(update)/2]
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0
		for i, update := range p.updates {
			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			if !slices.IsSortedFunc(update, p.sortFunc()) {
				slices.SortFunc(update, p.sortFunc())
				sum += update[len(update)/2]
			}
		}
		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
