package d1

import (
	"advent2024/pkg/solver"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

var (
	inputTest = `3   4
4   3
2   5
1   3
3   9
3   3`
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "11"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolver()
			_ = puzzle.Init(strings.NewReader(c.input))
			got, _ := puzzle.Solve(1)

			if got != c.want {
				t.Errorf("Got %s expected %s", got, c.want)
			}
		})
	}
}
func TestPart2(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "31"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolver()
			_ = puzzle.Init(strings.NewReader(c.input))
			got, _ := puzzle.Solve(2)

			if got != c.want {
				t.Errorf("Got %s expected %s", got, c.want)
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
	invalidInput2 := `1 2
a 3`

	cases := []struct {
		name  string
		input string
	}{
		{"empty input", ``},
		{"invalid input 1", `1 2 3`},
		{"invalid input 2", invalidInput2},
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

func TestPart1WithCtx(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "11"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolverWithCtx()
			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))
			got, _ := puzzle.SolveCtx(context.Background(), 1)

			if got != c.want {
				t.Errorf("Got %s expected %s", got, c.want)
			}
		})
	}
}

func TestPart2WithCtx(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "31"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolverWithCtx()
			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))
			got, _ := puzzle.SolveCtx(context.Background(), 2)

			if got != c.want {
				t.Errorf("Got %s expected %s", got, c.want)
			}
		})
	}
}

func TestCtxCancel(t *testing.T) {
	cases := []struct {
		name, input string
	}{
		{name: "test input", input: inputTest},
	}

	want := solver.ErrTimeout

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			puzzle := NewSolverWithCtx()
			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))

			ctx, cancel := context.WithTimeout(context.Background(), 250 * time.Millisecond)
			defer cancel()
			time.Sleep(500 * time.Millisecond)

			_, got := puzzle.SolveCtx(ctx, 1)

			if got != want {
				t.Errorf("Got %s expected %s", got, want)
			}
		})
	}
}
