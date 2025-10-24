package d9

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

		front_array_idx := 0
		front_block_idx := 0
		back_array_idx := len(*p.inputInts) - 1
		back_block_idx := len(*p.inputInts) / 2

		space := false

	outer:
		// go through the disk by index
		for disk_idx := 0; ; disk_idx++ {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}
			// until front is empty
			for (*p.inputInts)[front_array_idx] == 0 {
				// move to next block
				front_array_idx++

				// if block is out of bounds break
				if front_array_idx >= len(*p.inputInts) {
					break outer
				}

				// block > space, space > block
				space = !space

				// if next block file, increase the file index
				if !space {
					front_block_idx++
				}
			}

			// if we moved back array past the front array, break
			if front_array_idx > back_array_idx {
				break
			}

			if !space {
				// increase sum
				sum += disk_idx * front_block_idx

				// lower the front block count
				(*p.inputInts)[front_array_idx] = (*p.inputInts)[front_array_idx] - 1
				// if space
			} else {
				// increase the sum by the back block
				sum += disk_idx * back_block_idx

				// lower the back block count
				(*p.inputInts)[back_array_idx] = (*p.inputInts)[back_array_idx] - 1
				// lower the front space count
				(*p.inputInts)[front_array_idx] = (*p.inputInts)[front_array_idx] - 1

				// if we are out of block at the back move to next back block
				if (*p.inputInts)[back_array_idx] == 0 {
					// block < space < block
					back_array_idx -= 2
					// block id -1
					back_block_idx -= 1
				}
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for back := p.blockList.Back(); back != p.blockList.Front(); back = back.Prev() {

			select {
			case <-ctx.Done():
				return "", solver.ErrTimeout
			default:
			}

			block := back.Value.(*Block)

			if block.space {
				continue
			}

			back, _ = tryToMove(&p.blockList, back)
		}

		sum = Checksum(&p.blockList)

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}
