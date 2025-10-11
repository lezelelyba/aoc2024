package d9

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

type Disk struct {
}
type PuzzleStruct struct {
	input    string
	intInput []int
	blocks   map[int]int
}

func init() {
	solver.Register("d9", func() solver.PuzzleSolver {
		return NewSolver()
	})
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	input, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	p.input = input
	intInput := make([]int, len(input))

	for i, c := range input {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			log.Print(err)
			return err
		}
		intInput[i] = n
	}

	p.intInput = intInput

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		front_array_idx := 0
		front_block_idx := 0
		back_array_idx := len(p.intInput) - 1
		back_block_idx := len(p.intInput) / 2

		space := false

	outer:
		// go through the disk by index
		for disk_idx := 0; ; disk_idx++ {
			// until front is empty
			for p.intInput[front_array_idx] == 0 {
				// move to next block
				front_array_idx++

				// if block is out of bounds break
				if front_array_idx >= len(p.intInput) {
					break outer
				}

				// block > space, space > block
				space = !space

				// if next block file, increase the file index
				if !space {
					front_block_idx++
				}
			}

			if !space {
				// increase sum
				sum += disk_idx * front_block_idx

				// lower the front block count
				p.intInput[front_array_idx] = p.intInput[front_array_idx] - 1
				// if space
			} else {
				// increase the sum by the back block
				sum += disk_idx * back_block_idx

				// lower the back block count
				p.intInput[back_array_idx] = p.intInput[back_array_idx] - 1
				// lower the front space count
				p.intInput[front_array_idx] = p.intInput[front_array_idx] - 1

				// if we are out of block at the back move to next back block
				if p.intInput[back_array_idx] == 0 {
					// block < space < block
					back_array_idx -= 2
					// block id -1
					back_block_idx -= 1
				}
			}

			// if we moved back array past the front array, break
			if front_array_idx > back_array_idx {
				break
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("unknown Part %d", part)
}

func parseInput(sc *bufio.Scanner) (string, error) {

	var line string

	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())
	}

	return line, nil
}
