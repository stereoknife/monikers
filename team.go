package main

type team struct {
	members []int
	last    int
	score   int
}

func (t team) Next() int {
	r := t.members[t.last]
	t.last = (t.last + 1) % len(t.members)
	return r
}
