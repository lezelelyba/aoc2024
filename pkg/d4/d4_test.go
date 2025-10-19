package d4

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	input := `MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`

	want := "18"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := `MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`

	want := "9"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestEmptyLine(t *testing.T) {
	input := `MMMSXXMASM

MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`

	want := "9"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestBadPart(t *testing.T) {
	input := `MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`

	_ = "9"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	_, err := puzzle.Solve(5)

	if err == nil {
		t.Errorf("Expected wrong part")
	}

}
