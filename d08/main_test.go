package main

import "testing"

func TestFollow(t *testing.T) {
	aaa := &node{name: "AAA"}
	bbb := &node{name: "BBB"}
	zzz := &node{name: "ZZZ"}
	aaa.next[left] = bbb
	aaa.next[right] = bbb
	bbb.next[left] = aaa
	bbb.next[right] = zzz
	zzz.next[left] = zzz
	zzz.next[right] = zzz
	g := guide{
		instructions: []direction{left, left, right},
		network:      aaa,
	}
	steps := g.followToZZZ()
	if steps != 6 {
		t.Fatal()
	}
}
