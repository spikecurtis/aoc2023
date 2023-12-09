package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	seqs := make([]seq, 0)
	for _, line := range getInput() {
		seqs = append(seqs, parseInts(line))
	}
	p1 := 0
	p2 := 0
	for _, s := range seqs {
		p1 += s.next()
		p2 += s.prev()
	}
	fmt.Println(p1)
	fmt.Println(p2)
}

type seq []int

func (s seq) next() int {
	if s.constant() {
		return s[0]
	}
	d := s.diff()
	return s[len(s)-1] + d.next()
}

func (s seq) prev() int {
	if s.constant() {
		return s[0]
	}
	d := s.diff()
	return s[0] - d.prev()
}

func (s seq) diff() seq {
	d := make(seq, len(s)-1)
	for i := 0; i < len(s)-1; i++ {
		d[i] = s[i+1] - s[i]
	}
	return d
}

func (s seq) constant() bool {
	k := s[0]
	for _, x := range s {
		if x != k {
			return false
		}
	}
	return true
}

func parseInts(l string) []int {
	out := make([]int, 0)
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
		out = append(out, n)
	}
	return out
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
