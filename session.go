package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type session struct {
	clients []*client
	teams   [2]team
	_team   int
	status  int
	_deck   map[uint32]string
	deck    map[uint32]string
}

func (ses *session) Listen(c *client) ([]uint32, bool) {
	var d []uint32
l:
	for len(ses.deck) > 0 {
		select {
		case m := <-c.readCh:
			i := binary.BigEndian.Uint32(m)
			if i == 0 {
				ses.Broadcast(1, []byte("skipped"))
				continue
			}
			_, ok := ses.deck[i]
			if ok {
				d = append(d, i)
				delete(ses.deck, i)
				fmt.Printf("deleted %v \n", i)
			}
			fmt.Println(d)
			log.Println(ses.deck)
		case <-time.After(time.Minute):
			fmt.Println("stopped listening")
			break l
		}
	}
	return d, len(d) > 0
}

func (ses *session) MakeTeams() {
	l := len(ses.clients)
	h := int(math.Floor(float64(l) * 0.5))
	p := rand.Perm(l)
	ses.teams[0] = team{
		members: p[:h],
	}
	ses.teams[1] = team{
		members: p[h:],
	}
}

func (ses *session) Broadcast(messageType int, msg []byte) {
	for _, c := range ses.clients {
		if err := c.conn.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
		}
	}
}

func (ses *session) Phase() {
	var guess [2][]uint32
	for len(ses.deck) > 0 {
		cl, team := ses.NextPlayer()
		app, ok := ses.Listen(cl)
		if ok {
			guess[team] = append(guess[team], app...)
		}
	}
}

func (ses *session) Start() {
	// Setup phase
	// Game phase x3
	// Cleanup & DC
	ses.MakeTeams()
	ses.deck = ses._deck
}

func (ses *session) NextPlayer() (*client, int) {
	t := ses.NextTeamIndex()
	return ses.clients[ses.teams[t].Next()], t
}

func (ses *session) NextTeamIndex() int {
	r := ses._team
	ses._team = (ses._team + 1) % 2
	return r
}
