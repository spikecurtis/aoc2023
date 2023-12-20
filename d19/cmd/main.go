package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d19"
)

func main() {
	workflows, parts := d19.ParsePuzzle(d19.GetInput("./input"))
	p1 := 0
	for _, p := range parts {
		result := workflows.Process(p)
		if result == "A" {
			p1 += p.X + p.M + p.A + p.S
		}
	}
	fmt.Println(p1)
	fmt.Println(workflows.Combinations())
}
