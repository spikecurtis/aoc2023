package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d23"
)

func main() {
	board := d23.ParseBoard(d23.GetInput("./input"))
	fmt.Println(board.LongestPath1())
	fmt.Println(board.LongestPath2())
}
