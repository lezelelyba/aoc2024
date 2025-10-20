package d4

import (
	"advent2024/pkg/solver"
	"errors"
	"strings"
	"testing"
)

var (
	inputTest = `MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "18"},
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
		{name: "test input", input: inputTest, want: "9"},
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
	cases := []struct {
		name  string
		input string
	}{
		{"empty input", ``},
		{"invalid input", `Invalid Input`},
		{"unequal row lengths", "XMAS\nXMASXX"},
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
