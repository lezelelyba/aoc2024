package d7

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

var day = "d7"

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type Equation struct {
	result  int
	numbers []int
}
type PuzzleStruct struct {
	equations *[]Equation
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	equations, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	if err := validateInput(equations); err != nil {
		log.Print(err)
		return err
	}

	p.equations = equations

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for _, e := range *p.equations {
			if solvable(e) {
				sum += e.result
			}
		}
		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for _, e := range *p.equations {
			if solvablePart2(e) {
				sum += e.result
			}
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func parseInput(sc *bufio.Scanner) (*[]Equation, error) {

	input := make([]Equation, 0)

	for sc.Scan() {

		if strings.TrimSpace(sc.Text()) == "" {
			continue
		}

		result_numbers := strings.Split(strings.TrimSpace(sc.Text()), ":")

		if len(result_numbers) != 2 {
			maxOutput := min(len(sc.Text()), 80)
			return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
		}

		result := strings.TrimSpace(result_numbers[0])
		numbers := strings.Fields(result_numbers[1])

		r, err := strconv.Atoi(result)

		if err != nil {
			maxOutput := min(len(sc.Text()), 80)
			return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
		}

		ints := make([]int, len(numbers))

		for i := range len(numbers) {
			n, err := strconv.Atoi(numbers[i])

			if err != nil {
				maxOutput := min(len(sc.Text()), 80)
				return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
			}

			ints[i] = n
		}

		input = append(input, Equation{result: r, numbers: ints})

	}

	return &input, nil
}

func validateInput(equations *[]Equation) error {
	if equations == nil {
		return fmt.Errorf("%s empty equations: %w", day, solver.ErrInvalidInput)
	} else if len(*equations) == 0 {
		return fmt.Errorf("%s empty equations: %w", day, solver.ErrInvalidInput)
	}

	return nil
}

func solvable(e Equation) bool {
	numbers := e.numbers

	if len(numbers) == 0 {
		return false
	}

	if len(numbers) == 1 {
		return false
	}
	return _solvable(e.result, numbers[0], numbers[1:])
}

func _solvable(res, acc int, nums []int) bool {

	if res < acc {
		return false
	}

	if len(nums) == 0 {
		return res == acc
	}

	if len(nums) == 1 {
		return _solvable(res, acc+nums[0], []int{}) || _solvable(res, acc*nums[0], []int{})
	}

	return _solvable(res, acc+nums[0], nums[1:]) || _solvable(res, acc*nums[0], nums[1:])
}

func solvablePart2(e Equation) bool {
	numbers := e.numbers

	if len(numbers) == 0 {
		return false
	}

	if len(numbers) == 1 {
		return false
	}
	return _solvablePart2(e.result, numbers[0], numbers[1:])
}

func _concat(i, j int) int {

	result := i
	c := j

	for c > 0 {
		c = c / 10
		result *= 10
	}

	return result + j
}

func _solvablePart2(res, acc int, nums []int) bool {

	if res < acc {
		return false
	}

	if len(nums) == 0 {
		return res == acc
	}

	if len(nums) == 1 {
		return _solvablePart2(res, acc+nums[0], []int{}) || _solvablePart2(res, acc*nums[0], []int{}) || _solvablePart2(res, _concat(acc, nums[0]), []int{})
	}

	return _solvablePart2(res, acc+nums[0], nums[1:]) || _solvablePart2(res, acc*nums[0], nums[1:]) || _solvablePart2(res, _concat(acc, nums[0]), nums[1:])
}
