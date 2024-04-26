package main

import (
	"encoding/json"
	"log"
)

const messageReady string = "Ready"
const messageStart string = "Start"
const messageEnemy string = "Enemy"

var ReadyPlayers = []bool{false, false}

type Player struct {
	Name  string
	Enemy *Player
	Score int
}
type GameMessage struct {
	Command string
	Name    string
	Score   int32
}

func NewPlayer(name string) *Player {
	player := &Player{Name: name}
	return player
}
func PairPlayers(p1 *Player, p2 *Player) {
	p1.Enemy, p2.Enemy = p2, p1
}
func (p *playerConn) Command(command []byte) {
	log.Print("Command: '", command, "' received by player: ", p.Name)
	var gameMsg GameMessage
	err := json.Unmarshal(command, &gameMsg)
	if err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}
	switch gameMsg.Command {
	case messageReady:
		p.room.UpdateReady(p.Player)
	}
}
func (p *Player) GetState() string {
	return "Game state for Player: " + p.Name
}
func (p *Player) GiveUp() {
	log.Print("Player gave up: ", p.Name)
}
