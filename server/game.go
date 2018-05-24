package server

import (
	uuid "github.com/satori/go.uuid"
)

// Game ゲームです
type Game interface {
	JoinUser(user User)
	LeaveUser(user User)
	CreateParty() Party
}

type game struct {
	joinUser       chan User
	leaveUser      chan User
	addParty       chan Party
	removeParty    chan Party
	addGameMode    chan GameMode
	removeGameMode chan GameMode
	users          map[User]struct{}
	parties        map[string]Party
	gameModes      map[string]GameMode
}

// NewGame ゲームを作成します
//
func NewGame(gameModeList []GameMode) Game {
	gameModes := map[string]GameMode{}
	for _, gameMode := range gameModeList {
		id := uuid.NewV4().String()
		gameModes[id] = gameMode
	}

	return &game{
		joinUser:       make(chan User),
		leaveUser:      make(chan User),
		addParty:       make(chan Party),
		removeParty:    make(chan Party),
		addGameMode:    make(chan GameMode),
		removeGameMode: make(chan GameMode),
		users:          map[User]struct{}{},
		parties:        map[string]Party{},
		gameModes:      gameModes,
	}
}

func (g *game) CreateParty() Party {
	id := uuid.NewV4().String()
	party := &party{
		id:    id,
		join:  make(chan User),
		leave: make(chan User),
	}
	g.addParty <- party
	return party
}

func (g *game) start() {
	for _, gm := range g.gameModes {
		gm := gm.(*gameMode)
		go gm.Start()
	}
	for {
		select {
		case user := <-g.joinUser:
			g.users[user] = struct{}{}
		case user := <-g.leaveUser:
			delete(g.users, user)
		case party := <-g.addParty:
			g.parties[party.GetID()] = party
		case party := <-g.removeParty:
			delete(g.parties, party.GetID())
		case gameMode := <-g.addGameMode:
			g.gameModes[gameMode.GetID()] = gameMode
		case gameMode := <-g.removeGameMode:
			delete(g.gameModes, gameMode.GetID())
		}
	}
}

func (g *game) JoinUser(user User) {
	g.joinUser <- user
}

func (g *game) LeaveUser(user User) {
	g.leaveUser <- user
}
