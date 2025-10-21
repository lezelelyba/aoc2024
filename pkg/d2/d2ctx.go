package d2

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

		for i, report := range *p.reports {

			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}
			if report_safe(&report) {
				sum += 1
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for i, report := range *p.reports {
			if i%100000 == 0 {
				select {
				case <-ctx.Done():
					return "", solver.ErrTimeout
				default:
				}
			}

			if report_safe(&report) {
				sum += 1
			} else {
				modified_report := make([]int, len(report))
				for i := range len(report) {
					copy(modified_report, report)
					modified_report := append(modified_report[:i], modified_report[i+1:]...)
					if report_safe(&modified_report) {
						sum += 1
						break
					}
				}
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
