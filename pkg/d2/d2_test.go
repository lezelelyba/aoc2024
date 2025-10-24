package d2

import (
	"advent2024/pkg/solver"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

var (
	inputTest = `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`
)

func TestValid(t *testing.T) {
	cases := []struct {
		name, input string
		part        int
		want        string
	}{
		{"test input part 1", inputTest, 1, "2"},
		{"test input part 2", inputTest, 2, "4"},
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
	inputComplex := `invalidComplex`

	cases := []struct {
		name  string
		input string
	}{
		{"empty input", ``},
		{"invalid input", "inputInvalid"},
		{"complex input", inputComplex},
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
		{"test input part 1", inputTest, 1, "2"},
		{"test input part 2", inputTest, 2, "4"},
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

			ctx, cancel := context.WithTimeout(context.Background(), 250 * time.Millisecond)
			defer cancel()

			time.Sleep(500 * time.Millisecond)

			_, got := puzzle.SolveCtx(ctx, c.part)

			if got != want {
				t.Errorf("Got %s expected %s", got, want)
			}
		})
	}
}
