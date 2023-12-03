package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	s := getSchematic()
	p1PNs := s.p1PartNumbers()
	p1 := 0
	for _, pn := range p1PNs {
		p1 += pn.number
	}
	fmt.Println(p1)
	gears := s.gears()
	p2 := 0
	for _, g := range gears {
		p2 += g.ratio
	}
	fmt.Println(p2)
}

type point struct {
	x, y int
}

func (p point) adjacentPoints(maxX, maxY int) []point {
	points := []point{}
	for y := max(0, p.y-1); y < min(p.y+2, maxY); y++ {
		for x := max(0, p.x-1); x < min(p.x+2, maxX); x++ {
			points = append(points, point{x, y})
		}
	}
	return points
}

type partNumber struct {
	number     int
	start, end point
}

func (p partNumber) adjacentPoints(maxX, maxY int) []point {
	points := []point{}
	for y := max(0, p.start.y-1); y < min(p.end.y+2, maxY); y++ {
		for x := max(0, p.start.x-1); x < min(p.end.x+1, maxX); x++ {
			points = append(points, point{x, y})
		}
	}
	return points
}

func (p partNumber) within(pnt point) bool {
	if pnt.y != p.start.y {
		return false
	}
	if pnt.x >= p.end.x {
		return false
	}
	if pnt.x < p.start.x {
		return false
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type schematic struct {
	partNumbers []partNumber
	symbols     map[point]rune
	maxX        int
	maxY        int
}

type gear struct {
	ratio    int
	location point
}

func (s schematic) p1PartNumbers() []partNumber {
	out := []partNumber{}
parts:
	for _, pn := range s.partNumbers {
		for _, pnt := range pn.adjacentPoints(s.maxX, s.maxY) {
			if _, ok := s.symbols[pnt]; ok {
				out = append(out, pn)
				continue parts
			}
		}
	}
	return out
}

func (s schematic) gears() []gear {
	out := []gear{}
	for loc, m := range s.symbols {
		if m != '*' {
			continue
		}
		pns := []partNumber{}
	parts:
		for _, pn := range s.partNumbers {
			for _, pnt := range loc.adjacentPoints(s.maxX, s.maxY) {
				if pn.within(pnt) {
					pns = append(pns, pn)
					continue parts
				}
			}
		}
		if len(pns) == 2 {
			out = append(out, gear{
				ratio:    pns[0].number * pns[1].number,
				location: loc,
			})
		}
	}
	return out
}

var digits = map[rune]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
}

func getSchematic() schematic {
	lines := getInput()
	s := schematic{
		partNumbers: make([]partNumber, 0),
		symbols:     make(map[point]rune),
		maxY:        len(lines),
		maxX:        len(lines[0]),
	}
	for y, line := range lines {
		parseLine(line, y, &s)
	}
	return s
}

func parseLine(line string, y int, s *schematic) {
	isPn := false
	curPn := []rune{}
	var start point
	rl := []rune(line)
	for x, r := range rl {
		if digits[r] {
			curPn = append(curPn, r)
			if !isPn {
				start = point{x, y}
			}
			isPn = true
			continue
		} else {
			if isPn {
				pn, err := strconv.ParseInt(string(curPn), 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				s.partNumbers = append(s.partNumbers, partNumber{
					number: int(pn),
					start:  start,
					end:    point{x, y},
				})
			}
			isPn = false
			curPn = []rune{}
		}
		if r == '.' {
			continue
		}
		s.symbols[point{x, y}] = r
	}
	if isPn {
		pn, err := strconv.ParseInt(string(curPn), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		s.partNumbers = append(s.partNumbers, partNumber{
			number: int(pn),
			start:  start,
			end:    point{len(rl), y},
		})
	}
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
