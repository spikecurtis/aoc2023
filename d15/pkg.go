package d15

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Box struct {
	Lenses []Lens
}

type Lens struct {
	Label       string
	FocalLength int
}

type step struct {
	label       string
	op          string
	focalLength int
}

func parseStep(s string) step {
	var out step
	var err error
	i := strings.Index(s, "=")
	if i > 0 {
		out.op = "="
		out.label = s[:i]
		out.focalLength, err = strconv.Atoi(s[i+1:])
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
	i = strings.Index(s, "-")
	out.op = "-"
	out.label = s[:i]
	return out
}

type HASHMAP [256]Box

func (h *HASHMAP) Do(s string) {
	st := parseStep(s)
	b := HASH(st.label)
	box := h[b]
	switch st.op {
	case "-":
		for i, l := range box.Lenses {
			if l.Label == st.label {
				var remaining []Lens
				if i+1 < len(box.Lenses) {
					remaining = box.Lenses[i+1:]
				}
				box.Lenses = append(box.Lenses[:i], remaining...)
				break
			}
		}
	case "=":
		present := false
		for i, l := range box.Lenses {
			if l.Label == st.label {
				box.Lenses[i].FocalLength = st.focalLength
				present = true
				break
			}
		}
		if !present {
			box.Lenses = append(box.Lenses, Lens{st.label, st.focalLength})
		}
	}
	h[b] = box
}

func (h *HASHMAP) FocusingPower() int {
	total := 0
	for b, box := range h {
		for slot, l := range box.Lenses {
			total += (b + 1) * (slot + 1) * l.FocalLength
		}
	}
	return total
}

func HASH(s string) int {
	v := 0
	for _, c := range s {
		v += int(c)
		v *= 17
		v %= 256
	}
	return v
}

func GetInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	all, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(all), ",")
}
