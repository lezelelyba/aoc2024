package d1

import (
	"strings"
	"testing"
)

func TestDifference(t *testing.T) {

	tests := []struct {
		i, j, expected int
	}{
		{1, 1, 0},
		{1, 2, 1},
		{2, 1, 1},
		{2, 4, 2},
	}

	for _, test := range tests {
		if got := difference(test.i, test.j); got != test.expected {
			t.Errorf("Got %d expected %d. Test: %v", got, test.expected, test)
		}
	}
}

func TestPart1(t *testing.T) {
	input := `3   4
4   3
2   5
1   3
3   9
3   3`

	want := "11"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := `3   4
4   3
2   5
1   3
3   9
3   3`

	want := "31"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestUnknownPart(t *testing.T) {
	input := `
3   4
4   3
2   5
1   3
3   9
3   3`

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	_, err := puzzle.Solve(5)

	if err == nil {
		t.Errorf("No Error received")
	}
}

func TestParsing(t *testing.T) {
	tests := []string{
		`1 2 3`,
		`1 2
a 12`,
	}

	for _, test := range tests {
		puzzle := NewSolver()
		err := puzzle.Init(strings.NewReader(test))

		if err == nil {
			t.Errorf("No Error received")
		}

	}
}
