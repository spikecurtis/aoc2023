package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d16"
)

func main() {
	b := d16.ParseBoard(d16.GetInput("./input"))
	b.Energize(d16.Beam{D: d16.Right, P: d16.Point{0, 0}})
	fmt.Println(b.CountEnergized())
	fmt.Println(b.BestEnergized())
}
