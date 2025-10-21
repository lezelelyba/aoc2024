package d4

import (
	"advent2024/pkg/solver"
	"context"
	"fmt"
	"io"
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

		for y := 0; y < p.dy; y++ {

			if y == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}

			for x := 0; x < p.dx; x++ {
				sum += p.xmas(x, y)
			}
		}

		return strconv.Itoa(sum), nil

	case 2:
		sum := 0

		for y := 0; y < p.dy; y++ {
			if y == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			for x := 0; x < p.dx; x++ {
				if p.xmasPart2(x, y) {
					sum += 1
				}
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
