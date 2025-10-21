package d1

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strconv"

	"advent2024/pkg/solver"
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
		input_copy := p.inputCopy()

		slices.Sort(input_copy[0])
		slices.Sort(input_copy[1])

		diff := 0

		for i := 0; i < len(input_copy[0]); i++ {

			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			diff += difference(input_copy[0][i], input_copy[1][i])
		}

		return strconv.Itoa(diff), nil
	case 2:

		h := histogram(p.input[1])

		sim := 0

		for i := 0; i < len(p.input[0]); i++ {
			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			sim += similarity(p.input[0][i], h)
		}

		return strconv.Itoa(sim), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
