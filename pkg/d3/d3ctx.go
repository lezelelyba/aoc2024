package d3

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
		for i, entry := range *p.entries {

			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}

			if entry.instruction == "mul" {
				sum += entry.i * entry.j
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0
		sum_enabled := true
		for i, entry := range *p.entries {

			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}

			if entry.instruction == "mul" && sum_enabled {
				sum += entry.i * entry.j
			} else if entry.instruction == "do" {
				sum_enabled = true
			} else if entry.instruction == "don't" {
				sum_enabled = false
			} else {
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
