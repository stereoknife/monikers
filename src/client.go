package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type client struct {
	conn   *websocket.Conn
	readCh chan []byte
	name   string
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

func (c *client) ListenFor(duration time.Duration, f func([]byte) bool) {
	for {
		select {
		case m := <-c.readCh:
			if f(m) {
				return
			}
		case <-time.After(duration):
			fmt.Println("stopped listening")
			return
		}
	}
}
