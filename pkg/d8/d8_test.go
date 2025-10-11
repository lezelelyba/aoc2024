package d8

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	input := `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`

	want := "14"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(1)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`

	want := "34"

	puzzle := NewSolver()
	_ = puzzle.Init(strings.NewReader(input))
	got, _ := puzzle.Solve(2)

	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}

func TestAntinodes(t *testing.T) {
	tests := []struct {
		thisx, thisy   int
		otherx, othery int
		a1x, a1y       int
		a2x, a2y       int
	}{
		{1, 1, 3, 2, -1, 0, 5, 3},
	}

	for _, tt := range tests {
		this := Coords{x: tt.thisx, y: tt.thisy}
		other := Coords{x: tt.otherx, y: tt.othery}

		want1 := Coords{x: tt.a1x, y: tt.a1y}
		want2 := Coords{x: tt.a2x, y: tt.a2y}

		as := FindAntinodes(this, other)

		got1 := as[0]
		got2 := as[1]

		if !Eq(want1, got1) || !Eq(want2, got2) {
			t.Errorf("Got %v %v expected %v %v", got1, got2, want1, want2)
		}
	}
}
