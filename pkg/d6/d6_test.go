package d6

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`

	want := "41"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`

	want := "6"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}
func TestPartWrongPart(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`

	_ = "6"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	_, err := puzzle.Solve(5)

	if err == nil {
		t.Errorf("Expected wrong part")
	}
}

func TestPartWithoutGuard(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#........
........#.
#.........
......#...`

	_ = "6"

	puzzle := NewSolver()
	err := puzzle.Init(strings.NewReader(input))

	if err == nil {
		t.Errorf("Expected err")
	}
}

func TestPartGuardRotation(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..<.....
........#.
#.........
......#...`

	want := "26"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPartGuardRotation2(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..>.....
........#.
#.........
......#...`

	want := "6"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPartGuardRotation3(t *testing.T) {
	input := `....#.....
.........#
..........
..#.......
.......#..
..........
.#..v.....
........#.
#.........
......#...`

	want := "4"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}
