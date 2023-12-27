package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d22"
)

func main() {
	bricks := d22.ParseBricks(d22.GetInput("./input"))
	bricks = d22.DropBricks(bricks)
	supports := d22.SupportedBy(bricks)
	fmt.Println(d22.CountDisintegrate(supports, len(bricks)))
	fmt.Println(d22.CountTotalChainDisintegrate(supports, len(bricks)))
}
