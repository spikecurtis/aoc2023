package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d18"
)

func main() {
	dp := d18.ParseDigPlan(d18.GetInput("./input"), 1)
	fmt.Println(dp.CountInterior2())
	dp = d18.ParseDigPlan(d18.GetInput("./input"), 2)
	fmt.Println(dp.CountInterior2())
}
