package d13

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type Pattern struct {
	Rows []uint64
	Cols []uint64
}

type reflection struct {
	r     int
	isRow bool
}

func (r reflection) score() int {
	k := r.r + 1
	if r.isRow {
		k *= 100
	}
	return k
}

func (p Pattern) P1Score() (int, error) {
	for _, r := range p.findReflections() {
		return r.score(), nil
	}
	return 0, NoRefection
}

func (p Pattern) findReflections() []reflection {
	var out []reflection
	for _, r := range findReflection(p.Cols) {
		out = append(out, reflection{r: r})
	}
	for _, r := range findReflection(p.Rows) {
		out = append(out, reflection{r: r, isRow: true})
	}
	return out
}

func (p Pattern) P2Score() (int, error) {
	initReflections := p.findReflections()
	for x := 0; x < len(p.Cols); x++ {
		for y := 0; y < len(p.Rows); y++ {
			ps := p.smudge(x, y)
			newReflections := ps.findReflections()
			for _, nr := range newReflections {
				found := false
				for _, ir := range initReflections {
					if nr == ir {
						found = true
					}
				}
				if !found {
					return nr.score(), nil
				}
			}
		}
	}
	return 0, NoRefection
}

func (p Pattern) smudge(x, y int) Pattern {
	rows := make([]uint64, len(p.Rows))
	copy(rows, p.Rows)
	rows[y] ^= 1 << x
	cols := make([]uint64, len(p.Cols))
	copy(cols, p.Cols)
	cols[x] ^= 1 << y
	return Pattern{rows, cols}
}

func reflectedIndex(i, r int) int {
	return (2 * r) - i + 1
}

var NoRefection = errors.New("no reflection")

func findReflection(nums []uint64) []int {
	var out []int
reflection:
	for r := 0; r < len(nums)-1; r++ {
		for i := 0; i <= r; i++ {
			ri := reflectedIndex(i, r)
			if ri > len(nums)-1 {
				continue
			}
			if nums[i] != nums[ri] {
				continue reflection
			}
		}
		out = append(out, r)
	}
	return out
}

func ParsePatterns(lines []string) []Pattern {
	var out []Pattern
	var cur Pattern

	storePattern := func() {
		out = append(out, cur)
		cur = Pattern{}
	}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if line == "" {
			storePattern()
			continue
		}
		y := len(cur.Rows)
		if y == 0 {
			// init cols
			cur.Cols = make([]uint64, len(line))
		}
		r := uint64(0)
		for x, c := range line {
			if c == '#' {
				r |= 1 << x
				cur.Cols[x] |= 1 << y
			}
		}
		cur.Rows = append(cur.Rows, r)
	}
	storePattern()
	return out
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
