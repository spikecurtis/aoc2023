package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input := getInput()
	p1 := 0
	p2 := 0
	for _, l := range input {
		p1 += parseCalibrationP1(l)
		p2 += parseCalibrationP2(l)
	}
	fmt.Println(p1)
	fmt.Println(p2)
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

func parseCalibrationP1(l string) int {
	first := -1
	last := -1
	for _, c := range l {
		if c < 0x30 {
			continue
		}
		if c > 0x39 {
			continue
		}
		d := int(c - 0x30)
		if first == -1 {
			first = d
		}
		last = d
	}
	return first*10 + last
}

var digits = map[string]int{
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func parseCalibrationP2(l string) int {
	first := -1
	var err error
	for s, _ := range l {
		first, err = findPrefix(l[s:])
		if err == nil {
			break
		}
	}
	last := -1
	for s := len(l); s > 0; s-- {
		last, err = findSuffix(l[:s])
		if err == nil {
			break
		}
	}
	return first*10 + last
}

func findPrefix(l string) (int, error) {
	for t, i := range digits {
		if strings.HasPrefix(l, t) {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}

func findSuffix(l string) (int, error) {
	for t, i := range digits {
		if strings.HasSuffix(l, t) {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}
