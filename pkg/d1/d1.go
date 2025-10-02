package d1

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

type PuzzleStruct struct {
	input *[2][]int
}

func init() {
	solver.Register("d1", func() solver.PuzzleSolver {
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

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		input_copy := p.inputCopy()

		slices.Sort(input_copy[0])
		slices.Sort(input_copy[1])

		diff := 0

		for i := 0; i < len(input_copy[0]); i++ {
			diff += difference(input_copy[0][i], input_copy[1][i])
		}

		return strconv.Itoa(diff), nil
	case 2:

		h := histogram(p.input[1])

		sim := 0

		for i := 0; i < len(p.input[0]); i++ {
			sim += similarity(p.input[0][i], h)
		}

		return strconv.Itoa(sim), nil
	}

	return "", fmt.Errorf("unknown Part %d", part)
}

func parseInput(sc *bufio.Scanner) (*[2][]int, error) {

	var left, right []int

	for sc.Scan() {

		if sc.Text() == "" {
			continue
		}

		vs := strings.Fields(sc.Text())

		if len(vs) != 2 {
			return nil, fmt.Errorf("unable to parse \"%v\"", sc.Text())
		}

		l, lerr := strconv.Atoi(vs[0])
		r, rerr := strconv.Atoi(vs[1])

		if lerr != nil || rerr != nil {
			return nil, fmt.Errorf("unable to parse \"%v\"", sc.Text())
		}

		left = append(left, l)
		right = append(right, r)
	}

	return &[2][]int{left, right}, nil
}

func histogram(l []int) *map[int]int {
	h := make(map[int]int)

	for _, v := range l {
		h[v] += 1
	}

	return &h
}

func difference(i int, j int) int {
	v := i - j
	if v < 0 {
		v = -v
	}

	return v
}

func similarity(v int, h *map[int]int) int {
	if c, ok := (*h)[v]; ok {
		return v * c
	}

	return 0
}

func (p *PuzzleStruct) inputCopy() *[2][]int {
	left := make([]int, len(p.input[0]))
	right := make([]int, len(left))

	copy(left, p.input[0])
	copy(right, p.input[1])

	return &[2][]int{left, right}
}
