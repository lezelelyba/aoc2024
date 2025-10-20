package d3

import (
	"advent2024/pkg/solver"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

var day = "d3"

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type puzzleEntry struct {
	instruction string
	i, j        int
}

type PuzzleStruct struct {
	entries *[]puzzleEntry
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{entries: &[]puzzleEntry{}}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {

	s, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	p.entries = parseInput(string(s))

	if err := validateInput(p.entries); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {

	switch part {
	case 1:
		sum := 0
		for _, entry := range *p.entries {
			if entry.instruction == "mul" {
				sum += entry.i * entry.j
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0
		sum_enabled := true
		for _, entry := range *p.entries {
			if entry.instruction == "mul" && sum_enabled {
				sum += entry.i * entry.j
			} else if entry.instruction == "do" {
				sum_enabled = true
			} else if entry.instruction == "don't" {
				sum_enabled = false
			} else {
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func parseInput(s string) *[]puzzleEntry {

	entries := []puzzleEntry{}

	re := regexp.MustCompile(`(don't)\(\)|(do)\(\)|(mul)\((\d{1,3}),(\d{1,3})\)`)
	// re := regexp.MustCompile(`(don't)\(\)|(do)\(\)`)

	matches := re.FindAllStringSubmatch(s, -1)

	for _, match := range matches {

		if match[3] == "mul" {
			si, sj := match[4], match[5]
			i, _ := strconv.Atoi(si)
			j, _ := strconv.Atoi(sj)

			entries = append(entries, puzzleEntry{instruction: "mul", i: i, j: j})
		} else if match[2] == "do" {
			entries = append(entries, puzzleEntry{instruction: "do", i: 0, j: 0})
		} else if match[1] == "don't" {
			entries = append(entries, puzzleEntry{instruction: "don't", i: 0, j: 0})
		} else {
		}
	}

	return &entries
}

func validateInput(entries *[]puzzleEntry) error {
	if entries == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	} else if len(*entries) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	return nil
}
