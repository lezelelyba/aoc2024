package d6

import (
	"advent2024/pkg/solver"
	"errors"
	"strings"
	"testing"
)

var (
	inputTest = `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "41"},
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
		{name: "test input", input: inputTest, want: "6"},
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
		{"missing guard", "..\n##\n"},
		{"unequal rows", "..\n##\n.v.\n"},
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
func TestValidInput(t *testing.T) {

	cases := []struct {
		name  string
		input string
	}{
		{"guard orientation 1", "..\n##\n^.\n"},
		{"guard orientation 2", "..\n##\n>.\n"},
		{"guard orientation 3", "..\n##\nv.\n"},
		{"guard orientation 4", "..\n##\n<.\n"},
	}

	var want error = nil

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
