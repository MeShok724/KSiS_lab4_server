package main

import (
	"log"
	"strconv"
	"sync"
)

var allRooms = make(map[string]*room)
var freeRooms = make(map[string]*room)
var roomsCount int
var wsMutex sync.Mutex

type room struct {
	name              string
	scoreFirstPlayer  int
	scoreSecondPlayer int

	// Registered connections.
	playerConns map[*playerConn]bool

	// Update state for all conn.
	updateAll chan bool

	// Register requests from the connections.
	join chan *playerConn

	// Unregister requests from connections.
	leave        chan *playerConn
	readyPlayers int
}

// Run the room in goroutine
func (r *room) run() {

	for {
		select {
		case c := <-r.join:
			r.playerConns[c] = true
			//r.updateAllPlayers()

			// if room is full - delete from freeRooms
			if len(r.playerConns) == 2 {
				delete(freeRooms, r.name)
				// pair players
				var p []*Player
				for k, _ := range r.playerConns {
					p = append(p, k.Player)
				}
				PairPlayers(p[0], p[1])
				for k, _ := range r.playerConns {
					k.SendMessage(GameMessage{Command: messageEnemy, Name: k.Enemy.Name, Score: k.Enemy.Score})
				}
			}

		case c := <-r.leave:
			//c.GiveUp()
			//r.updateAllPlayers()
			delete(r.playerConns, c)
			if len(r.playerConns) == 0 {
				goto Exit
			}
		case <-r.updateAll:
			//r.updateAllPlayers()
		}
	}

Exit:

	// delete room
	delete(allRooms, r.name)
	delete(freeRooms, r.name)
	roomsCount -= 1
	log.Print("room closed:", r.name)
}

//	func (r *room) updateAllPlayers() {
//		for c := range r.playerConns {
//			message := GameMessage{Command: messageUpdate, Score: c.Score}
//			c.SendMessage(message)
//		}
//	}

func NewRoom(name string) *room {
	if name == "" {
		name = "1"
		for IsExist(name) {
			name = name + strconv.FormatInt(int64(len(name)), 10)
		}
	}

	room := &room{
		name:        name,
		playerConns: make(map[*playerConn]bool),
		updateAll:   make(chan bool),
		join:        make(chan *playerConn),
		leave:       make(chan *playerConn),
	}

	allRooms[name] = room
	freeRooms[name] = room

	// run room
	go room.run()

	roomsCount += 1

	return room
}

func (r *room) UpdateReady(p *Player) {
	r.readyPlayers++
	log.Print("Ready players:", r.readyPlayers)
	if r.readyPlayers == 2 {
		log.Print("Both ready")
		for i, _ := range r.playerConns {
			wsMutex.Lock()
			i.SendMessage(GameMessage{Command: messageStart})
			wsMutex.Unlock()
		}
		x, y := GenerateNewFish(XMax, YMax)
		r.SendFishCords(x, y)
	}
}
func IsExist(nameToCheck string) bool {
	for _, curr := range allRooms {
		if curr.name == nameToCheck {
			return true
		}
	}
	return false
}
