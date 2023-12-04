package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type card struct {
	number  int
	winners map[int]bool
	have    map[int]bool
}

func (c card) p1Score() int {
	k := c.matches()
	if k == 0 {
		return 0
	}
	return 1 << (k - 1)
}

func (c card) matches() int {
	k := 0
	for n := range c.have {
		if c.winners[n] {
			k++
		}
	}
	return k
}

func main() {
	cards := make([]card, 0)
	for _, line := range getInput() {
		cards = append(cards, parseCard(line))
	}
	p1 := 0
	for _, c := range cards {
		p1 += c.p1Score()
	}
	fmt.Println(p1)
	p2 := p2Score(cards)
	fmt.Println(p2)
}

func parseCard(l string) card {
	numS := strings.Trim(l[4:8], " ")
	cn, err := strconv.Atoi(numS)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(l[9:], "|")
	winners := parseInts(parts[0])
	have := parseInts(parts[1])
	return card{
		number:  cn,
		winners: winners,
		have:    have,
	}
}

func parseInts(l string) map[int]bool {
	out := make(map[int]bool)
	l = strings.Trim(l, " ")
	for _, s := range strings.Split(l, " ") {
		s = strings.Trim(s, " ")
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		out[n] = true
	}
	return out
}

func p2Score(cards []card) int {
	copies := make([]int, len(cards)+1)
	for i := 1; i <= len(cards); i++ {
		copies[i] += 1
	}
	for _, c := range cards {
		m := c.matches()
		i := copies[c.number]
		for j := 0; j < m; j++ {
			copies[c.number+j+1] += i
		}
	}
	p2 := 0
	for _, k := range copies {
		p2 += k
	}
	return p2
}

func getInput() []string {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
