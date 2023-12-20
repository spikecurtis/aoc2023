package d19

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Workflows map[string][]Rule

func (w Workflows) Process(p Part) string {
	target := "in"
	for {
		if target == "A" || target == "R" {
			return target
		}
		rules, ok := w[target]
		if !ok {
			log.Fatal("could not find workflow")
		}
		match := false
		for _, r := range rules {
			if r.Matches(p) {
				target = r.Result
				match = true
				break
			}
		}
		if !match {
			log.Fatal("no match")
		}
	}
}

func (w Workflows) Combinations() int {
	targets := []RangeResult{
		{
			PartRange: PartRange{
				X: NumRange{1, 4000},
				M: NumRange{1, 4000},
				A: NumRange{1, 4000},
				S: NumRange{1, 4000},
			},
			Result: "in",
		},
	}
	combos := 0
	for len(targets) > 0 {
		target := targets[len(targets)-1]
		targets = targets[:len(targets)-1]
		rules := w[target.Result]
		unresolved := []RangeResult{target}
		resolved := []RangeResult{}
		for _, r := range rules {
			newUnresolved := []RangeResult{}
			for _, t := range unresolved {
				nextTargets := r.MatchPartRange(t.PartRange)
				for _, nt := range nextTargets {
					if nt.Result == "next" {
						newUnresolved = append(newUnresolved, nt)
						continue
					}
					resolved = append(resolved, nt)
				}
			}
			unresolved = newUnresolved
		}
		if len(unresolved) != 0 {
			log.Fatal("didn't resolve all rules")
		}
		for _, rr := range resolved {
			if rr.Result == "R" {
				continue
			}
			if rr.Result == "A" {
				combos += rr.PartRange.Combos()
				continue
			}
			targets = append(targets, rr)
		}
	}
	return combos
}

type Rule struct {
	Var          func(p Part) int
	VarRange     func(pr PartRange) NumRange
	ReplaceRange func(pr PartRange, nr NumRange) PartRange
	Cmp          func(i int) bool
	MatchRange   func(nr NumRange) []NumRange // first is match (if any)
	Result       string
}

func (r Rule) Matches(p Part) bool {
	return r.Cmp(r.Var(p))
}

func (r Rule) MatchPartRange(pr PartRange) []RangeResult {
	nr := r.VarRange(pr)
	rngs := r.MatchRange(nr)
	var out []RangeResult
	for i, rng := range rngs {
		if rng.Empty() {
			continue
		}
		if i == 0 {
			// match!
			out = append(out, RangeResult{
				PartRange: r.ReplaceRange(pr, rng),
				Result:    r.Result,
			})
			continue
		}
		out = append(out, RangeResult{
			PartRange: r.ReplaceRange(pr, rng),
			Result:    "next",
		})
	}
	return out
}

type Part struct {
	X, M, A, S int
}

type PartRange struct {
	X, M, A, S NumRange
}

func (pr PartRange) Combos() int {
	c := 1
	c *= pr.X.Size()
	c *= pr.M.Size()
	c *= pr.A.Size()
	c *= pr.S.Size()
	return c
}

type NumRange struct {
	Min, Max int
}

func (nr NumRange) Size() int {
	if nr.Empty() {
		return 0
	}
	return nr.Max - nr.Min + 1
}

func (nr NumRange) Empty() bool {
	if nr.Min == 0 {
		return true
	}
	if nr.Max == 0 {
		return true
	}
	if nr.Min > nr.Max {
		return true
	}
	return false
}

type RangeResult struct {
	PartRange
	Result string
}

func ParsePuzzle(lines []string) (Workflows, []Part) {
	workflows := make(Workflows)
	var parts []Part
	i := 0
	for {
		line := lines[i]
		i++
		if line == "" {
			break
		}
		name, rules := ParseWorkflow(line)
		workflows[name] = rules
	}
	for {
		line := lines[i]
		i++
		parts = append(parts, ParsePart(line))
		if i == len(lines) {
			break
		}
	}
	return workflows, parts
}

func ParseWorkflow(line string) (string, []Rule) {
	parts := strings.Split(line, "{")
	name := parts[0]
	var rules []Rule
	ruleStrings := strings.Split(strings.Trim(parts[1], "}"), ",")
	for _, rs := range ruleStrings {
		rules = append(rules, ParseRule(rs))
	}
	return name, rules
}

func ParseRule(s string) Rule {
	parts := strings.Split(s, ":")
	if len(parts) == 1 {
		return Rule{
			Var:          func(Part) int { return 0 },
			VarRange:     func(PartRange) NumRange { return NumRange{1, 4000} }, // trivially non-empty
			ReplaceRange: func(pr PartRange, _ NumRange) PartRange { return pr },
			Cmp:          func(int) bool { return true },
			MatchRange:   func(nr NumRange) []NumRange { return []NumRange{nr} },
			Result:       s,
		}
	}
	r := Rule{Result: parts[1]}
	switch parts[0][0] {
	case 'x':
		r.Var = func(p Part) int { return p.X }
		r.VarRange = func(pr PartRange) NumRange { return pr.X }
		r.ReplaceRange = func(pr PartRange, nr NumRange) PartRange {
			pr.X = nr
			return pr
		}
	case 'm':
		r.Var = func(p Part) int { return p.M }
		r.VarRange = func(pr PartRange) NumRange { return pr.M }
		r.ReplaceRange = func(pr PartRange, nr NumRange) PartRange {
			pr.M = nr
			return pr
		}
	case 'a':
		r.Var = func(p Part) int { return p.A }
		r.VarRange = func(pr PartRange) NumRange { return pr.A }
		r.ReplaceRange = func(pr PartRange, nr NumRange) PartRange {
			pr.A = nr
			return pr
		}
	case 's':
		r.Var = func(p Part) int { return p.S }
		r.VarRange = func(pr PartRange) NumRange { return pr.S }
		r.ReplaceRange = func(pr PartRange, nr NumRange) PartRange {
			pr.S = nr
			return pr
		}
	default:
		log.Fatal("unknown var")
	}
	j, err := strconv.Atoi(parts[0][2:])
	if err != nil {
		log.Fatal(err)
	}
	switch parts[0][1] {
	case '<':
		r.Cmp = func(i int) bool { return i < j }
		r.MatchRange = func(nr NumRange) []NumRange {
			if nr.Min < j {
				return []NumRange{
					{
						Min: nr.Min,
						Max: min(nr.Max, j-1),
					},
					{
						Min: j,
						Max: nr.Max,
					},
				}
			}
			return []NumRange{}
		}
	case '>':
		r.Cmp = func(i int) bool { return i > j }
		r.MatchRange = func(nr NumRange) []NumRange {
			if nr.Max > j {
				return []NumRange{
					{
						Min: max(nr.Min, j+1),
						Max: nr.Max,
					},
					{
						Min: nr.Min,
						Max: j,
					},
				}
			}
			return []NumRange{}
		}
	default:
		log.Fatal("unknown cmp")
	}
	return r
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

func ParsePart(line string) Part {
	line = strings.Trim(line, "{}")
	parts := strings.Split(line, ",")
	x, err := strconv.Atoi(parts[0][2:])
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(parts[1][2:])
	if err != nil {
		log.Fatal(err)
	}
	a, err := strconv.Atoi(parts[2][2:])
	if err != nil {
		log.Fatal(err)
	}
	s, err := strconv.Atoi(parts[3][2:])
	if err != nil {
		log.Fatal(err)
	}
	return Part{
		X: x,
		M: m,
		A: a,
		S: s,
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
