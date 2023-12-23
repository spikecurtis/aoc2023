package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d21"
)

func main() {
	b := d21.ParseBoard(d21.GetInput("./input"))
	fmt.Println(len(b.PlotsForSteps(64)))
	fmt.Println(b.Part2(26501365))
}
