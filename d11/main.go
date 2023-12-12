package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	u := parseGalaxies(getInput())
	u1 := u.expand(1)
	fmt.Println(u1.sumOfShortestPaths())
	u2 := u.expand(999999)
	fmt.Println(u2.sumOfShortestPaths())
}

type point struct {
	x, y int
}

type universe struct {
	galaxies   []point
	maxX, maxY int
}

func (u universe) expand(c int) universe {
	cols := make([]int, 0)
column:
	for x := 0; x <= u.maxX; x++ {
		for _, g := range u.galaxies {
			if g.x == x {
				continue column
			}
		}
		cols = append(cols, x)
	}
	rows := make([]int, 0)
row:
	for y := 0; y <= u.maxY; y++ {
		for _, g := range u.galaxies {
			if g.y == y {
				continue row
			}
		}
		rows = append(rows, y)
	}
	maxX := u.maxX + (len(cols) * c)
	maxY := u.maxY + (len(rows) * c)
	galaxies := make([]point, len(u.galaxies))
	for i, g := range u.galaxies {
		galaxies[i] = point{
			x: g.x + expansion(cols, g.x, c),
			y: g.y + expansion(rows, g.y, c),
		}
	}
	return universe{
		galaxies: galaxies,
		maxX:     maxX,
		maxY:     maxY,
	}
}

func (u universe) sumOfShortestPaths() int {
	sum := 0
	for i := 0; i < len(u.galaxies); i++ {
		for j := i + 1; j < len(u.galaxies); j++ {
			// shortest path is the "city block" distance
			sum += abs(u.galaxies[i].x, u.galaxies[j].x)
			sum += abs(u.galaxies[i].y, u.galaxies[j].y)
		}
	}
	return sum
}

func abs(i, j int) int {
	k := i - j
	if k < 0 {
		return -k
	}
	return k
}

func expansion(is []int, k, c int) int {
	e := 0
	for _, i := range is {
		if i < k {
			e += c
		}
	}
	return e
}

func parseGalaxies(lines []string) universe {
	galaxies := make([]point, 0)
	maxX := len(lines[0]) - 1
	maxY := len(lines) - 1
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				galaxies = append(galaxies, point{x, y})
			}
		}
	}
	return universe{
		galaxies: galaxies,
		maxX:     maxX,
		maxY:     maxY,
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
