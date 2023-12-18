package d17

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"strconv"
)

type Board [][]int

func ParseBoard(lines []string) Board {
	b := make(Board, len(lines))
	for y, line := range lines {
		b[y] = make([]int, len(line))
		for x, c := range line {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			b[y][x] = i
		}
	}
	return b
}

func GetInput(name string) []string {
	f, err := os.Open(name)
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

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Seg struct {
	x, y     int
	d        Direction
	straight int
}

type Mode int

const (
	Regular Mode = iota
	Ultra
)

func (b Board) MinLoss(m Mode) int {
	gX := len(b[0]) - 1
	gY := len(b) - 1
	h := func(s Seg) int {
		return abs(s.x, gX) + abs(s.y, gY)
	}
	open := make(map[Seg]bool)
	gLoss := make(map[Seg]int)
	fLoss := make(map[Seg]int)
	for _, s := range []Seg{
		{0, 0, East, 0},
		{0, 0, South, 0},
	} {
		open[s] = true
		gLoss[s] = 0
		fLoss[s] = h(s)
	}
	for len(open) > 0 {
		var cur Seg
		bestL := math.MaxInt
		for s := range open {
			l := fLoss[s]
			if l < bestL {
				cur = s
				bestL = l
			}
		}
		delete(open, cur)

		if cur.x == gX && cur.y == gY {
			if m == Regular || cur.straight >= 3 {
				return gLoss[cur]
			}
		}

		for _, s := range b.neighbors(cur, m) {
			tl := gLoss[cur] + b[s.y][s.x]
			bl, ok := gLoss[s]
			if !ok || tl < bl {
				gLoss[s] = tl
				fLoss[s] = tl + h(s)
				open[s] = true
			}
		}
	}
	log.Fatal("failed to find")
	return -1
}

func (b Board) neighbors(s Seg, m Mode) []Seg {
	r := (s.d + 1) % 4
	l := (s.d + 3) % 4
	n := make([]Seg, 0, 3)
	if m == Regular || s.straight >= 3 {
		for _, d := range []Direction{r, l} {
			x, y, err := b.step(s.x, s.y, d)
			if err != nil {
				continue
			}
			n = append(n, Seg{
				x: x,
				y: y,
				d: d,
			})
		}
	}
	if (m == Regular && s.straight < 2) || (m == Ultra && s.straight < 9) {
		x, y, err := b.step(s.x, s.y, s.d)
		if err == nil {
			n = append(n, Seg{
				x:        x,
				y:        y,
				d:        s.d,
				straight: s.straight + 1,
			})
		}
	}
	return n
}

var offBoard = errors.New("off board")

func (b Board) step(x, y int, d Direction) (nx, ny int, err error) {
	nx = x
	ny = y
	switch d {
	case North:
		ny -= 1
	case South:
		ny += 1
	case East:
		nx += 1
	case West:
		nx -= 1
	}
	if nx < 0 || ny < 0 {
		return 0, 0, offBoard
	}
	if nx > len(b[0])-1 || ny > len(b)-1 {
		return 0, 0, offBoard
	}
	return
}

func abs(a, b int) int {
	r := a - b
	if r < 0 {
		return -r
	}
	return r
}
