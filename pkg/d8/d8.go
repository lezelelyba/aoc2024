package d8

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

type PuzzleStruct struct {
	field    [][]byte
	antennas map[byte][]Coords
}

type Coords struct {
	x, y int
}

func init() {
	solver.Register("d8", func() solver.PuzzleSolver {
		return NewSolver()
	})
}

func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

func (p *PuzzleStruct) Init(reader io.Reader) error {
	field, err := parseInput(bufio.NewScanner(reader))

	if err != nil {
		log.Print(err)
		return err
	}

	p.field = *field
	p.findAntennas()

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		antinodesMap := make(map[Coords]struct{})

		for freq := range p.antennas {
			antennas := p.antennas[freq]
			for a1 := 0; a1 < len(antennas); a1++ {
				for a2 := a1 + 1; a2 < len(antennas); a2++ {
					antinodes := FindAntinodes(antennas[a1], antennas[a2])

					for _, an := range antinodes {
						if p.inField(an) {
							antinodesMap[an] = struct{}{}
						}
					}
				}
			}
		}

		sum = len(antinodesMap)

		return strconv.Itoa(sum), nil
	case 2:

		sum := 0

		antinodesMap := make(map[Coords]struct{})

		for freq := range p.antennas {
			antennas := p.antennas[freq]
			for a1 := 0; a1 < len(antennas); a1++ {
				for a2 := a1 + 1; a2 < len(antennas); a2++ {

					// line in one direction
					iter := NewLineIter(antennas[a1], antennas[a2])

					for an := iter.Next(); ; an = iter.Next() {
						if !p.inField(an) {
							break
						}
						antinodesMap[an] = struct{}{}
					}

					// line in the other direction
					iter = NewLineIter(antennas[a2], antennas[a1])

					for an := iter.Next(); ; an = iter.Next() {
						if !p.inField(an) {
							break
						}
						antinodesMap[an] = struct{}{}
					}
				}
			}
		}

		sum = len(antinodesMap)

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("unknown Part %d", part)
}

func parseInput(sc *bufio.Scanner) (*[][]byte, error) {

	result := make([][]byte, 0)

	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		bs := []byte(s)

		result = append(result, bs)
	}

	return &result, nil
}

func (p *PuzzleStruct) findAntennas() {

	p.antennas = make(map[byte][]Coords)

	for y := range len(p.field) {
		for x := range p.field[y] {
			c := p.field[y][x]
			if c != '.' {
				v, ok := p.antennas[c]
				if !ok {
					p.antennas[c] = make([]Coords, 0)
				}

				p.antennas[c] = append(v, Coords{x: x, y: y})
			}
		}
	}
}

func (p *PuzzleStruct) inField(c Coords) bool {

	if c.y >= 0 && c.y < len(p.field) {
		if c.x >= 0 && c.x < len(p.field[c.y]) {
			return true
		}
	}

	return false
}

func FindAntinodes(this, other Coords) [2]Coords {

	delta := Sub(this, other)

	a1 := Add(this, delta)
	a2 := Add(other, Invert(delta))

	return [2]Coords{a1, a2}
}

func Add(this, other Coords) Coords {
	return Coords{x: this.x + other.x, y: this.y + other.y}
}

func Sub(this, other Coords) Coords {
	return Coords{x: this.x - other.x, y: this.y - other.y}
}

func Eq(this, other Coords) bool {
	return this.x == other.x && this.y == other.y
}

func Invert(this Coords) Coords {
	return Coords{x: -this.x, y: -this.y}
}

type LineIter struct {
	curr, delta Coords
}

func NewLineIter(this, other Coords) LineIter {
	return LineIter{curr: this, delta: Sub(other, this)}
}

func (i *LineIter) Next() Coords {
	i.curr = Add(i.curr, i.delta)
	return i.curr
}
