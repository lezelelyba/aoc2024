package d0

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	input := ``

	want := "0"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := ``

	want := "0"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}
