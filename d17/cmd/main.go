package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d17"
)

func main() {
	b := d17.ParseBoard(d17.GetInput("./input"))
	fmt.Println(b.MinLoss(d17.Regular))
	fmt.Println(b.MinLoss(d17.Ultra))
}
