package d8

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"advent2024/pkg/solver"
)

var day = "d8"

func init() {
	solver.Register("d8", func() solver.PuzzleSolver {
		return NewSolver()
	})
}

type PuzzleStruct struct {
	field    [][]byte
	antennas map[byte][]Coords
}

type Coords struct {
	x, y int
}
type LineIter struct {
	curr, delta Coords
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

	if err := validateInput(field); err != nil {
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
			// time.Sleep(time.Second / 1000)
			antennas := p.antennas[freq]
			for a1 := 0; a1 < len(antennas); a1++ {
				for a2 := a1 + 1; a2 < len(antennas); a2++ {

					// line in one direction
					iter := NewLineIter(antennas[a1], antennas[a2])

					for an, ok := iter.Next(); ok; an, ok = iter.Next() {
						if !p.inField(an) {
							break
						}
						antinodesMap[an] = struct{}{}
					}

					// line in the other direction
					iter = NewLineIter(antennas[a2], antennas[a1])

					for an, ok := iter.Next(); ok; an, ok = iter.Next() {
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

	case 21:

		// Rerwite using Go Routines

		sum := 0
		antinodesMap := make(map[Coords]struct{})

		tasks := make(chan byte, 10)
		results := make(chan Coords, 100)

		var wg sync.WaitGroup

		numWorkers := runtime.NumCPU()

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for t := range tasks {
					// time.Sleep(time.Second / 1000)
					freq := t
					antennas := p.antennas[freq]
					for a1 := 0; a1 < len(antennas); a1++ {
						for a2 := a1 + 1; a2 < len(antennas); a2++ {

							// line in one direction
							iter := NewLineIter(antennas[a1], antennas[a2])

							for an, ok := iter.Next(); ok; an, ok = iter.Next() {
								if !p.inField(an) {
									break
								}
								results <- an
							}

							// line in the other direction
							iter = NewLineIter(antennas[a2], antennas[a1])

							for an, ok := iter.Next(); ok; an, ok = iter.Next() {
								if !p.inField(an) {
									break
								}
								results <- an
							}
						}
					}

				}
			}()
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		go func() {
			for t := range p.antennas {
				tasks <- t
			}
			close(tasks)
		}()

		for an := range results {
			antinodesMap[an] = struct{}{}
		}

		sum = len(antinodesMap)

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
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

func validateInput(field *[][]byte) error {
	if field == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	} else if len(*field) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	rowLength := len((*field)[0])

	for y := 0; y < len(*field); y++ {
		if rowLength != len((*field)[y]) {
			return fmt.Errorf("%s rows have unequal lengths: %w", day, solver.ErrInvalidInput)
		}
	}

	return nil
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

func NewLineIter(this, other Coords) LineIter {
	return LineIter{curr: this, delta: Sub(other, this)}
}

func (i *LineIter) Next() (Coords, bool) {
	i.curr = Add(i.curr, i.delta)
	return i.curr, true
}
