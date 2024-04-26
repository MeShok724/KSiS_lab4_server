package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var ADDR = ":8080"

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	playerName := "Player"
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}
	// Get or create a room
	var room *room
	if len(freeRooms) > 0 {
		for _, r := range freeRooms {
			room = r
			break
		}
	} else {
		room = NewRoom("")
	}
	// Create Player and Conn
	player := NewPlayer(playerName)
	pConn := NewPlayerConn(ws, player, room)
	// Join Player to room
	room.join <- pConn
	log.Printf("Player: %s has joined to room: %s", pConn.Name, room.name)
}
func main() {
	http.HandleFunc("/ws", wsHandler)

	if err := http.ListenAndServe(ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
