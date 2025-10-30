package d10

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"

	"advent2024/pkg/solver"
)

var day = "d10"

type PuzzleStruct struct {
	input string
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

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func parseInput(sc *bufio.Scanner) (string, error) {

	for sc.Scan() {
		// for test return fed input
		return sc.Text(), nil
	}

	if sc.Err() != nil {
		return "", fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return "", nil
}

func validateInput(entry string) error {
	if len(entry) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	if entry == "inputInvalid" {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}
	return nil
}

// test ci2
