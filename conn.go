package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type playerConn struct {
	ws *websocket.Conn
	*Player
	room *room
}

// Receive msg from ws in goroutine
func (pc *playerConn) receiver() {
	for {
		_, command, err := pc.ws.ReadMessage()
		if err != nil {
			break
		}
		// execute a command
		pc.Command(command)
		// update all conn
		pc.room.updateAll <- true
	}
	pc.room.leave <- pc
	pc.ws.Close()
}

//	func (pc *playerConn) sendState() {
//		go func() {
//			msg := pc.GetState()
//			err := pc.ws.WriteMessage(websocket.TextMessage, []byte(msg))
//			if err != nil {
//				pc.room.leave <- pc
//				pc.ws.Close()
//			}
//		}()
//	}
func (pc *playerConn) SendMessage(message GameMessage) {
	go func() {
		log.Println("Sending message to player:", pc.Player, ":", message)
		msg, _ := json.Marshal(message)
		err := pc.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			pc.room.leave <- pc
			pc.ws.Close()
		}
	}()
}
func NewPlayerConn(ws *websocket.Conn, player *Player, room *room) *playerConn {
	pc := &playerConn{ws, player, room}
	go pc.receiver()
	return pc
}
