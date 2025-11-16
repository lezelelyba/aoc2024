package d12

import (
	"advent2024/pkg/solver"
	"container/list"
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

	// preprocess - determine regions

	regionCnt := NewCounter()

	for y := 0; y < len(*p.field); y++ {
		for x := 0; x < len((*p.field)[0]); x++ {
			if visited(p.field, Coord{x: x, y: y}) {
				continue
			}

			regionCnt.Next()

			(*p.field)[y][x].Region = regionCnt.Value()

			stack := list.New()
			stack.PushFront(Coord{x: x, y: y})

			for stack.Len() > 0 {
				cur := stack.Front().Value.(Coord)
				stack.Remove(stack.Front())

				neighbors := filterNeighbors(p.field, Coord{x: cur.x, y: cur.y}, sameByte)

				for _, n := range neighbors {
					if !visited(p.field, n) {
						(*p.field)[n.y][n.x].Region = regionCnt.Value()
						stack.PushFront(n)
					}
				}
			}
		}
	}

	select {
	case <-ctx.Done():
		return "", solver.ErrTimeout
	default:
	}

	switch part {
	case 1:
		sum := 0

		var regions = make(map[int]Region)

		for y := 0; y < len(*p.field); y++ {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}

			for x := 0; x < len((*p.field)[0]); x++ {
				id := (*p.field)[y][x].Region

				if _, ok := regions[id]; !ok {
					regions[id] = Region{Area: 0, Perimeter: 0, Edges: 0}
				}

				region := regions[id]
				region.Area += 1
				region.Perimeter += borderCount(p.field, Coord{x: x, y: y})

				regions[id] = region
			}
		}

		for _, region := range regions {
			sum += region.Area * region.Perimeter
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		var regions = make(map[int]Region)

		for y := 0; y < len(*p.field); y++ {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}

			for x := 0; x < len((*p.field)[0]); x++ {
				id := (*p.field)[y][x].Region

				if _, ok := regions[id]; !ok {
					regions[id] = Region{Area: 0, Perimeter: 0, Edges: 0}
				}

				region := regions[id]
				region.Area += 1
				region.Edges += newEdges(p.field, Coord{x: x, y: y})

				regions[id] = region
			}
		}

		for _, region := range regions {
			sum += region.Area * region.Edges
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
