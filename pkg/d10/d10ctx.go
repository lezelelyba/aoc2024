package d10

import (
	"advent2024/pkg/solver"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"
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

		for {
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for {
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
