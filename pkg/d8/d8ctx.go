package d8

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

		antinodesMap := make(map[Coords]struct{})

		for freq := range p.antennas {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}

			antennas := p.antennas[freq]
			for a1 := 0; a1 < len(antennas); a1++ {
				for a2 := a1 + 1; a2 < len(antennas); a2++ {
					antinodes := FindAntinodes(antennas[a1], antennas[a2])

					for _, an := range antinodes {
						if p.inField(an) {
							antinodesMap[an] = struct{}{}
						}
					}
				}
			}
		}

		sum = len(antinodesMap)

		return strconv.Itoa(sum), nil
	case 2:

		sum := 0

		antinodesMap := make(map[Coords]struct{})

		for freq := range p.antennas {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}

			antennas := p.antennas[freq]
			for a1 := 0; a1 < len(antennas); a1++ {
				for a2 := a1 + 1; a2 < len(antennas); a2++ {

					// line in one direction
					iter := NewLineIter(antennas[a1], antennas[a2])

					for an, ok := iter.Next(); ok; an, ok = iter.Next() {
						if !p.inField(an) {
							break
						}
						antinodesMap[an] = struct{}{}
					}

					// line in the other direction
					iter = NewLineIter(antennas[a2], antennas[a1])

					for an, ok := iter.Next(); ok; an, ok = iter.Next() {
						if !p.inField(an) {
							break
						}
						antinodesMap[an] = struct{}{}
					}
				}
			}
		}

		sum = len(antinodesMap)

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
