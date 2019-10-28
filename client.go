package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type client struct {
	conn   *websocket.Conn
	readCh chan []byte
}

func (c *client) Read() {
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.readCh <- p
	}
}

func (c *client) Join(ses *session) {
	ses.clients = append(ses.clients, c)
	ses.Broadcast(1, []byte("Client joined session"))
}

func listen(ch chan []byte, d time.Duration, f func(m []byte)) {
	for {
		select {
		case m := <-ch:
			f(m)
		case <-time.After(d):
			return
		}
	}
}
