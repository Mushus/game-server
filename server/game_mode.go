package server

import "log"

type GameMode interface {
	GetID() string
}

type gameMode struct {
	name         string
	id           string
	maxUser      int
	maxPartyUser int
	matchParty   chan Party
	leaveParty   chan Party
	rooms        map[Room]struct{}
}

func NewGameMode(name string, maxUser int, maxPartyUser int) GameMode {
	return &gameMode{
		name:         name,
		maxUser:      maxUser,
		maxPartyUser: maxPartyUser,
		rooms:        map[Room]struct{}{},
	}
}

func (g *gameMode) GetID() string {
	return g.id
}

func (g *gameMode) Start() {
	for {
		select {
		case party := <-g.matchParty:
			rooms := []Room{}
			for room := range g.rooms {
				room.GetUserCount()
				rooms = append(rooms, room)
			}
			log.Printf("%#v", party)
			// TODO:
		case party := <-g.leaveParty:
			log.Printf("%#v", party)
			// TODO:
		}
	}
}

func (g *gameMode) CreateRoom(party) Room {
	room := &room{}
	go room.start()
	return room
}
