package d5

import (
	"advent2024/pkg/solver"
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
)

var day = "d5"

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type Rules struct {
	before, after map[int]bool
}

func newRules() Rules {
	return Rules{before: make(map[int]bool), after: make(map[int]bool)}
}

type PuzzleStruct struct {
	rules   map[int]Rules
	updates [][]int
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	input, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		return err
	}

	p.rules = processRules(input[0])
	p.updates = input[1]

	if err := validateInput(&p.rules, &p.updates); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for _, update := range p.updates {
			if slices.IsSortedFunc(update, p.sortFunc()) {
				sum += update[len(update)/2]
			}
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0
		for _, update := range p.updates {
			if !slices.IsSortedFunc(update, p.sortFunc()) {
				slices.SortFunc(update, p.sortFunc())
				sum += update[len(update)/2]
			}
		}
		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func (p *PuzzleStruct) sortFunc() func(a, b int) int {
	return func(a, b int) int {
		_, bAfterA := p.rules[a].after[b]
		_, bBeforeA := p.rules[a].before[b]

		switch {
		case bAfterA:
			return -1
		case bBeforeA:
			return 1
		default:
			return 0
		}
	}
}

type section int

const (
	RULES   section = 0
	RECORDS section = 1
)

func parseInput(sc *bufio.Scanner) (*[2][][]int, error) {

	section := RULES

	rules := make([][]int, 0)
	updates := make([][]int, 0)

	for sc.Scan() {
		line_string := strings.TrimSpace(sc.Text())

		if line_string == "" {
			section = RECORDS
			continue
		}

		switch section {
		case RULES:
			rule_string := strings.Split(sc.Text(), "|")
			if len(rule_string) != 2 {
				maxOutput := min(len(sc.Text()), 80)
				return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
			}

			rule := make([]int, 0)

			for _, si := range rule_string {
				v, err := strconv.Atoi(si)
				if err != nil {
					maxOutput := min(len(sc.Text()), 80)
					return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
				}
				rule = append(rule, v)
			}

			rules = append(rules, rule)

		case RECORDS:
			update_string := strings.Split(sc.Text(), ",")

			update := make([]int, 0)

			for _, si := range update_string {
				v, err := strconv.Atoi(si)
				if err != nil {
					maxOutput := min(len(sc.Text()), 80)
					return nil, fmt.Errorf("%s unable to parse \"%v\": %w", day, sc.Text()[:maxOutput], solver.ErrInvalidInput)
				}
				update = append(update, v)
			}

			updates = append(updates, update)
		}
	}

	var result [2][][]int
	result[0] = rules
	result[1] = updates
	return &result, nil
}

func processRules(parsedRules [][]int) map[int]Rules {

	result := make(map[int]Rules)

	for _, parsedRule := range parsedRules {
		before, after := parsedRule[0], parsedRule[1]

		rules, ok := result[before]
		if !ok {
			rules = newRules()
		}

		rules.updateRules([]int{}, []int{after})
		result[before] = rules

		rules, ok = result[after]
		if !ok {
			rules = newRules()
		}

		rules.updateRules([]int{before}, []int{})
		result[after] = rules
	}

	return result
}

func validateInput(rules *map[int]Rules, updates *[][]int) error {
	if rules == nil {
		return fmt.Errorf("%s empty rules: %w", day, solver.ErrInvalidInput)
	} else if len(*rules) == 0 {
		return fmt.Errorf("%s empty rules: %w", day, solver.ErrInvalidInput)
	}

	if updates == nil {
		return fmt.Errorf("%s empty updates: %w", day, solver.ErrInvalidInput)
	} else if len(*updates) == 0 {
		return fmt.Errorf("%s empty updates: %w", day, solver.ErrInvalidInput)
	}

	return nil
}

func (r *Rules) updateRules(before, after []int) {

	for _, b := range before {
		r.before[b] = true
	}
	for _, a := range after {
		r.after[a] = true
	}
}
