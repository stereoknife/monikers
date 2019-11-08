package main

type team struct {
	members []int
	_player int
	score   int
}

func (t team) NextPlayer() (int, int) {
	rp, ri := t.members[t._player], t._player
	t._player = (t._player + 1) % len(t.members)
	return rp, ri
}
