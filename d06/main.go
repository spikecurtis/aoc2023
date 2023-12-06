package main

import (
	"fmt"
	"math"
)

func main() {
	races := []race{
		//{7, 9}, // unit test 😂
		{55, 246},
		{82, 1441},
		{64, 1012},
		{90, 1111},
	}
	p1 := 1
	for _, r := range races {
		p1 *= r.waysToWin()
	}
	fmt.Println(p1)
	fmt.Println(race{55826490, 246144110121111}.waysToWin())
}

type race struct {
	time   int
	record int
}

func (r race) waysToWin() int {
	b := float64(r.time)
	c := float64(-r.record)
	h := (-b - math.Sqrt(math.Pow(b, 2)+4*c)) / -2
	l := (-b + math.Sqrt(math.Pow(b, 2)+4*c)) / -2
	w := math.Floor(h) - math.Ceil(l) + 1
	return int(w)
}
