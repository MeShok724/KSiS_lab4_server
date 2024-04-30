package main

import (
	"encoding/json"
	"log"
	"math/rand"
)

const messageReady string = "Ready"
const messageStart string = "Start"
const messageEnemy string = "Enemy"
const messageFish string = "Fish"
const messageUpdate string = "Update"

//var fishOfFirst int = 0

// var fishOfSecond int = 0
var XMax = 1600

var YMax = 800

var ReadyPlayers = []bool{false, false}

type Player struct {
	Name  string
	Enemy *Player
	Score int32
}
type GameMessage struct {
	Command string
	Name    string
	Score   int32
	FishX   int32
	FishY   int32
}

func NewPlayer(name string) *Player {
	player := &Player{Name: name}
	return player
}
func PairPlayers(p1 *Player, p2 *Player) {
	p1.Enemy, p2.Enemy = p2, p1
}
func (p *playerConn) Command(command []byte) {

	var gameMsg GameMessage
	err := json.Unmarshal(command, &gameMsg)
	if err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}
	log.Print("Command: '", gameMsg.Command, gameMsg.Name, gameMsg.Score, "' received by player: ", p.Name)
	switch gameMsg.Command {
	case messageReady:
		p.room.UpdateReady(p.Player)
		log.Print("+ ready")
	case messageFish:
		p.Player.Score++
		x, y := GenerateNewFish(XMax, YMax)
		p.room.SendFishCords(x, y)
	}
}
func GenerateNewFish(maxX int, maxY int) (int32, int32) {
	var x = int32(rand.Intn(maxX + 1)) // Добавляем 1, чтобы включить максимальное значение
	var y = int32(rand.Intn(maxY + 1)) // Добавляем 1, чтобы включить максимальное значение
	return x, y
}
