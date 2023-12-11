package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	d := parseDiagram(getInput())
	ll := d.loopLen()
	fmt.Println(ll / 2)
	p2 := d.part2()
	fmt.Println(p2)
}

type point struct {
	x, y int
}

var (
	north = point{0, -1}
	south = point{0, 1}
	east  = point{1, 0}
	west  = point{-1, 0}
)

func (p point) plus(k point) point {
	return point{p.x + k.x, p.y + k.y}
}

type tile struct {
	symbol string
	point
	connections []point
}

type diagram struct {
	tiles      map[point]tile
	start      point
	maxX, maxY int
	loop       map[point]bool
}

func (d *diagram) inBounds(p point) bool {
	switch {
	case p.x < 0:
		return false
	case p.y < 0:
		return false
	case p.x > d.maxX:
		return false
	case p.y > d.maxY:
		return false
	default:
		return true
	}
}

func (d *diagram) loopLen() int {
cardinals:
	for _, dir := range []point{north, south, east, west} {
		prev := d.start
		this := prev.plus(dir)
		loopTiles := map[point]bool{this: true}
		for this != d.start {
			next, err := d.step(prev, this)
			if err != nil {
				continue cardinals
			}
			loopTiles[next] = true
			prev = this
			this = next
		}
		d.loop = loopTiles
		return len(loopTiles)
	}
	log.Fatal("didn't work")
	return -1
}

func (d *diagram) part2() int {
	unvisited := []point{{0, 0}}
	//visited := make(map[point]bool)
	discovered := map[point]bool{point{0, 0}: true}
	for len(unvisited) > 0 {
		this := unvisited[len(unvisited)-1]
		unvisited = unvisited[:len(unvisited)-1]
		//visited[this] = true
		candidates := d.connectedOutside(this)
		for _, c := range candidates {
			if !discovered[c] {
				unvisited = append(unvisited, c)
			}
			discovered[c] = true
		}
	}
	outside := 0
	for x := 0; x <= d.maxX; x++ {
		for y := 0; y <= d.maxY; y++ {
			allVisited := true
			for _, p := range []point{{0, 0}, {0, 1}, {1, 0}, {1, 1}} {
				if _, ok := discovered[point{x, y}.plus(p)]; !ok {
					allVisited = false
					break
				}
			}
			if allVisited {
				outside++
			}
		}
	}
	all := (d.maxX + 1) * (d.maxY + 1)
	return all - len(d.loop) - outside
}

func (d *diagram) connectedOutside(p point) []point {
	var candidates []point
	t, ok := d.tiles[p]
	if ok && d.loop[t.point] {
		// east
		if !in(t.plus(north), t.connections) {
			candidates = append(candidates, p.plus(east))
		}
		// south
		if !in(t.plus(west), t.connections) {
			candidates = append(candidates, p.plus(south))
		}
	} else {
		candidates = append(candidates, p.plus(east), p.plus(south))
	}
	t, ok = d.tiles[p.plus(point{-1, -1})]
	if ok && d.loop[t.point] {
		// east
		if !in(t.plus(east), t.connections) {
			candidates = append(candidates, p.plus(north))
		}
		// south
		if !in(t.plus(south), t.connections) {
			candidates = append(candidates, p.plus(west))
		}
	} else {
		candidates = append(candidates, p.plus(north), p.plus(west))
	}
	var out []point
	for _, c := range candidates {
		if c.x < 0 {
			continue
		}
		if c.y < 0 {
			continue
		}
		if c.x > d.maxX+1 {
			continue
		}
		if c.y > d.maxY+1 {
			continue
		}
		out = append(out, c)
	}
	return out
}

func in[P comparable](p P, s []P) bool {
	for _, ss := range s {
		if p == ss {
			return true
		}
	}
	return false
}

var errNotConnected = errors.New("not connected")
var errBadPath = errors.New("bad path")

func (d *diagram) step(prev, this point) (next point, err error) {
	if !d.connected(prev, this) {
		return point{}, errNotConnected
	}
	t := d.tiles[this]
	for _, n := range t.connections {
		if n == prev {
			continue
		}
		return n, nil
	}
	return point{}, errBadPath
}

func (d *diagram) connected(i, j point) bool {
	has := false
	for _, k := range d.tiles[i].connections {
		if k == j {
			has = true
		}
	}
	if !has {
		return false
	}
	for _, k := range d.tiles[j].connections {
		if k == i {
			return true
		}
	}
	return false
}

func parseDiagram(lines []string) *diagram {
	d := diagram{
		tiles: make(map[point]tile),
		maxY:  len(lines) - 1,
		maxX:  len(lines[0]) - 1,
	}
	for y, line := range lines {
		for x, c := range line {
			p := point{x, y}
			t := tile{point: p, symbol: string(c)}
			switch c {
			case 'S':
				d.start = p
				t.connections = connections(p, d, north, south, east, west)
			case '|':
				t.connections = connections(p, d, north, south)
			case '-':
				t.connections = connections(p, d, east, west)
			case 'L':
				t.connections = connections(p, d, north, east)
			case 'J':
				t.connections = connections(p, d, north, west)
			case '7':
				t.connections = connections(p, d, south, west)
			case 'F':
				t.connections = connections(p, d, south, east)
			case '.':
				// ground
			default:
				log.Fatal("unknown")
			}
			d.tiles[t.point] = t
		}
	}
	return &d
}

func connections(p point, d diagram, dir ...point) []point {
	var out []point
	for _, k := range dir {
		if j := p.plus(k); d.inBounds(j) {
			out = append(out, j)
		}
	}
	return out
}

func getInput() []string {
	f, err := os.Open("./input")
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
