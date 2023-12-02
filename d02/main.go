package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := getInput()
	games := make([]game, 0, 100)
	for _, l := range lines {
		games = append(games, parseGame(l))
	}
	p1 := int64(0)
	p2 := int64(0)
	for _, g := range games {
		if g.p1Possible() {
			p1 += g.n
		}
		p2 += g.minimumSet().power()
	}
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
}

type game struct {
	n    int64
	sets []set
}

type set struct {
	red, green, blue int64
}

func parseGame(l string) game {
	parts := strings.Split(l, ":")
	n, err := strconv.ParseInt(parts[0][5:], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	setStrings := strings.Split(parts[1], ";")
	sets := []set{}
	for _, ss := range setStrings {
		sets = append(sets, parseSet(ss))
	}
	return game{n: n, sets: sets}
}

func parseSet(s string) set {
	out := set{}
	cubes := strings.Split(s, ",")
	for _, cube := range cubes {
		cube = strings.Trim(cube, " ")
		parts := strings.Split(cube, " ")
		k, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		switch parts[1] {
		case "red":
			out.red = k
		case "green":
			out.green = k
		case "blue":
			out.blue = k
		default:
			log.Fatalf("unknown color: %s", parts[1])
		}
	}
	return out
}

func (g game) p1Possible() bool {
	for _, s := range g.sets {
		if s.red > 12 {
			return false
		}
		if s.green > 13 {
			return false
		}
		if s.blue > 14 {
			return false
		}
	}
	return true
}

func (g game) minimumSet() set {
	minSet := set{}
	for _, s := range g.sets {
		minSet.red = max(minSet.red, s.red)
		minSet.blue = max(minSet.blue, s.blue)
		minSet.green = max(minSet.green, s.green)
	}
	return minSet
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (s set) power() int64 {
	return s.green * s.blue * s.red
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
