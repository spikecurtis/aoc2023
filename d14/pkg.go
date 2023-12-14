package d14

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Platform struct {
	loc  [][]int32
	maxX int
	maxY int
}

func (p *Platform) P1Load() int {
	t := 0
	for y, row := range p.loc {
		w := p.maxY - y + 1
		for _, c := range row {
			if c == 'O' {
				t += w
			}
		}
	}
	return t
}

func (p *Platform) Cycle(n int) {
	states := make(map[string]int)
	for c := 0; c < n; c++ {
		if c%100000 == 0 {
			fmt.Printf("cycle No %d\n", c)
		}
		ss := p.String()
		if prevC, ok := states[ss]; ok {
			// super cycle!
			l := c - prevC
			fmt.Printf("found super cycle of length %d at cycle %d\n", l, c)
			left := n - c
			k := left / l
			c += k * l
			fmt.Printf("fast forwarded to cycle %d\n", c)
		} else {
			states[ss] = c
		}
		p.TiltNorth()
		p.TiltWest()
		p.TiltSouth()
		p.TiltEast()
	}
}

func (p *Platform) TiltNorth() {
	for y := 1; y <= p.maxY; y++ {
		for x := 0; x <= p.maxX; x++ {
			if p.loc[y][x] == 'O' {
				p.loc[y][x] = '.'
				xf, yf := p.roll(x, y, 0, -1)
				p.loc[yf][xf] = 'O'
			}
		}
	}
}

func (p *Platform) TiltSouth() {
	for y := p.maxY; y >= 0; y-- {
		for x := 0; x <= p.maxX; x++ {
			if p.loc[y][x] == 'O' {
				p.loc[y][x] = '.'
				xf, yf := p.roll(x, y, 0, 1)
				p.loc[yf][xf] = 'O'
			}
		}
	}
}

func (p *Platform) TiltEast() {
	for x := p.maxX; x >= 0; x-- {
		for y := 0; y <= p.maxY; y++ {
			if p.loc[y][x] == 'O' {
				p.loc[y][x] = '.'
				xf, yf := p.roll(x, y, 1, 0)
				p.loc[yf][xf] = 'O'
			}
		}
	}
}

func (p *Platform) TiltWest() {
	for x := 0; x <= p.maxX; x++ {
		for y := 0; y <= p.maxY; y++ {
			if p.loc[y][x] == 'O' {
				p.loc[y][x] = '.'
				xf, yf := p.roll(x, y, -1, 0)
				p.loc[yf][xf] = 'O'
			}
		}
	}
}

func (p Platform) roll(x, y, dx, dy int) (xf, yf int) {
	xf = x
	yf = y
	for {
		if xf+dx < 0 || xf+dx > p.maxX || yf+dy < 0 || yf+dy > p.maxY {
			return
		}
		switch p.loc[yf+dy][xf+dx] {
		case 'O', '#':
			return
		default:
			yf += dy
			xf += dx
		}
	}
}

func (p Platform) String() string {
	b := strings.Builder{}
	for _, row := range p.loc {
		for _, c := range row {
			b.WriteRune(c)
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func ParsePlatform(lines []string) Platform {
	maxX := len(lines[0]) - 1
	maxY := len(lines) - 1
	var loc [][]int32
	for _, line := range lines {
		var row []int32
		for _, c := range line {
			row = append(row, c)
		}
		loc = append(loc, row)
	}
	return Platform{
		loc:  loc,
		maxX: maxX,
		maxY: maxY,
	}
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
