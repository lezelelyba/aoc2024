// Skeleton Package for new days
package d0

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"

	"advent2024/pkg/solver"
)

// Solver name
var day = "d0"

// PuzzleStruct
type PuzzleStruct struct {
	input string
}

// Registers day wih the registry
func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

// Constructor
func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

// Initializes the PuzzleStruct with input
// Return nil on success
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

// Solves the puzzle
// Accepts part as parameter
// Returns string containing the solution of the puzzle
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

// Parses provided input
// Returns parsed string
func parseInput(sc *bufio.Scanner) (string, error) {

	for sc.Scan() {
	}

	if sc.Err() != nil {
		return "", fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return "", nil
}

// Validates parsed input
// Returns nil in case of successfull validation
func validateInput(entry string) error {
	if len(entry) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	return nil
}
