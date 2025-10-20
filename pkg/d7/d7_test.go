package d7

import (
	"advent2024/pkg/solver"
	"errors"
	"strings"
	"testing"
)

var (
	inputTest = `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "3749"},
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
		{name: "test input", input: inputTest, want: "11387"},
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
		{"string result", `asdf: 123 123`},
		{"string parts", `1234: 123 asdf`},
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
func TestValidInputButUnsolvable(t *testing.T) {

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"no numbers", `123:`, "0"},
		{"only 1 number", `123: 123`, "0"},
		{"unsovable", `123: 123 10`, "0"},
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
