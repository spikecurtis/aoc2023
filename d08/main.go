package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	g := parseGuide(getInput())
	p1 := g.followToZZZ()
	fmt.Println(p1)
	p2 := g.ghostFollowToZ()
	fmt.Println(p2)
}

type direction int

const (
	left direction = iota
	right
)

type guide struct {
	instructions []direction
	network      *node
	startNodes   []*node
}

type node struct {
	name string
	next [2]*node
}

func parseGuide(lines []string) guide {
	instructions := make([]direction, 0)
	for _, c := range lines[0] {
		switch c {
		case 'L':
			instructions = append(instructions, left)
		case 'R':
			instructions = append(instructions, right)
		default:
			log.Fatal(c)
		}
	}
	nodes := make(map[string]*node)
	startNodes := make([]*node, 0)
	for _, line := range lines[2:] {
		nn := line[:3]
		ln := line[7:10]
		rn := line[12:15]
		n, ok := nodes[nn]
		if !ok {
			n = &node{name: nn}
			nodes[nn] = n
		}
		l, ok := nodes[ln]
		if !ok {
			l = &node{name: ln}
			nodes[ln] = l
		}
		n.next[left] = l
		r, ok := nodes[rn]
		if !ok {
			r = &node{name: rn}
			nodes[rn] = r
		}
		n.next[right] = r
		if strings.HasSuffix(nn, "A") {
			startNodes = append(startNodes, n)
		}
	}
	return guide{
		instructions: instructions,
		network:      nodes["AAA"],
		startNodes:   startNodes,
	}
}

type start struct {
	phase int
	name  string
}

type visit struct {
	steps int
	phase int
}

type cycle struct {
	steps int
	phase int
}

func (g guide) followToZZZ() int {
	steps := 0
	cur := g.network
	lastVisit := make(map[string]visit)
	cycles := make(map[start]cycle)
	lastEmit := -1
	for {
		for phase := 0; phase < len(g.instructions); phase++ {
			if v, ok := lastVisit[cur.name]; ok {
				// cycle!

				cycles[start{v.phase, cur.name}] = cycle{steps - v.steps, phase}
			}
			lastVisit[cur.name] = visit{steps, phase}
			if c, ok := cycles[start{phase, cur.name}]; ok {
				phase = c.phase
				steps += c.steps
				delete(lastVisit, cur.name)
			}
			if steps-lastEmit > 1000 {
				//fmt.Printf("took steps: %d, cur node: %s, num cycles: %d\n", steps, cur.name, len(cycles))
				lastEmit = steps
			}
			if cur.name == "ZZZ" {
				return steps
			}
			cur = cur.next[g.instructions[phase]]
			steps++
		}
	}
}

func (g guide) ghostFollowToZ() int {
	var fs []fullCycle
	for _, ghost := range g.startNodes {
		fs = append(fs, g.findFullCycle(ghost))
	}
	longestL := 0
	var longestFS fullCycle
	for _, f := range fs {
		if f.length > longestL {
			longestFS = f
			longestL = f.length
		}
	}
	if len(longestFS.goalPhases) != 1 {
		log.Fatal("uh oh")
	}
	steps := longestFS.stepsToEnter + longestFS.goalPhases[0]
	for {
		if !longestFS.isGoal(steps) {
			log.Fatal("bad cycle math")
		}
		goal := true
		for _, f := range fs {
			if !f.isGoal(steps) {
				goal = false
				break
			}
		}
		if goal {
			return steps
		}
		steps += longestFS.length
	}
	return 0
}

type fullCycle struct {
	stepsToEnter int
	length       int
	goalPhases   []int
}

func (f fullCycle) isGoal(steps int) bool {
	phase := (steps - f.stepsToEnter) % f.length
	for _, gp := range f.goalPhases {
		if phase == gp {
			return true
		}
	}
	return false
}

func (g guide) findFullCycle(s *node) fullCycle {
	visits := make(map[start]int)
	steps := 0
	goalSteps := make(map[string]int)
	for {
		for phase, d := range g.instructions {
			if vs, ok := visits[start{phase, s.name}]; ok {
				// found cycle
				stepsToEnter := vs
				length := steps - stepsToEnter
				goalPhases := make([]int, 0)
				for _, gs := range goalSteps {
					if gs < stepsToEnter {
						continue
					}
					goalPhases = append(goalPhases, gs-stepsToEnter)
				}
				return fullCycle{
					stepsToEnter: stepsToEnter,
					length:       length,
					goalPhases:   goalPhases,
				}
			}
			visits[start{phase, s.name}] = steps
			if strings.HasSuffix(s.name, "Z") {
				goalSteps[s.name] = steps
			}
			s = s.next[d]
			steps++
		}
	}
}

func ghostsOnZ(ghosts []*node) bool {
	for _, g := range ghosts {
		if !strings.HasSuffix(g.name, "Z") {
			return false
		}
	}
	return true
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
