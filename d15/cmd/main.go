package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d15"
)

func main() {
	tokens := d15.GetInput("./input")
	p1 := 0
	for _, t := range tokens {
		p1 += d15.HASH(t)
	}
	fmt.Println(p1)
	var h d15.HASHMAP
	for _, t := range tokens {
		h.Do(t)
	}
	fmt.Println(h.FocusingPower())
}
