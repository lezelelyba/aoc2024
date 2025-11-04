package d11

import (
	"advent2024/pkg/solver"
	"container/list"
	"context"
	"fmt"
	"io"
	"strconv"
)

// PuzzleStruct with Context
type PuzzleStructWithCtx struct {
	PuzzleStruct
}

// Registers day wih the registry
func init() {
	solver.RegisterWithCtx(day, func() solver.PuzzleSolverWithCtx {
		return NewSolverWithCtx()
	})
}

// Constructor
func NewSolverWithCtx() *PuzzleStructWithCtx {
	return &PuzzleStructWithCtx{}
}

// Initializes the PuzzleStruct with input
func (p *PuzzleStructWithCtx) InitCtx(ctx context.Context, reader io.Reader) error {
	return p.PuzzleStruct.Init(reader)
}

// Solves the puzzle
// Accepts part as parameter
// Returns string containing the solution of the puzzle
func (p *PuzzleStructWithCtx) SolveCtx(ctx context.Context, part int) (string, error) {
	switch part {
	case 1:
		sum, err := BlinkWithCtx(p.l, 25, ctx)

		if err != nil {
			return "", err
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum, err := BlinkWithCtx(p.l, 75, ctx)

		if err != nil {
			return "", err
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func BlinkWithCtx(l *list.List, cnt int, ctx context.Context) (int, error) {
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

	for i, e := 0, stack.Front(); e != nil; i++ {

		if i%100000 == 0 {
			select {
			case <-ctx.Done():
				return -1, solver.ErrTimeout
			default:
			}
		}

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

	return sum, nil
}
