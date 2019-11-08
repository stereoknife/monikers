package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var sesh = session{
	_deck: map[uint32]string{1: "hello", 2: "there", 3: "how's", 4: "it", 5: "going"},
}

func socketSetup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := client{
		ws,
		make(chan []byte),
		r.FormValue("name"),
	}

	c.Join(&sesh)
	go c.Read()
	if sesh.status == 0 {
		sesh.Start()
	}
}

func setRoutes(r *mux.Router) {
	r.HandleFunc("/ws/join", socketSetup).Methods("GET")
}

func main() {
	fmt.Println("Monikers Server v0.0.0")

	router := mux.NewRouter()
	setRoutes(router)

	log.Fatal(http.ListenAndServe(":1337", router))
}
