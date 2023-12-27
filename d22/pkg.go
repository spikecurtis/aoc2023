package d22

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	Coord [3]int
}

func (p Point) Drop(n int) Point {
	return Point{Coord: [3]int{p.Coord[0], p.Coord[1], p.Coord[2] - n}}
}

type Brick struct {
	Start, End Point
}

func (b Brick) Bottom() int {
	return min(b.Start.Coord[2], b.End.Coord[2])
}

func (b Brick) Min(d int) int {
	return min(b.Start.Coord[d], b.End.Coord[d])
}

func (b Brick) Max(d int) int {
	return max(b.Start.Coord[d], b.End.Coord[d])
}

func (b Brick) Overlaps(o Brick) bool {
	for d := range b.Start.Coord {
		if b.Min(d) > o.Max(d) || b.Max(d) < o.Min(d) {
			return false
		}
	}
	return true
}

func (b Brick) Drop(n int) Brick {
	return Brick{
		Start: b.Start.Drop(n),
		End:   b.End.Drop(n),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func ParseBricks(lines []string) []Brick {
	var out []Brick
	for _, line := range lines {
		out = append(out, ParseBrick(line))
	}
	return out
}

func DropBricks(b []Brick) []Brick {
	bricks := make([]Brick, len(b))
	copy(bricks, b)
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].Bottom() < bricks[j].Bottom()
	})
	for i := range bricks {
		this := bricks[i]
	drop:
		for {
			if this.Bottom() == 1 {
				break
			}
			next := this.Drop(1)
			for j := 0; j < i; j++ {
				that := bricks[j]
				if next.Overlaps(that) {
					break drop
				}
			}
			this = next
		}
		bricks[i] = this
	}
	return bricks
}

func SupportedBy(bricks []Brick) map[int][]int {
	out := make(map[int][]int)
	for i, this := range bricks {
		thisD := this.Drop(1)
		for j := 0; j < i; j++ {
			if thisD.Overlaps(bricks[j]) {
				out[i] = append(out[i], j)
			}
		}
	}
	return out
}

func CountDisintegrate(supports map[int][]int, n int) int {
	cant := make(map[int]bool)
	for i := 0; i < n; i++ {
		s := supports[i]
		if len(s) == 1 {
			// can't disintegrate a brick if it's the only thing supporting
			// another
			cant[s[0]] = true
		}
	}
	return n - len(cant)
}

func CountTotalChainDisintegrate(supports map[int][]int, n int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += CountChanDisintegrate(supports, n, i)
	}
	return total
}

func CountChanDisintegrate(supports map[int][]int, n, i int) int {
	fall := make([]bool, n)
	fall[i] = true
	for {
		fell := 0
		for j := 0; j < n; j++ {
			if fall[j] { // already fallen
				continue
			}
			s := supports[j]
			if len(s) > 0 { // only bricks supported by bricks can fall
				stillSupported := false
				for _, k := range s {
					if !fall[k] {
						stillSupported = true
					}
				}
				if !stillSupported {
					fell++
					fall[j] = true
				}
			}
		}
		if fell == 0 {
			break
		}
	}
	falls := 0
	for j := 0; j < n; j++ {
		if j == i {
			continue
		}
		if fall[j] {
			falls++
		}
	}
	return falls
}

func ParseBrick(line string) Brick {
	parts := strings.Split(line, "~")
	return Brick{
		Start: ParsePoint(parts[0]),
		End:   ParsePoint(parts[1]),
	}
}

func ParsePoint(s string) Point {
	parts := strings.Split(s, ",")
	var xyz [3]int
	for n, p := range parts {
		i, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		xyz[n] = i
	}
	return Point{Coord: xyz}
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
