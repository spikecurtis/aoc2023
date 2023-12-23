package d20

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Module interface {
	Recv(p Pulse) []Pulse
	Connections() []string
}

type Pulse struct {
	Src, Dst string
	Level    bool
}

const Low = false
const High = true

type broadcaster struct {
	this        string
	connections []string
}

func (b broadcaster) Recv(p Pulse) []Pulse {
	var send []Pulse
	for _, c := range b.connections {
		send = append(send, Pulse{
			Src:   b.this,
			Dst:   c,
			Level: p.Level,
		})
	}
	return send
}

func (b broadcaster) Connections() []string {
	return b.connections
}

type flipFlop struct {
	this        string
	state       bool
	connections []string
}

func (f *flipFlop) Recv(p Pulse) []Pulse {
	var send []Pulse
	if p.Level == High {
		return send
	}
	f.state = !f.state
	for _, c := range f.connections {
		send = append(send, Pulse{
			Src:   f.this,
			Dst:   c,
			Level: f.state,
		})
	}
	return send
}

func (f *flipFlop) Connections() []string {
	return f.connections
}

type conjunction struct {
	this        string
	memory      map[string]bool
	connections []string
}

func (c *conjunction) Recv(p Pulse) []Pulse {
	var send []Pulse
	c.memory[p.Src] = p.Level
	output := Low
	for _, l := range c.memory {
		if l == Low {
			output = High
		}
	}
	if output == Low && (c.this == "vr" || c.this == "xd" || c.this == "pf" || c.this == "ts") {
		log.Println(c.this, "sent low!")
	}
	for _, conn := range c.connections {
		send = append(send, Pulse{
			Src:   c.this,
			Dst:   conn,
			Level: output,
		})
	}
	return send
}

func (c *conjunction) Connections() []string {
	return c.connections
}

func (c *conjunction) MemoryState() uint64 {
	var inputs []string
	for i := range c.memory {
		inputs = append(inputs, i)
	}
	sort.Strings(inputs)
	m := uint64(0)
	for j, name := range inputs {
		if c.memory[name] != High {
			m |= 1 << j
		}
	}
	return m
}

func PressButton(mods map[string]Module, times int) (lows, highs int) {
	for i := 0; i < times; i++ {
		l, h, _ := PressOnce(mods)
		lows += l
		highs += h
	}
	return
}

func LowToRx(mods map[string]Module) int {
	presses := 0
	cycleLen := make(map[string]int)
	nodes := []string{"vr", "xd", "pf", "ts"}
	for {
		if presses%1000000 == 0 {
			log.Println("presses", presses)
		}

		if len(cycleLen) == len(nodes) {
			break
		}
		_, _, conj := PressOnce(mods)
		presses++
		for _, n := range nodes {
			if conj[n] {
				cycleLen[n] = presses
				log.Println("cycle for", n, presses)
			}
		}
		if len(cycleLen) == len(nodes) {
			mult := make([]int, 0, len(nodes))
			for _, l := range cycleLen {
				mult = append(mult, l)
			}
			return LCM(mult[0], mult[1], mult[2:]...)
		}
	}
	return -1
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func PressOnce(mods map[string]Module) (lows, highs int, conj map[string]bool) {
	conj = make(map[string]bool)
	memories := map[string][]string{
		"ts": {"mt", "pb", "xg", "kx", "gk", "lq", "vn", "cd", "rr", "bf", "qm", "jr"},
	}
	var q []Pulse
	q = append(q, Pulse{"button", "broadcaster", Low})
	for len(q) > 0 {
		ts := QueryMemory(mods, memories["ts"])
		if ts == "111111111111" {
			fmt.Println("ts all high")
		}
		p := q[0]
		q = q[1:]
		switch p.Level {
		case Low:
			lows++
		case High:
			highs++
		}
		if p.Level == Low && (p.Src == "ts" || p.Src == "pf" || p.Src == "xd" || p.Src == "vr") {
			conj[p.Src] = true
		}
		m, ok := mods[p.Dst]
		if !ok {
			//log.Println("unknown destination ", p.Dst)
			continue
		}
		pulses := m.Recv(p)
		q = append(q, pulses...)
	}
	return
}

func QueryMemory(mods map[string]Module, flops []string) string {
	b := strings.Builder{}
	for _, f := range flops {
		fl := mods[f].(*flipFlop)
		if fl.state == Low {
			b.WriteString("0")
		} else {
			b.WriteString("1")
		}
	}
	return b.String()
}

func ParseModules(lines []string) map[string]Module {
	modules := make(map[string]Module)
	conjunctions := make(map[string]*conjunction)
	for _, line := range lines {
		if strings.HasPrefix(line, "%") {
			m := ParseFlipFlop(line)
			modules[m.this] = m
			continue
		}
		if strings.HasPrefix(line, "&") {
			m := ParseConjunction(line)
			modules[m.this] = m
			conjunctions[m.this] = m
			continue
		}
		m := ParseBroadcaster(line)
		modules[m.this] = m
	}
	for name, m := range modules {
		for _, c := range m.Connections() {
			if conj, ok := conjunctions[c]; ok {
				conj.memory[name] = Low
			}
		}
	}
	return modules
}

func ParseFlipFlop(s string) *flipFlop {
	s = s[1:] // strip prefix
	parts := strings.Split(s, " ")
	name := parts[0]
	var conns []string
	for i := 2; i < len(parts); i++ {
		c := strings.Trim(parts[i], ",")
		conns = append(conns, c)
	}
	return &flipFlop{
		this:        name,
		state:       Low,
		connections: conns,
	}
}

func ParseConjunction(s string) *conjunction {
	s = s[1:] // strip prefix
	parts := strings.Split(s, " ")
	name := parts[0]
	var conns []string
	for i := 2; i < len(parts); i++ {
		c := strings.Trim(parts[i], ",")
		conns = append(conns, c)
	}
	return &conjunction{
		this:        name,
		memory:      make(map[string]bool),
		connections: conns,
	}
}

func ParseBroadcaster(s string) broadcaster {
	parts := strings.Split(s, " ")
	name := parts[0]
	var conns []string
	for i := 2; i < len(parts); i++ {
		c := strings.Trim(parts[i], ",")
		conns = append(conns, c)
	}
	return broadcaster{
		this:        name,
		connections: conns,
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

func DOT(mods map[string]Module) string {
	dot := strings.Builder{}
	dot.WriteString("digraph {\n")

	for name, m := range mods {
		dot.WriteString("  ")
		dot.WriteString(name)
		switch m.(type) {
		case broadcaster:
			dot.WriteString(" [color = red]")
		case *flipFlop:
			dot.WriteString(" [color = blue]")
		case *conjunction:
			dot.WriteString(" [color = green]")
		}
		dot.WriteString("\n")
		for _, c := range m.Connections() {
			dot.WriteString("  ")
			dot.WriteString(name)
			dot.WriteString(" -> ")
			dot.WriteString(c)
			dot.WriteString("\n")
		}
	}

	dot.WriteString("}\n")
	return dot.String()
}
