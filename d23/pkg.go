package d23

import (
	"bufio"
	"errors"
	"log"
	"os"
)

type Point struct {
	X, Y int
}

type Board struct {
	Tiles      [][]int32
	Start      Point
	Goal       Point
	MaxX, MaxY int
}

type node struct {
	start, end Point
	len        int
	next       []Point
	goal       bool
}

func p1c(at int32, cur, next Point) bool {
	if at == '^' && next.Y < cur.Y {
		return true
	}
	if at == 'v' && next.Y > cur.Y {
		return true
	}
	if at == '<' && next.X < cur.X {
		return true
	}
	if at == '>' && next.X > cur.X {
		return true
	}
	return false
}

func p2c(at int32, _, _ Point) bool {
	if at == '^' {
		return true
	}
	if at == 'v' {
		return true
	}
	if at == '<' {
		return true
	}
	if at == '>' {
		return true
	}
	return false
}

var nodeCache map[Point]map[Point]node

func (b Board) getNode(prev, start Point, criterion func(at int32, cur, next Point) bool) (result node) {
	if m, ok := nodeCache[prev]; ok {
		if mm, okok := m[start]; okok {
			return mm
		}
	}
	defer func() {
		m, ok := nodeCache[prev]
		if !ok {
			m = make(map[Point]node)
			nodeCache[prev] = m
		}
		m[start] = result
	}()
	cur := start
	steps := 1
hike:
	for {
		r := node{start: start, end: cur, len: steps}
		if cur.Y == b.MaxY {
			r.goal = true
			return r
		}
		for _, n := range b.next(cur) {
			if n == prev {
				continue
			}
			if b.at(n) == '.' {
				steps++
				prev = cur
				cur = n
				continue hike
			}
			if criterion(b.at(n), cur, n) {
				r.next = append(r.next, n)
			}
		}
		return r
	}
}

func (b Board) next(p Point) []Point {
	out := make([]Point, 0, 4)
	if p.X > 1 {
		out = append(out, Point{p.X - 1, p.Y})
	}
	if p.X < b.MaxX {
		out = append(out, Point{p.X + 1, p.Y})
	}
	if p.Y > 1 {
		out = append(out, Point{p.X, p.Y - 1})
	}
	if p.Y < b.MaxY {
		out = append(out, Point{p.X, p.Y + 1})
	}
	return out
}

func (b Board) at(p Point) int32 {
	return b.Tiles[p.Y][p.X]
}

func (b Board) LongestPath1() int {
	nodeCache = make(map[Point]map[Point]node)
	visited := make(map[Point]bool)
	n := b.getNode(b.Start, Point{b.Start.X, b.Start.Y + 1}, p1c)
	l, err := b.longestPath(visited, n, p1c)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func (b Board) LongestPath2() int {
	nodeCache = make(map[Point]map[Point]node)
	visited := make(map[Point]bool)
	n := b.getNode(b.Start, Point{b.Start.X, b.Start.Y + 1}, p2c)
	l, err := b.longestPath(visited, n, p2c)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func (b Board) longestPath(visited map[Point]bool, n node, criterion func(at int32, cur, next Point) bool) (int, error) {
	if n.goal == true {
		return n.len, nil
	}
	best := -1
	vc := make(map[Point]bool)
	for k := range visited {
		vc[k] = true
	}
	for _, next := range n.next {
		vc[next] = true
	}
	for _, next := range n.next {
		if visited[next] {
			continue
		}
		nn := b.getNode(n.end, next, criterion)
		k, err := b.longestPath(vc, nn, criterion)
		if err != nil {
			continue
		}
		best = max(k, best)
	}
	if best == -1 {
		return 0, errors.New("dead end")
	}
	return best + n.len, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ParseBoard(lines []string) Board {
	tiles := make([][]int32, len(lines))
	for i, l := range lines {
		tiles[i] = []int32(l)
	}
	var start Point
	for x, c := range tiles[0] {
		if c == '.' {
			start = Point{x, 0}
			break
		}
	}
	var end Point
	for x, c := range tiles[len(tiles)-1] {
		if c == '.' {
			end = Point{x, len(tiles) - 1}
		}
	}
	return Board{
		Tiles: tiles,
		Start: start,
		Goal:  end,
		MaxY:  len(tiles) - 1,
		MaxX:  len(tiles[0]) - 1,
	}
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
