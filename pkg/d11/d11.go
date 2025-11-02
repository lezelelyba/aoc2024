package d11

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"container/list"

	"advent2024/pkg/solver"
)

// Solver name
var day = "d11"

// PuzzleStruct
type PuzzleStruct struct {
	l *list.List
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
	inputList, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	if err := validateInput(inputList); err != nil {
		log.Print(err)
		return err
	}

	p.l = inputList

	return nil
}

// Solves the puzzle
// Accepts part as parameter
// Returns string containing the solution of the puzzle
func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := Blink(p.l, 25)
		return strconv.Itoa(sum), nil
	case 2:
		sum := Blink(p.l, 75)
		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

// Parses provided input
// Returns parsed list
func parseInput(sc *bufio.Scanner) (*list.List, error) {
	var line string
	var resultList = list.New()

	// expecting only 1 line

	for i := 0; sc.Scan(); i++ {

		if i > 0 {
			return nil, fmt.Errorf("%s multiline input: %w", day, solver.ErrInvalidInput)
		}

		line = strings.TrimSpace(sc.Text())
		fields := strings.Fields(line)

		for _, field := range fields {
			n, err := strconv.Atoi(string(field))
			if err != nil {
				log.Print(err)
				return nil, fmt.Errorf("%s unable to convert %s to int: %w", day, string(field), solver.ErrInvalidInput)
			}

			resultList.PushBack(n)
		}
	}

	if sc.Err() != nil {
		return nil, fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return resultList, nil

}

// Validates parsed input
// Returns nil in case of successfull validation
func validateInput(l *list.List) error {
	if l == nil {
		return fmt.Errorf("%s empty input: %w", day, solver.ErrInvalidInput)
	} else if l.Len() == 0 {
		return fmt.Errorf("%s empty input: %w", day, solver.ErrInvalidInput)
	}

	return nil
}

func Blink(l *list.List, cnt int) int {

	var sum int

	var stack = list.New()

	type stackElem struct {
		depth int
		value int
	}

	var memoizationMap = make(map[stackElem]int)

	var new, new1, new2, curr stackElem

	for e := l.Front(); e != nil; e = e.Next() {
		stack.PushBack(stackElem{depth: 0, value: e.Value.(int)})
	}

	for e := stack.Front(); e != nil; {

		curr = e.Value.(stackElem)

		// if value is memoized, drop the value
		// if its original list, count it
		if val, ok := memoizationMap[curr]; ok {
			stack.Remove(e)
			e = stack.Front()

			if curr.depth == 0 {
				sum += val
			}
			continue
		}

		// full depth in
		// memoize the value as 1
		if curr.depth >= cnt {
			memoizationMap[curr] = 1
			stack.Remove(e)
			e = stack.Front()
			continue
		}

		// not full depth in
		// calculate the next value
		if curr.value == 0 {
			new = stackElem{depth: curr.depth + 1, value: 1}

			// if calculated value is memoized, saved this value
			if val, ok := memoizationMap[new]; ok {
				memoizationMap[curr] = val
				// if not then calculate the value
			} else {
				stack.PushFront(new)
			}

		} else if noOfDigits(curr.value)%2 == 0 {
			n1, n2 := splitNumber(e.Value.(stackElem).value)
			new1 = stackElem{depth: curr.depth + 1, value: n1}
			new2 = stackElem{depth: curr.depth + 1, value: n2}

			val1, ok1 := memoizationMap[new1]
			val2, ok2 := memoizationMap[new2]

			if !ok2 {
				stack.PushFront(new2)
			}

			if !ok1 {
				stack.PushFront(new1)
			}

			if ok1 && ok2 {
				memoizationMap[curr] = val1 + val2
			}
		} else {
			new = stackElem{depth: curr.depth + 1, value: curr.value * 2024}
			if val, ok := memoizationMap[new]; ok {
				memoizationMap[curr] = val
			} else {
				stack.PushFront(new)
			}
		}

		// work from front
		e = stack.Front()
	}

	return sum
}

func noOfDigits(number int) int {

	var digitCount int

	if number == 0 {
		return 1
	}

	for number > 0 {
		number = number / 10
		digitCount++
	}

	return digitCount
}

func splitNumber(number int) (int, int) {

	digitCount := noOfDigits(number)

	var mod = 1

	for i := 0; i < digitCount/2; i++ {
		mod = mod * 10
	}

	return number / mod, number % mod

}
