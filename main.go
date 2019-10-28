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
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	c := client{ws, make(chan []byte)}
	c.Join(&sesh)
	go c.Read()
	sesh.Start()
	fmt.Println("end of loop")
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
