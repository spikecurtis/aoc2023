package d18

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type DigPlan struct {
	Trenches []Trench
}

type Point struct {
	x, y int
}

func (p Point) Move(d Direction, l int) Point {
	switch d {
	case Up:
		return Point{p.x, p.y - l}
	case Down:
		return Point{p.x, p.y + l}
	case Left:
		return Point{p.x - l, p.y}
	case Right:
		return Point{p.x + l, p.y}
	}
	panic("unknown direction")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func (dp DigPlan) CountInterior() int {
	cur := Point{0, 0}
	dug := map[Point]bool{cur: true}
	prevD := dp.Trenches[len(dp.Trenches)-1].D
	cirality := 0
	for _, t := range dp.Trenches {
		c := t.D - prevD
		if c < -1 {
			c += 4
		}
		if c > 1 {
			c -= 4
		}
		if c < -1 || c > 1 {
			panic("cirality assumption violated")
		}
		cirality += int(c)
		prevD = t.D
		for i := 0; i < t.Len; i++ {
			cur = cur.Move(t.D, 1)
			dug[cur] = true
		}
	}
	if cur.x != 0 || cur.y != 0 {
		panic("didn't reach origin")
	}
	if cirality > 0 {
		cirality = 1
	} else {
		cirality = -1
	}

	// add adjacent interior points to open set
	open := make(map[Point]bool)
	for _, t := range dp.Trenches {
		for i := 0; i < t.Len; i++ {
			if i > 1 && i < t.Len-1 {
				d := (t.D + Direction(cirality) + 4) % 4
				open[cur.Move(d, 1)] = true
			}
			cur = cur.Move(t.D, 1)
		}
	}

	// fill
	pop := func() Point {
		for p := range open {
			delete(open, p)
			return p
		}
		panic("empty")
	}
	for len(open) > 0 {
		f := pop()
		dug[f] = true
		for _, d := range []Direction{Up, Down, Left, Right} {
			fp := f.Move(d, 1)
			if !dug[fp] {
				open[fp] = true
			}
		}
	}
	return len(dug)
}

func (dp DigPlan) CountInterior2() int {
	cur := Point{0, 0}
	xm := make(map[int]bool)
	ym := make(map[int]bool)

	segments := make([]Segment, 0)
	for _, t := range dp.Trenches {
		xm[cur.x] = true
		ym[cur.y] = true
		x0 := cur.x
		y0 := cur.y
		cur = cur.Move(t.D, t.Len)
		segments = append(segments, Segment{
			Trench: t,
			Xmin:   min(cur.x, x0),
			Xmax:   max(cur.x, x0),
			Ymin:   min(cur.y, y0),
			Ymax:   max(cur.y, y0),
		})
	}
	xs := make([]int, 0, len(xm))
	for x := range xm {
		xs = append(xs, x)
	}
	sort.Ints(xs)
	xs = append(xs, xs[len(xs)-1])
	ys := make([]int, 0, len(ym))
	for y := range ym {
		ys = append(ys, y)
	}
	sort.Ints(ys)
	ys = append(ys, ys[len(ys)-1])

	interior := 0
	for i := 0; i < len(xs)-1; i++ {
		x0 := xs[i]
		x1 := xs[i+1]

		for j := 0; j < len(ys)-1; j++ {
			y0 := ys[j]
			y1 := ys[j+1]
			sq := Square{x0, x1, y0, y1}
			interior += sq.CountInterior(segments)
		}
	}
	return interior
}

type Square struct {
	x0, x1 int
	y0, y1 int
}

// CountInterior counts the number of interior tiles in a Square.
//
// We only count the top-left corner, top, left and middle of the square, leaving
// off the bottom-left, bottom-right and top-right corners, the bottom, and the right
// this allows us to avoid double counting, since those tiles are in other squares in the array
//
// *--------*-----*--
// |        |     |
// *--------*-----*--
// |        |     |
func (s Square) CountInterior(segs []Segment) int {
	interior := 0
	mx := (s.x1-s.x0)/2 + s.x0
	my := (s.y1-s.y0)/2 + s.y0
	if IsInterior(Point{s.x0, s.y0}, segs) {
		interior++
	}
	// Left
	if my != s.y0 && IsInterior(Point{s.x0, my}, segs) {
		interior += s.y1 - s.y0 - 1
	}
	// Top
	if mx != s.x0 && IsInterior(Point{mx, s.y0}, segs) {
		interior += s.x1 - s.x0 - 1
	}
	// Middle
	if mx != s.x0 && my != s.y0 && IsInterior(Point{mx, my}, segs) {
		interior += (s.x1 - s.x0 - 1) * (s.y1 - s.y0 - 1)
	}
	return interior
}

type Trench struct {
	D   Direction
	Len int
}

type Segment struct {
	Trench
	Xmin, Xmax int
	Ymin, Ymax int
}

func IsInterior(p Point, segs []Segment) bool {
	wn := 0
	for _, seg := range segs {
		if seg.Contains(p) {
			return true
		}
		wn += WindingNumberInc(p, seg)
	}
	return wn != 0
}

func (s Segment) Contains(p Point) bool {
	return p.x >= s.Xmin && p.x <= s.Xmax && p.y >= s.Ymin && p.y <= s.Ymax
}

func WindingNumberInc(p Point, s Segment) int {
	if s.D != Up && s.D != Down {
		return 0
	}
	if p.x >= s.Xmin {
		return 0
	}
	if p.y <= s.Ymin {
		return 0
	}
	if p.y > s.Ymax {
		return 0
	}
	if s.D == Up {
		return 1
	}
	return -1
}

func ParseDigPlan(lines []string, part int) DigPlan {
	dp := DigPlan{}
	for _, line := range lines {
		dp.Trenches = append(dp.Trenches, ParseTrench(line, part))
	}
	return dp
}

func ParseTrench(line string, part int) Trench {
	if part == 1 {
		return ParseTrench1(line)
	}
	return ParseTrench2(line)
}

func ParseTrench1(line string) Trench {
	parts := strings.Split(line, " ")
	t := Trench{}
	switch parts[0] {
	case "U":
		t.D = Up
	case "D":
		t.D = Down
	case "L":
		t.D = Left
	case "R":
		t.D = Right
	}
	l, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	t.Len = l
	return t
}

func ParseTrench2(line string) Trench {
	parts := strings.Split(line, " ")
	t := Trench{}
	l, err := strconv.ParseInt(parts[2][2:7], 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	t.Len = int(l)
	switch parts[2][7] {
	case '0':
		t.D = Right
	case '1':
		t.D = Down
	case '2':
		t.D = Left
	case '3':
		t.D = Up
	}
	return t
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
