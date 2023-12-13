package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d13"
	"log"
)

func main() {
	lines := d13.GetInput("./input")
	patterns := d13.ParsePatterns(lines)
	p1 := 0
	p2 := 0
	for _, p := range patterns {
		score, err := p.P1Score()
		if err != nil {
			log.Fatal(err)
		}
		p1 += score
		score, err = p.P2Score()
		if err != nil {
			log.Fatal(err)
		}
		p2 += score
	}
	fmt.Println(p1)
	fmt.Println(p2)
}
