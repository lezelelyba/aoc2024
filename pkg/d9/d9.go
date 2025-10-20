package d9

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

var day = "d9"

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type Block struct {
	size    int
	blockid int
	space   bool
	moved   bool
}

type PuzzleStruct struct {
	inputStr  *string
	inputInts *[]int

	blockList list.List
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	inputStr, inputInts, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	p.inputStr = inputStr
	p.inputInts = inputInts

	if err := validateInput(inputInts); err != nil {
		log.Print(err)
		return err
	}

	space := false
	blockid := 0

	for i := 0; i < len(*p.inputInts); i++ {
		if !space {
			p.blockList.PushBack(&Block{size: (*p.inputInts)[i],
				blockid: blockid,
				space:   false,
				moved:   false})
		} else {
			p.blockList.PushBack(&Block{size: (*p.inputInts)[i],
				blockid: 0,
				space:   true,
				moved:   false})
		}

		space = !space
		if !space {
			blockid++
		}
	}

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
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

func parseInput(sc *bufio.Scanner) (*string, *[]int, error) {

	var line string
	var inputInts []int

	// expecting only 1 line
	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())

		inputInts = make([]int, len(line))

		for i, c := range line {
			n, err := strconv.Atoi(string(c))
			if err != nil {
				log.Print(err)
				return nil, nil, fmt.Errorf("%s unable to convert %s to int: %w", day, string(c), solver.ErrInvalidInput)
			}
			inputInts[i] = n
		}
	}

	if sc.Err() != nil {
		return nil, nil, fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return &line, &inputInts, nil
}

func validateInput(ints *[]int) error {
	if ints == nil {
		return fmt.Errorf("%s empty rules: %w", day, solver.ErrInvalidInput)
	} else if len(*ints) == 0 {
		return fmt.Errorf("%s empty rules: %w", day, solver.ErrInvalidInput)
	} else if len(*ints)%2 != 1 {
		return fmt.Errorf("%s non odd input length: %w", day, solver.ErrInvalidInput)
	}

	return nil
}

func tryToMove(l *list.List, e *list.Element) (*list.Element, error) {

	head := l.Front()

	if e.Value.(*Block).space {
		return e, fmt.Errorf("trying to move space")
	}

	if e.Value.(*Block).moved {
		return e, fmt.Errorf("Already moved")
	}

	for next := head; next != e; next = next.Next() {
		if !next.Value.(*Block).space {
			continue
		}

		if next.Value.(*Block).size < e.Value.(*Block).size {
			continue
		}

		new_space := Block{space: true, size: e.Value.(*Block).size, moved: false, blockid: 0}

		ret := l.InsertAfter(&new_space, e)
		l.MoveBefore(e, next)

		next.Value.(*Block).size -= new_space.size
		if next.Value.(*Block).size == 0 {
			_ = l.Remove(next)
		}

		e.Value.(*Block).moved = true

		return ret, nil
	}

	return e, fmt.Errorf("Cannot move")
}

func Checksum(l *list.List) int {

	sum := 0
	idx := 0

	for c := l.Front(); c != l.Back(); c = c.Next() {
		if c.Value.(*Block).space {
			idx += c.Value.(*Block).size
		} else {
			blockid := c.Value.(*Block).blockid
			for i := idx; i < c.Value.(*Block).size+idx; i++ {
				sum += i * blockid
			}
			idx += c.Value.(*Block).size
		}
	}

	return sum
}
