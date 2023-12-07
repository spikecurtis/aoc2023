package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var hnds hands
	for _, line := range getInput() {
		hnds = append(hnds, parseHand(line))
	}
	sort.Sort(hnds)
	p1 := 0
	for i, h := range hnds {
		p1 += (i + 1) * h.bid
	}
	fmt.Println(p1)
	var hnds2 handsP2
	for _, line := range getInput() {
		hnds2 = append(hnds2, parseHand(line))
	}
	sort.Sort(hnds2)
	p2 := 0
	for i, h := range hnds2 {
		p2 += (i + 1) * h.bid
	}
	fmt.Println(p2)
}

type hand struct {
	cards []int
	bid   int
}

type hands []hand

func (h hands) Len() int      { return len(h) }
func (h hands) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h hands) Less(i, j int) bool {
	a := h[i]
	b := h[j]
	at := a.getType()
	bt := b.getType()
	if at < bt {
		return true
	}
	if at > bt {
		return false
	}
	for k := 0; k < 5; k++ {
		if a.cards[k] < b.cards[k] {
			return true
		}
		if a.cards[k] > b.cards[k] {
			return false
		}
	}
	return false
}

type handsP2 []hand

func (h handsP2) Len() int      { return len(h) }
func (h handsP2) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h handsP2) Less(i, j int) bool {
	a := h[i]
	b := h[j]
	at := a.getTypeJokers()
	bt := b.getTypeJokers()
	if at < bt {
		return true
	}
	if at > bt {
		return false
	}
	for k := 0; k < 5; k++ {
		ak := a.cards[k]
		if ak == 11 {
			ak = 1
		}
		bk := b.cards[k]
		if bk == 11 {
			bk = 1
		}
		if ak < bk {
			return true
		}
		if ak > bk {
			return false
		}
	}
	return false
}

func (h hand) getType() int {
	labels := make(map[int]int)
	for _, card := range h.cards {
		labels[card] += 1
	}
	matches := make([]int, 0)
	for _, amt := range labels {
		matches = append(matches, amt)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(matches)))
	switch {
	case matches[0] == 5:
		return 6
	case matches[0] == 4:
		return 5
	case matches[0] == 3 && matches[1] == 2:
		return 4
	case matches[0] == 3:
		return 3
	case matches[0] == 2 && matches[1] == 2:
		return 2
	case matches[0] == 2:
		return 1
	default:
		return 0
	}
}

func (h hand) getTypeJokers() int {
	labels := make(map[int]int)
	jokers := 0
	for _, card := range h.cards {
		if card == 11 {
			jokers++
			continue
		}
		labels[card] += 1
	}
	matches := make([]int, 0)
	for _, amt := range labels {
		matches = append(matches, amt)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(matches)))
	if jokers == 5 {
		return 6
	}
	matches[0] += jokers
	switch {
	case matches[0] == 5:
		return 6
	case matches[0] == 4:
		return 5
	case matches[0] == 3 && matches[1] == 2:
		return 4
	case matches[0] == 3:
		return 3
	case matches[0] == 2 && matches[1] == 2:
		return 2
	case matches[0] == 2:
		return 1
	default:
		return 0
	}
}

var values = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

func parseHand(l string) hand {
	parts := strings.Split(l, " ")
	var cards []int
	runes := []rune(parts[0])
	for _, r := range runes {
		cards = append(cards, values[r])
	}
	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	return hand{cards, bid}
}

//func compareHand(a, b hand) int {
//
//}

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
