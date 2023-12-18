package d16

import (
	"bufio"
	"log"
	"os"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var ioMap = map[int32][4][]Direction{
	'.':  {{Up}, {Down}, {Left}, {Right}},
	'/':  {{Right}, {Left}, {Down}, {Up}},
	'\\': {{Left}, {Right}, {Up}, {Down}},
	'-':  {{Left, Right}, {Left, Right}, {Left}, {Right}},
	'|':  {{Up}, {Down}, {Up, Down}, {Up, Down}},
}

type tile struct {
	IO  [4][]Direction
	Out [4]bool
}

type Board struct {
	tiles      map[Point]*tile
	maxX, maxY int
}

type Point struct {
	X, Y int
}

func (p Point) next(d Direction) Point {
	switch d {
	case Up:
		return Point{p.X, p.Y - 1}
	case Down:
		return Point{p.X, p.Y + 1}
	case Left:
		return Point{p.X - 1, p.Y}
	case Right:
		return Point{p.X + 1, p.Y}
	}
	panic("unknown Direction")
}

func ParseBoard(lines []string) *Board {
	b := &Board{
		tiles: make(map[Point]*tile),
		maxX:  len(lines[0]) - 1,
		maxY:  len(lines) - 1,
	}
	for y, line := range lines {
		for x, c := range line {
			b.tiles[Point{x, y}] = &tile{
				IO: ioMap[c],
			}
		}
	}
	return b
}

type Beam struct {
	D Direction
	P Point
}

func (b *Board) Energize(start Beam) {
	open := make(map[Beam]bool)
	closed := make(map[Beam]bool)
	open[start] = true
	pop := func() Beam {
		for bm := range open {
			delete(open, bm)
			return bm
		}
		panic("empty")
	}
	for len(open) > 0 {
		bm := pop()
		t := b.tiles[bm.P]
		for _, d := range t.IO[bm.D] {
			t.Out[d] = true
			newP := bm.P.next(d)
			if !b.inBounds(newP) {
				continue
			}
			newBm := Beam{d, newP}
			if _, ok := closed[newBm]; !ok {
				open[newBm] = true
			}
		}
		closed[bm] = true
	}
}

func (b *Board) Deenergize() {
	for y := 0; y <= b.maxY; y++ {
		for x := 0; x <= b.maxX; x++ {
			t := b.tiles[Point{x, y}]
			for i := range t.Out {
				t.Out[i] = false
			}
		}
	}
}

func (b *Board) CountEnergized() int {
	count := 0
	for y := 0; y <= b.maxY; y++ {
		for x := 0; x <= b.maxX; x++ {
			for _, o := range b.tiles[Point{x, y}].Out {
				if o {
					count++
					break
				}
			}
		}
	}
	return count
}

func (b *Board) inBounds(p Point) bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}
	if p.X > b.maxX || p.Y > b.maxY {
		return false
	}
	return true
}

func (b *Board) BestEnergized() int {
	b.Deenergize()
	best := 0
	for x := 0; x <= b.maxX; x++ {
		// top
		b.Energize(Beam{Down, Point{x, 0}})
		e := b.CountEnergized()
		if e > best {
			best = e
		}
		b.Deenergize()
		// bottom
		b.Energize(Beam{Up, Point{x, b.maxY}})
		e = b.CountEnergized()
		if e > best {
			best = e
		}
		b.Deenergize()
	}
	for y := 0; y <= b.maxY; y++ {
		// leftest
		b.Energize(Beam{Right, Point{0, y}})
		e := b.CountEnergized()
		if e > best {
			best = e
		}
		b.Deenergize()
		// rightest
		b.Energize(Beam{Left, Point{b.maxX, y}})
		e = b.CountEnergized()
		if e > best {
			best = e
		}
		b.Deenergize()
	}
	return best
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
