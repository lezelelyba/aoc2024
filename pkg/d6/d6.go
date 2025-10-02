package d6

import (
	"advent2024/pkg/solver"
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type Orientation int

const (
	UP    Orientation = 0
	RIGHT Orientation = 1
	DOWN  Orientation = 2
	LEFT  Orientation = 3
)

func (o Orientation) String() string {
	strs := []string{"^", ">", "v", "<"}
	return strs[o]
}

func (o Orientation) Byte() byte {
	bytes := []byte{'^', '>', 'v', '<'}
	return bytes[o]
}

func toOrientation(b byte) (Orientation, error) {
	switch b {
	case '^':
		return UP, nil
	case '>':
		return RIGHT, nil
	case 'v':
		return DOWN, nil
	case '<':
		return LEFT, nil
	}

	return UP, fmt.Errorf("unable to determine orientation %b", b)
}

type Coords struct {
	x, y int
}

type Guard struct {
	c       Coords
	o       Orientation
	visited map[Coords][4]bool
}

func NewGuard(x, y int, o Orientation) Guard {
	return Guard{Coords{x, y}, o, map[Coords][4]bool{}}
}

type NotInFieldError struct {
	Resource string
}

func (e NotInFieldError) Error() string {
	return fmt.Sprintf("Guard out of field")
}

type LoopingError struct {
	Resource string
}

func (e LoopingError) Error() string {
	return fmt.Sprintf("Guard is looping")
}

func (g *Guard) Move(field *[][]byte) error {
	x, y := g.c.x, g.c.y
	f := *field

	if y < len(f) && y > -1 {
		if x < len(f[y]) && x > -1 {
			next := Coords{x, y}

			switch g.o {
			case UP:
				next = Coords{x, y - 1}
			case RIGHT:
				next = Coords{x + 1, y}
			case DOWN:
				next = Coords{x, y + 1}
			case LEFT:
				next = Coords{x - 1, y}
			}

			next_v := byte('.')

			if next.x > -1 && next.x < len(f[y]) && next.y > -1 && next.y < len(f) {
				next_v = f[next.y][next.x]
			}

			switch next_v {
			case '.':
				visited := g.visited[g.c]

				if visited[g.o] {
					return LoopingError{}
				}

				visited[g.o] = true
				g.visited[g.c] = visited

				g.c.x = next.x
				g.c.y = next.y
			case '#':
				switch g.o {
				case UP:
					g.o = RIGHT
				case RIGHT:
					g.o = DOWN
				case DOWN:
					g.o = LEFT
				case LEFT:
					g.o = UP
				}
			}

			return nil
		}
	}

	return NotInFieldError{}
}

type PuzzleStruct struct {
	field [][]byte
	guard Guard
}

func init() {
	solver.Register("d6", func() solver.PuzzleSolver {
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

	gx, gy, err := findGuard(field)

	if err != nil {
		log.Print(err)
		return err
	}

	orientation, err := toOrientation(p.field[gy][gx])

	if err != nil {
		log.Print(err)
		return err
	}

	p.field[gy][gx] = '.'

	p.guard = NewGuard(gx, gy, orientation)

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for p.guard.Move(&p.field) == nil {
		}

		sum = len(p.guard.visited)

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		og := NewGuard(p.guard.c.x, p.guard.c.y, p.guard.o)

		for p.guard.Move(&p.field) == nil {
		}

		visited := p.guard.visited

		for coord := range visited {
			// get original guard
			p.guard = NewGuard(og.c.x, og.c.y, og.o)

			// skip initial field
			if p.guard.c == coord {
				continue
			}

			// put obstacle in place
			p.field[coord.y][coord.x] = '#'

			var err error
			// loop
			for err = p.guard.Move(&p.field); err == nil; err = p.guard.Move(&p.field) {
			}

			switch err.(type) {
			case LoopingError:
				sum += 1
			default:
			}

			// remove obstacle
			p.field[coord.y][coord.x] = '.'
		}

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

func findGuard(field *[][]byte) (int, int, error) {
	for y, line := range *field {
		for x, c := range line {
			switch c {
			case UP.Byte(), LEFT.Byte(), DOWN.Byte(), RIGHT.Byte():
				return x, y, nil
			default:
				continue
			}
		}
	}

	return -1, -1, fmt.Errorf("Guard not found")
}
