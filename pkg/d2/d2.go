package d2

import (
	"advent2024/pkg/solver"
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

var day = "d2"

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type PuzzleStruct struct {
	reports *[][]int
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {

	reports, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		return err
	}

	if err := validateInput(reports); err != nil {
		log.Print(err)
		return err
	}

	p.reports = reports

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:

		sum := 0

		for _, report := range *p.reports {
			if report_safe(&report) {
				sum += 1
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for _, report := range *p.reports {
			if report_safe(&report) {
				sum += 1
			} else {
				modified_report := make([]int, len(report))
				for i := range len(report) {
					copy(modified_report, report)
					modified_report := append(modified_report[:i], modified_report[i+1:]...)
					if report_safe(&modified_report) {
						sum += 1
						break
					}
				}
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func parseInput(sc *bufio.Scanner) (*[][]int, error) {

	reports := make([][]int, 0)

	for sc.Scan() {

		report := make([]int, 0)

		if sc.Text() == "" {
			continue
		}

		vs := strings.Fields(sc.Text())

		for _, v := range vs {

			i, err := strconv.Atoi(v)

			if err != nil {
				return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text(), solver.ErrInvalidInput)
			}

			report = append(report, i)
		}

		reports = append(reports, report)
	}

	return &reports, nil
}

func validateInput(reports *[][]int) error {

	if reports == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	} else if len(*reports) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}
	return nil
}

func report_safe(report *[]int) bool {

	direction := 0
	prev := -1

	max := 3

	for i, v := range *report {

		if i > 0 {
			if v == prev {
				return false
			}

			if direction == 0 {
				if prev > v {
					direction = -1
				} else if prev < v {
					direction = 1
				} else {
				}
			}

			if direction == -1 && (prev < v || prev+max*direction > v) {
				return false
			} else if direction == 1 && (prev > v || prev+max*direction < v) {
				return false
			}
		}

		prev = v
	}

	return true
}
