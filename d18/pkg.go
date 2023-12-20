package d18

import (
	"bufio"
	"log"
	"os"
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
	prev := cur

	border := 0
	shoelace := 0
	for _, t := range dp.Trenches {
		cur = cur.Move(t.D, t.Len)
		border += t.Len
		// "shoelace" formula computes the area from point coordinates
		// can be derived from Green's theorem
		shoelace += cur.y * prev.x
		shoelace -= cur.x * prev.y
		prev = cur
	}
	if shoelace < 0 {
		shoelace = -shoelace
	}
	shoelace = shoelace / 2
	// shoelace gives the area of the polygon defined by the center of the trenches,
	// but we want to count border + interior points, since we dig 1m x1m cubes
	// at all these points.
	// Pick's theorem: A = i + b/2 - 1 where
	// A = area (calculated by shoelace formula above)
	// b = border points
	// i = interior points
	//
	// Some algebra gives: i + b = A + b/2 + 1
	return shoelace + (border / 2) + 1
}

type Trench struct {
	D   Direction
	Len int
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
