package d8

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"

	"advent2024/pkg/solver"
)

type PuzzleStruct struct {
	input string
}

func init() {
	solver.Register("d0", func() solver.PuzzleSolver {
		return NewSolver()
	})
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	_, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("unknown Part %d", part)
}

func parseInput(sc *bufio.Scanner) (string, error) {

	for sc.Scan() {
	}

	return "", nil
}
