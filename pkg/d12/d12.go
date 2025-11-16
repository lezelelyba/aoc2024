// Skeleton Package for new days
package d12

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"

	"container/list"

	"advent2024/pkg/solver"
)

// Solver name
var day = "d12"

// PuzzleStruct
type PuzzleStruct struct {
	field *[][]Point
}

type BorderName int

const (
	Top    BorderName = 0
	Right             = 1
	Bottom            = 2
	Left              = 3
)

type Coord struct {
	x, y int
}

type Point struct {
	Byte   byte
	Region int
}

type Region struct {
	Area      int
	Perimeter int
	Edges     int
}

type Counter struct {
	cnt int
}

func NewCounter() Counter {
	return Counter{cnt: 0}
}

func (c *Counter) Next() int {
	ret := c.cnt
	c.cnt++
	return ret
}

func (c *Counter) Value() int {
	return c.cnt
}

// Registers day wih the registry
func init() {
	solver.Register(day, func() solver.PuzzleSolver {
		return NewSolver()
	})
}

// Constructor
func NewSolver() *PuzzleStruct {
	return &PuzzleStruct{}
}

// Initializes the PuzzleStruct with input
// Return nil on success
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

// Solves the puzzle
// Accepts part as parameter
// Returns string containing the solution of the puzzle
func (p *PuzzleStruct) Solve(part int) (string, error) {
	// preprocess - determine regions

	regionCnt := NewCounter()

	for y := 0; y < len(*p.field); y++ {
		for x := 0; x < len((*p.field)[0]); x++ {
			if visited(p.field, Coord{x: x, y: y}) {
				continue
			}

			regionCnt.Next()

			(*p.field)[y][x].Region = regionCnt.Value()

			stack := list.New()
			stack.PushFront(Coord{x: x, y: y})

			for stack.Len() > 0 {
				cur := stack.Front().Value.(Coord)
				stack.Remove(stack.Front())

				neighbors := filterNeighbors(p.field, Coord{x: cur.x, y: cur.y}, sameByte)

				for _, n := range neighbors {
					if !visited(p.field, n) {
						(*p.field)[n.y][n.x].Region = regionCnt.Value()
						stack.PushFront(n)
					}
				}
			}
		}
	}

	switch part {
	case 1:
		sum := 0

		var regions = make(map[int]Region)

		for y := 0; y < len(*p.field); y++ {
			for x := 0; x < len((*p.field)[0]); x++ {
				id := (*p.field)[y][x].Region

				if _, ok := regions[id]; !ok {
					regions[id] = Region{Area: 0, Perimeter: 0, Edges: 0}
				}

				region := regions[id]
				region.Area += 1
				region.Perimeter += borderCount(p.field, Coord{x: x, y: y})

				regions[id] = region
			}
		}

		for _, region := range regions {
			sum += region.Area * region.Perimeter
		}

		return strconv.Itoa(sum), nil
	case 2:
		sum := 0

		var regions = make(map[int]Region)

		for y := 0; y < len(*p.field); y++ {
			for x := 0; x < len((*p.field)[0]); x++ {
				id := (*p.field)[y][x].Region

				if _, ok := regions[id]; !ok {
					regions[id] = Region{Area: 0, Perimeter: 0, Edges: 0}
				}

				region := regions[id]
				region.Area += 1
				region.Edges += newEdges(p.field, Coord{x: x, y: y})

				regions[id] = region
			}
		}

		for _, region := range regions {
			sum += region.Area * region.Edges
		}

		return strconv.Itoa(sum), nil
	}

	return "", fmt.Errorf("%s unknown part %d: %w", day, part, solver.ErrUnknownPart)
}

// Parses provided input
// Returns parsed string
func parseInput(sc *bufio.Scanner) (*[][]Point, error) {

	var field = make([][]Point, 0)

	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		row := make([]Point, 0, len(line))

		for _, b := range line {
			coord := Point{Byte: b, Region: -1}
			row = append(row, coord)
		}
		field = append(field, row)
	}

	if sc.Err() != nil {
		return nil, fmt.Errorf("%s scan error %s: %w", day, sc.Err(), solver.ErrInvalidInput)
	}

	return &field, nil
}

// Validates parsed input
// Returns nil in case of successfull validation
func validateInput(field *[][]Point) error {
	if field == nil {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	if len(*field) == 0 {
		return fmt.Errorf("%s empty records: %w", day, solver.ErrInvalidInput)
	}

	checkedLength := len((*field)[0])

	for i := 0; i < len(*field); i++ {
		if len((*field)[i]) != checkedLength {
			return fmt.Errorf("%s unequal line length: %w", day, solver.ErrInvalidInput)
		}
	}

	return nil
}

// Checks if field was already visited when determining region
// visited == Region is not -1
func visited(field *[][]Point, c Coord) bool {
	return (*field)[c.y][c.x].Region != -1
}

// Filters neighbor based on provided function
// Returns list of neighbors who pass the preficate f
func filterNeighbors(field *[][]Point, c Coord, f func(Point, Point) bool) []Coord {

	var ret = make([]Coord, 0, 4)

	neighbors := [4]Coord{
		{c.x, c.y - 1},
		{c.x + 1, c.y},
		{c.x, c.y + 1},
		{c.x - 1, c.y},
	}

	my := len(*field)
	mx := len((*field)[0])

	for _, n := range neighbors {
		// neighbor is out of range
		if n.x < 0 || n.x >= mx || n.y < 0 || n.y >= my {
			continue
		}

		this := (*field)[c.y][c.x]
		other := (*field)[n.y][n.x]

		if f(this, other) {
			ret = append(ret, n)
		}
	}

	return ret
}

// true if both Points have save Byte value
func sameByte(this, other Point) bool {
	return this.Byte == other.Byte
}

// true if both Points are in same region
func sameRegion(this, other Point) bool {
	return this.Region == other.Region
}

// returns count of borders for the point
// == amount of neighbors from different region
func borderCount(field *[][]Point, c Coord) int {
	return 4 - len(filterNeighbors(field, c, sameRegion))
}

// returns map of borders for the point
// {Top, Right, Bottom, Left}
func edgeMap(field *[][]Point, c Coord) [4]bool {
	var ret [4]bool

	neighbors := []Coord{
		{c.x, c.y - 1},
		{c.x + 1, c.y},
		{c.x, c.y + 1},
		{c.x - 1, c.y},
	}

	my := len(*field)
	mx := len((*field)[0])

	for i, n := range neighbors {
		if n.x < 0 || n.x >= mx || n.y < 0 || n.y >= my {
			ret[i] = true
			continue
		}

		this := (*field)[c.y][c.x]
		other := (*field)[n.y][n.x]

		if !sameRegion(this, other) {
			ret[i] = true
		}
	}

	return ret
}

// returns number of new edges introducted by this point
// new edge == edge which doesn't continue from previous neighbor
// previous neighbor == neighbor to the left or top, as we scan from left > right and top > bottom
func newEdges(field *[][]Point, c Coord) int {
	borders := borderCount(field, c)

	cMap := edgeMap(field, c)

	neighborTop := Coord{x: c.x, y: c.y - 1}
	neighborLeft := Coord{x: c.x - 1, y: c.y}

	if (cMap[Top] || cMap[Bottom]) && !cMap[Left] {
		nLMap := edgeMap(field, neighborLeft)

		if cMap[Top] && cMap[Top] == nLMap[Top] {
			borders--
		}

		if cMap[Bottom] && cMap[Bottom] == nLMap[Bottom] {
			borders--
		}

	}

	if (cMap[Left] || cMap[Right]) && !cMap[Top] {
		nTMap := edgeMap(field, neighborTop)

		if cMap[Left] && cMap[Left] == nTMap[Left] {
			borders--
		}

		if cMap[Right] && cMap[Right] == nTMap[Right] {
			borders--
		}
	}

	return borders
}
