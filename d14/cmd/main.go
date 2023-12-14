package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d14"
)

func main() {
	p := d14.ParsePlatform(d14.GetInput("./input"))
	p.TiltNorth()
	fmt.Println(p.P1Load())
	//fmt.Println(p.String())
	fmt.Println("Part 2:")

	p = d14.ParsePlatform(d14.GetInput("./input"))
	p.Cycle(1000000000)
	//fmt.Println(p.String())
	fmt.Println(p.P1Load())
}
