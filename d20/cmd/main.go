package main

import (
	"fmt"
	"github.com/spikecurtis/aoc2023/d20"
)

func main() {
	mods := d20.ParseModules(d20.GetInput("./input"))
	//d, err := os.Create("./input.dot")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//d.WriteString(d20.DOT(mods))
	//d.Close()
	lows, highs := d20.PressButton(mods, 1000)
	fmt.Println(lows * highs)
	mods = d20.ParseModules(d20.GetInput("./input"))
	fmt.Println(d20.LowToRx(mods))
}
