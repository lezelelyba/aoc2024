package d4

import (
	"advent2024/pkg/solver"
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

var day = "d4"

type PuzzleStruct struct {
	dx, dy int
	input  [][]byte
}

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
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

	if err := validateInput(input); err != nil {
		log.Print(err)
		return err
	}

	p.input = *input
	p.dy = len(p.input)
	if p.dy > 0 {
		p.dx = len(p.input[0])
	}

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for y := 0; y < p.dy; y++ {
			for x := 0; x < p.dx; x++ {
				sum += p.xmas(x, y)
			}
		}

		return strconv.Itoa(sum), nil

	case 2:
		sum := 0

		for y := 0; y < p.dy; y++ {
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

func (p *PuzzleStruct) xmas(x, y int) int {

	sum := 0

	if p.is(x, y, 'X') {
		neighbors := neighbors(x, y)
		for _, nm := range neighbors {
			if p.is(nm[0], nm[1], 'M') {
				na := next(x, y, nm[0], nm[1])
				if p.is(na[0], na[1], 'A') {
					ns := next(nm[0], nm[1], na[0], na[1])
					if p.is(ns[0], ns[1], 'S') {
						sum += 1
					}
				}
			}
		}
	}

	return sum
}

func (p *PuzzleStruct) xmasPart2(x, y int) bool {
	if p.is(x, y, 'A') {
		nbs := neighbors(x, y)

		tl := nbs[0]
		tr := nbs[2]
		bl := nbs[5]
		br := nbs[7]

		if (p.is(tl[0], tl[1], 'S')) &&
			(p.is(tr[0], tr[1], 'S')) &&
			(p.is(br[0], br[1], 'M')) &&
			(p.is(bl[0], bl[1], 'M')) {
			return true
		} else if (p.is(tl[0], tl[1], 'M')) &&
			(p.is(tr[0], tr[1], 'S')) &&
			(p.is(br[0], br[1], 'S')) &&
			(p.is(bl[0], bl[1], 'M')) {
			return true
		} else if (p.is(tl[0], tl[1], 'M')) &&
			(p.is(tr[0], tr[1], 'M')) &&
			(p.is(br[0], br[1], 'S')) &&
			(p.is(bl[0], bl[1], 'S')) {
			return true
		} else if (p.is(tl[0], tl[1], 'S')) &&
			(p.is(tr[0], tr[1], 'M')) &&
			(p.is(br[0], br[1], 'M')) &&
			(p.is(bl[0], bl[1], 'S')) {
			return true
		}
	}
	return false
}

func (p *PuzzleStruct) is(x, y int, c byte) bool {
	if x >= 0 && x < p.dx && y >= 0 && y < p.dy {
		return p.input[y][x] == c
	}

	return false
}

func neighbors(x, y int) [][2]int {
	return [][2]int{
		{x - 1, y - 1},
		{x, y - 1},
		{x + 1, y - 1},
		{x - 1, y},
		{x + 1, y},
		{x - 1, y + 1},
		{x, y + 1},
		{x + 1, y + 1},
	}
}

func next(x1, y1, x2, y2 int) [2]int {
	dx := x2 - x1
	dy := y2 - y1

	return [2]int{x2 + dx, y2 + dy}
}

func parseInput(sc *bufio.Scanner) (*[][]byte, error) {

	crossword := make([][]byte, 0)

	for sc.Scan() {
		line_string := strings.TrimSpace(sc.Text())
		if line_string == "" {
			continue
		}

		line := []byte(line_string)
		crossword = append(crossword, line)
	}

	return &crossword, nil
}

func validateInput(field *[][]byte) error {
	if field == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	} else if len(*field) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	rowLength := len((*field)[0])

	for y := 0; y < len(*field); y++ {

		if rowLength != len((*field)[y]) {
			return fmt.Errorf("%s rows have unequal lengths: %w", day, solver.ErrInvalidInput)
		}

		for x := 0; x < rowLength; x++ {
			switch (*field)[y][x] {
			case 'X', 'M', 'A', 'S':
				continue
			default:
				return fmt.Errorf("%s unknown character %s in input: %w", day, string((*field)[x][y]), solver.ErrInvalidInput)
			}
		}
	}

	return nil
}
