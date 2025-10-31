package d10

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"advent2024/pkg/solver"
)

var day = "d10"

type PuzzleStruct struct {
	field *[][]int
}

type Coord struct {
	x, y int
}

func init() {
	solver.Register(day, func() solver.PuzzleSolver {
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

	if err := validateInput(field); err != nil {
		log.Print(err)
		return err
	}

	p.field = field

	return nil
}

func (p *PuzzleStruct) Solve(part int) (string, error) {
	switch part {
	case 1:
		sum := 0

		for _, zero := range p.FindZeroes() {
			sum += p.ReachableSummits(zero)
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		for _, zero := range p.FindZeroes() {
			sum += p.PathsToSummits(zero)
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

func parseInput(sc *bufio.Scanner) (*[][]int, error) {

	result := make([][]int, 0)

	for sc.Scan() {
		// for test return fed input
		s := strings.TrimSpace(sc.Text())

		if s == "" {
			continue
		}

		var line = make([]int, 0, len(s))

		for _, c := range s {
			switch c {
			case '.':
				line = append(line, -1)
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				line = append(line, int(c-'0'))
			default:
				return nil, fmt.Errorf("%s invalid character %s in input: %w", day, string(c), solver.ErrInvalidInput)
			}
		}

		result = append(result, line)
	}

	if sc.Err() != nil {
		return nil, fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return &result, nil
}

func validateInput(field *[][]int) error {
	if field == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	if len(*field) == 0 {
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

func (p *PuzzleStruct) FindZeroes() []Coord {
	var result []Coord

	for y := 0; y < len(*p.field); y++ {
		for x := 0; x < len((*p.field)[y]); x++ {

			coord := Coord{x: x, y: y}
			value, err := p.ValueAt(coord)

			if err == nil && value == 0 {
				result = append(result, coord)
			}
		}
	}

	return result
}

func (p *PuzzleStruct) ReachableSummits(from Coord) int {
	return len(p.reachableSummits(from))
}

func (p *PuzzleStruct) PathsToSummits(from Coord) int {
	var sum int

	for _, v := range p.reachableSummits(from) {
		sum += v
	}

	return sum
}

func (p *PuzzleStruct) reachableSummits(coord Coord) map[Coord]int {
	var summits = make(map[Coord]int)

	value, err := p.ValueAt(coord)

	if err == nil && value == 9 {
		summits[coord] = 1
	} else {
		for _, next := range p.nextFrom(coord) {
			for summitCoord, newOccurences := range p.reachableSummits(next) {
				currentOccurences, ok := summits[summitCoord]

				if !ok {
					summits[summitCoord] = newOccurences
				} else {
					summits[summitCoord] = currentOccurences + newOccurences
				}
			}
		}
	}
	return summits
}

func (p *PuzzleStruct) ValueAt(coord Coord) (int, error) {
	if p == nil {
		return -1, fmt.Errorf("puzzle struct is nil")
	}

	if p.field == nil {
		return -1, fmt.Errorf("field is nil")
	}

	if coord.y < 0 || coord.y >= len(*p.field) {
		return -1, fmt.Errorf("coord %v outside of field", coord)
	}

	if coord.x < 0 || coord.x >= len((*p.field)[coord.y]) {
		return -1, fmt.Errorf("coord %v outside of field", coord)
	}

	return (*p.field)[coord.y][coord.x], nil
}

func (p *PuzzleStruct) nextFrom(coord Coord) []Coord {

	x, y := coord.x, coord.y
	coordValue := (*p.field)[y][x]

	var validNeighbors []Coord

	neighbors := []Coord{
		Coord{x: x, y: y - 1},
		Coord{x: x + 1, y: y},
		Coord{x: x, y: y + 1},
		Coord{x: x - 1, y: y},
	}

	for _, neighbor := range neighbors {
		neighborValue, err := p.ValueAt(neighbor)

		if err == nil && neighborValue-1 == coordValue {
			validNeighbors = append(validNeighbors, neighbor)
		}
	}

	return validNeighbors
}
