package src

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

func (ses *session) Broadcast(messageType int, msg []byte) {
	for _, c := range ses.clients {
		if err := c.conn.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
		}
	}
}

func (ses *session) ListenFor(c *client, duration time.Duration, f func([]byte)) {
	for len(ses.deck) > 0 {
		select {
		case m := <-c.readCh:
			f(m)
		case <-time.After(duration):
			fmt.Println("stopped listening")
			return
		}
	}
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

func (ses *session) NextPlayer() (*client, int) {
	t, ti := ses.NextTeam()
	p, _ := t.NextPlayer()
	return ses.clients[p], ti
}

func (ses *session) NextTeam() (*team, int) {
	rt, ri := &ses.teams[ses._team], ses._team
	ses._team = (ses._team + 1) % 2
	return rt, ri
}

// Game below

func (ses *session) Phase() {
	// Init array for correct cards
	var guess [2][]uint32

	// Phase loop: loop until remaining cards reaches 0
	for len(ses.deck) > 0 {
		// Get next player and team
		cl, team := ses.NextPlayer()
		fmt.Printf("%v is up\n", cl.name)

		// Listen for a minute
		cl.ListenFor(time.Minute, func(m []byte) bool {
			i := binary.BigEndian.Uint32(m)
			if i == 0 {
				ses.Broadcast(1, []byte("skipped"))
			} else {
				// Check that card hasn't been guessed
				if _, ok := ses.deck[i]; ok {
					// Remove from remaining cards and move it into correct guesses for the team
					guess[team] = append(guess[team], i)
					delete(ses.deck, i)
					fmt.Printf("Guessed %v \n", i)
				}
			}
			return len(ses.deck) <= 0
		})
	}

	// Add scores
	for i, team := range ses.teams {
		team.score += len(guess[i])
	}
}

func (ses *session) Start() {
	// Setup phase
	// Game phase x3
	// Cleanup & DC
	if ses.status > 0 {
		return
	}
	ses.status++
	<-ses.clients[0].readCh
	ses.MakeTeams()
	ses.deck = ses._deck
	ses.Phase()
}
