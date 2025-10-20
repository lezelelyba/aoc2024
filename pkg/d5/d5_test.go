package d5

import (
	"advent2024/pkg/solver"
	"errors"
	"strings"
	"testing"
)

var (
	inputTest = `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`
)

func TestPart1(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{name: "test input", input: inputTest, want: "143"},
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
		{name: "test input", input: inputTest, want: "123"},
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
		{"missing rules", `123, 123`},
		{"missing records", `123|123`},
		{"non numeric rules", "123|124\n124|adsf\n\n123, 123"},
		{"non numeric records", "123|124\n124|125\n\n123, asdf"},
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
