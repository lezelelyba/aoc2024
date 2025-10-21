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
	inputTest  = ``
	inputEntry = ``
)

func TestValid(t *testing.T) {
	cases := []struct {
		name, input string
		part        int
		want        string
	}{
		{"test input part 1", inputTest, 1, "0"},
		{"entry input part 1", inputEntry, 1, "0"},
		{"test input part 2", inputTest, 2, "0"},
		{"entry input part 2", inputEntry, 2, "0"},
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

// func TestValidWithCtx(t *testing.T) {
// 	cases := []struct {
// 		name, input string
// 		part        int
// 		want        string
// 	}{
// 		{"test input part 1", inputTest, 1, "0"},
// 		{"entry input part 1", inputEntry, 1, "0"},
// 		{"test input part 2", inputTest, 2, "0"},
// 		{"entry input part 2", inputEntry, 2, "0"},
// 	}
//
// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			puzzle := NewSolverWithCtx()
// 			_ = puzzle.InitCtx(context.Background(), strings.NewReader(c.input))
// 			got, _ := puzzle.SolveCtx(context.Background(), c.part)
//
// 			if got != c.want {
// 				t.Errorf("part %d: got %s expected %s", c.part, got, c.want)
// 			}
// 		})
// 	}
// }

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
