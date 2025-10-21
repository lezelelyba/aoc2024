package d6

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

		for i := 0; p.guard.Move(&p.field) == nil; i++ {
			if i%1000000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
		}

		sum = len(p.guard.visited)

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		og := NewGuard(p.guard.c.x, p.guard.c.y, p.guard.o)

		for i := 0; p.guard.Move(&p.field) == nil; i++ {
			if i%1000000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
		}

		visited := p.guard.visited

		for coord := range visited {
			// get original guard
			p.guard = NewGuard(og.c.x, og.c.y, og.o)

			// skip initial field
			if p.guard.c == coord {
				continue
			}

			// put obstacle in place
			p.field[coord.y][coord.x] = '#'

			var err error
			// loop
			for err = p.guard.Move(&p.field); err == nil; err = p.guard.Move(&p.field) {
			}

			switch err.(type) {
			case LoopingError:
				sum += 1
			default:
			}

			// remove obstacle
			p.field[coord.y][coord.x] = '.'
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
