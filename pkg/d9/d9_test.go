package d9

import (
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	table := []struct {
		input, want string
	}{
		{input: "12345", want: "60"},
		{input: "20202", want: "23"},
	}

	for _, tt := range table {
		t.Run(tt.input, func(t *testing.T) {
			puzzle := NewSolver()
			_ = puzzle.Init(strings.NewReader(tt.input))
			got, _ := puzzle.Solve(1)

			if got != tt.want {
				t.Errorf("Got %s expected %s", got, tt.want)
			}
		})
	}
}

func TestPart1(t *testing.T) {
	input := `2333133121414131402`

	want := "1928"

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
