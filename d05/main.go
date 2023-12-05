package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	a := parseAlmanac(getInput())
	fmt.Println(a.part1())
	fmt.Println(a.part2())
}

type mapRange struct {
	dst, src, len int
}

type sMap []mapRange

func (m sMap) get(x int) int {
	for _, mr := range m {
		if x < mr.src {
			continue
		}
		if x > mr.src+mr.len-1 {
			continue
		}
		return x - mr.src + mr.dst
	}
	return x
}

type almanac struct {
	seeds []int
	maps  []sMap
}

func (a almanac) location(seed int) int {
	x := seed
	for _, m := range a.maps {
		x = m.get(x)
	}
	return x
}

func (a almanac) part1() int {
	bestLoc := math.MaxInt
	for _, seed := range a.seeds {
		loc := a.location(seed)
		if loc < bestLoc {
			bestLoc = loc
		}
	}
	return bestLoc
}

func (a almanac) part2() int {
	bestLoc := math.MaxInt
	for i := 0; i < len(a.seeds); i += 2 {
		start := a.seeds[i]
		for j := 0; j < a.seeds[i+1]; j++ {
			seed := start + j
			loc := a.location(seed)
			if loc < bestLoc {
				bestLoc = loc
			}
		}
	}
	return bestLoc
}

func parseAlmanac(lines []string) almanac {
	seedS := parseInts(lines[0][7:])
	maps := make([]sMap, 0)
	var curMap sMap
	for _, l := range lines[2:] {
		ll := strings.Trim(l, " ")
		if strings.HasSuffix(ll, ":") {
			curMap = make(sMap, 0)
			continue
		}
		if ll == "" {
			maps = append(maps, curMap)
			continue
		}
		nums := parseInts(ll)
		curMap = append(curMap, mapRange{dst: nums[0], src: nums[1], len: nums[2]})
	}
	maps = append(maps, curMap)
	return almanac{
		seeds: seedS,
		maps:  maps,
	}
}

func parseInts(l string) []int {
	out := make([]int, 0)
	l = strings.Trim(l, " ")
	for _, s := range strings.Split(l, " ") {
		s = strings.Trim(s, " ")
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, n)
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
