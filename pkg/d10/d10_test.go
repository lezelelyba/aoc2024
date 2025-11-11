package d10

import (
	"advent2024/pkg/solver"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

var (
	inputTest1 = `0123
1234
8765
9876`
	inputTest2 = `...0...
...1...
...2...
6543456
7.....7
8.....8
9.....9`
	inputTest3 = `..90..9
...1.98
...2..7
6543456
765.987
876....
987....`
	inputTest4 = `10..9..
2...8..
3...7..
4567654
...8..3
...9..2
.....01`
	inputTest = `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`
	inputTest1Part2 = `.....0.
..4321.
..5..2.
..6543.
..7..4.
..8765.
..9....`
	inputTest2Part2 = `..90..9
...1.98
...2..7
6543456
765.987
876....
987....`
	inputTest3Part2 = `012345
123456
234567
345678
4.6789
56789.
`
)

func TestValid(t *testing.T) {
	cases := []struct {
		name, input string
		part        int
		want        string
	}{
		{"test input1 part 1", inputTest1, 1, "1"},
		{"test input2 part 1", inputTest2, 1, "2"},
		{"test input3 part 1", inputTest3, 1, "4"},
		{"test input4 part 1", inputTest4, 1, "3"},
		{"test input part 1", inputTest, 1, "36"},
		{"test input 1 part 2", inputTest1Part2, 2, "3"},
		{"test input 2 part 2", inputTest2Part2, 2, "13"},
		{"test input 3 part 2", inputTest3Part2, 2, "227"},
		{"test input part 2", inputTest, 2, "81"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolver()
			_ = puzzle.Init(strings.NewReader(c.input))
			got, _ := puzzle.Solve(c.part)

			if got != c.want {
				t.Errorf("part %d: got %s expected %s", c.part, got, c.want)
			}
		})
	}
}

func TestUnknownPart(t *testing.T) {
	invalidPart := 3

	want := solver.ErrUnknownPart

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(inputTest))
	_, got := puzzle.Solve(invalidPart)

	if !errors.Is(got, want) {
		t.Errorf("Got %v expected %v", got, want)
	}
}

func TestInvalidInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty input", ``},
		{"invalid input", "inputInvalid"},
	}

	want := solver.ErrInvalidInput

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolver()
			got := puzzle.Init(strings.NewReader(c.input))

			if !errors.Is(got, want) {
				t.Errorf("Got %v expected %v", got, want)
			}
		})
	}
}

func TestValidWithCtx(t *testing.T) {
	cases := []struct {
		name, input string
		part        int
		want        string
	}{
		{"test input1 part 1", inputTest1, 1, "1"},
		{"test input2 part 1", inputTest2, 1, "2"},
		{"test input3 part 1", inputTest3, 1, "4"},
		{"test input4 part 1", inputTest4, 1, "3"},
		{"test input part 1", inputTest, 1, "36"},
		{"test input 1 part 2", inputTest1Part2, 2, "3"},
		{"test input 2 part 2", inputTest2Part2, 2, "13"},
		{"test input 3 part 2", inputTest3Part2, 2, "227"},
		{"test input part 2", inputTest, 2, "81"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolverWithCtx()
			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))
			got, _ := puzzle.SolveCtx(context.Background(), c.part)

			if got != c.want {
				t.Errorf("part %d: got %s expected %s", c.part, got, c.want)
			}
		})
	}
}

func TestCtxTimeout(t *testing.T) {
	cases := []struct {
		name, input string
		part        int
	}{
		{"test input part 1", inputTest, 1},
		{"test input part 2", inputTest, 2},
	}

	want := solver.ErrTimeout

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolverWithCtx()
			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))

			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			defer cancel()
			time.Sleep(500 * time.Millisecond)

			_, got := puzzle.SolveCtx(ctx, c.part)

			if got != want {
				t.Errorf("Got %s expected %s", got, want)
			}
		})
	}
}
