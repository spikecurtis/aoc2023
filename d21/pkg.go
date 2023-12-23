package d21

import (
	"bufio"
	"log"
	"os"
)

type Tile int32

const (
	Plot Tile = '.'
	Rock Tile = '#'
)

type Point struct {
	X, Y int
}

type Board struct {
	Tiles      [][]Tile
	Start      Point
	MaxX, MaxY int
}

func ParseBoard(lines []string) Board {
	maxX := len(lines[0]) - 1
	maxY := len(lines) - 1
	b := Board{
		Tiles: make([][]Tile, len(lines)),
		MaxX:  maxX,
		MaxY:  maxY,
	}
	for y, line := range lines {
		b.Tiles[y] = make([]Tile, len(line))
		for x, c := range line {
			if c == 'S' {
				b.Start = Point{x, y}
				c = int32(Plot)
			}
			b.Tiles[y][x] = Tile(c)
		}
	}
	return b
}

func (b Board) PlotsForSteps(n int) map[Point]bool {
	this := map[Point]bool{b.Start: true}
	for i := 0; i < n; i++ {
		next := make(map[Point]bool)
		for p := range this {
			for _, n := range p.Step() {
				if b.IsPlot(n) {
					next[n] = true
				}
			}
		}
		this = next
	}
	return this
}

func (b Board) IsPlot(p Point) bool {
	y := p.Y % (b.MaxY + 1)
	if y < 0 {
		y = y + b.MaxY + 1
	}
	x := p.X % (b.MaxX + 1)
	if x < 0 {
		x = x + b.MaxX + 1
	}
	if b.Tiles[y][x] == Rock {
		return false
	}
	return true
}

func (p Point) Step() []Point {
	out := make([]Point, 4)
	out[0] = Point{p.X, p.Y + 1}
	out[1] = Point{p.X, p.Y - 1}
	out[2] = Point{p.X + 1, p.Y}
	out[3] = Point{p.X - 1, p.Y}
	return out
}

func (b Board) Part2(n int) int {
	// The board has special properties where up, down, left, right from start
	// are clear.  This makes the result quadratic for steps = pk + m where p
	// is the full board width and m is the half width.  So, we brute force
	// solve for k = 0, 1, 2, then solve the quadratic equation, and apply it
	// for n = pk + m (that is, k = (n - m)/p)
	p := b.MaxX + 1
	m := b.MaxX / 2
	c := len(b.PlotsForSteps(m))
	f1 := len(b.PlotsForSteps(m + p))
	f2 := len(b.PlotsForSteps(m + p + p))
	a := (f2 - 2*f1 + c) / 2
	bb := f1 - c - a
	k := (n - m) / p
	return a*k*k + bb*k + c
}

func GetInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
