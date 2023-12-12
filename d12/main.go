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
	rt := row{
		groups:      [][]condition{{unknown, unknown, unknown}, {broken, broken, broken}},
		groupsSizes: []int{1, 1, 3},
	}
	if rt.arrangements() != 1 {
		log.Fatal("not working")
	}
	rows := make([]row, 0)
	for _, line := range getInput("./input") {
		rows = append(rows, parseRow(line))
	}
	p1 := 0
	for _, r := range rows {
		p1 += r.arrangements()
	}
	fmt.Println(p1)
	// Pt 2
	rows = make([]row, 0)
	for _, line := range getInput("./input") {
		rows = append(rows, parseRow2(line))
	}
	p2 := 0
	for _, r := range rows {
		p2 += r.arrangements()
	}
	fmt.Println(p2)
}

type condition int

const (
	operational condition = iota
	broken
	unknown
)

type row struct {
	groups      [][]condition
	groupsSizes []int
}

var lookup = map[int32]condition{
	'?': unknown,
	'.': operational,
	'#': broken,
}

func parseRow(line string) row {
	parts := strings.Split(line, " ")
	groups := make([][]condition, 0)
	var curGroup []condition
	for _, c := range parts[0] {
		cnd := lookup[c]
		switch cnd {
		case operational:
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
			}
			curGroup = nil
		default:
			curGroup = append(curGroup, cnd)
		}
	}
	if len(curGroup) > 0 {
		groups = append(groups, curGroup)
	}
	groupSizes := make([]int, 0)
	gts := strings.Split(parts[1], ",")
	for _, gt := range gts {
		i, err := strconv.Atoi(gt)
		if err != nil {
			log.Fatal(err)
		}
		groupSizes = append(groupSizes, i)
	}
	return row{
		groups:      groups,
		groupsSizes: groupSizes,
	}
}

func parseRow2(line string) row {
	parts := strings.Split(line, " ")
	newline := fmt.Sprintf("%s?%s?%s?%s?%s %s,%s,%s,%s,%s",
		parts[0], parts[0], parts[0], parts[0], parts[0],
		parts[1], parts[1], parts[1], parts[1], parts[1],
	)
	return parseRow(newline)
}

func (r row) arrangements() int {
	// base cases
	if len(r.groups) == 0 {
		if len(r.groupsSizes) == 0 {
			return 1
		}
		return 0
	}
	if len(r.groupsSizes) == 0 {
		// only way for this is all unknowns in all groups
		for _, g := range r.groups {
			for _, cnd := range g {
				if cnd != unknown {
					return 0
				}
			}
		}
		return 1
	}
	if len(r.groups) == 1 {
		return singleUnknownGroup(r.groups[0], r.groupsSizes)
	}

	// recursive check
	split := len(r.groups) / 2
	lg := r.groups[:split]
	rg := r.groups[split:]
	ways := 0
	for k := 0; k <= len(r.groupsSizes); k++ {
		lgs := r.groupsSizes[:k]
		rgs := r.groupsSizes[k:]
		lr := row{lg, lgs}
		rr := row{rg, rgs}
		// compute the smaller one first in case it's 0, then we can skip
		if len(lgs) < len(rgs) {
			la := lr.arrangements()
			if la == 0 {
				continue
			}
			ways += la * rr.arrangements()
			continue
		}
		ra := rr.arrangements()
		if ra == 0 {
			continue
		}
		ways += ra * lr.arrangements()
	}
	return ways
}

func printSingleUnknown(group []condition, sizes []int) string {
	out := strings.Builder{}
	for _, cnd := range group {
		out.WriteString(fmt.Sprintf("%d", cnd))
	}
	out.WriteString("|")
	for _, s := range sizes {
		out.WriteString(fmt.Sprintf("%d", s))
		out.WriteString(",")
	}
	return out.String()
}

var cache map[string]int

func singleUnknownGroup(group []condition, sizes []int) (result int) {
	if cache == nil {
		cache = make(map[string]int)
	}
	key := printSingleUnknown(group, sizes)
	r, ok := cache[key]
	if ok {
		return r
	}
	defer func() {
		cache[key] = result
	}()
	if len(group) == 0 {
		if len(sizes) == 0 {
			return 1
		}
		return 0
	}
	minConditions := len(sizes) - 1
	for _, s := range sizes {
		minConditions += s
	}
	if len(group) < minConditions {
		return 0
	}
	if group[0] == broken {
		if len(sizes) == 0 {
			return 0
		}
		// comp off the first group
		group = group[sizes[0]:]
		sizes = sizes[1:]
		if len(sizes) > 0 {
			if len(group) < 2 {
				// not enough room for another group, since we need at least
				// one operational
				return 0
			}
			if group[0] == broken {
				// need an operational between
				return 0
			}
			return singleUnknownGroup(group[1:], sizes)
		}
		return singleUnknownGroup(group, sizes)
	}
	ways := 0
	// assume group[0] is operational
	ways += singleUnknownGroup(group[1:], sizes)
	// assume group[0] is broken
	newGroup := make([]condition, len(group))
	copy(newGroup, group)
	newGroup[0] = broken
	ways += singleUnknownGroup(newGroup, sizes)
	return ways
}

func getInput(name string) []string {
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
