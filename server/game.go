package server

import (
	uuid "github.com/satori/go.uuid"
)

// Game ゲームです
type Game interface {
	JoinUser(user User)
	LeaveUser(user User)
	CreateParty(req CreatePartyRequest)
	JoinParty(req JoinPartyRequest)
}

type game struct {
	joinUser       chan User
	leaveUser      chan User
	createParty    chan CreatePartyRequest
	joinParty      chan JoinPartyRequest
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
		createParty:    make(chan CreatePartyRequest),
		joinParty:      make(chan JoinPartyRequest),
		removeParty:    make(chan Party),
		addGameMode:    make(chan GameMode),
		removeGameMode: make(chan GameMode),
		users:          map[User]struct{}{},
		parties:        map[string]Party{},
		gameModes:      gameModes,
	}
}

func (g *game) CreateParty(req CreatePartyRequest) {
	g.createParty <- req
}

func (g *game) JoinParty(req JoinPartyRequest) {
	g.joinParty <- req
}

func (g *game) JoinUser(user User) {
	g.joinUser <- user
}

func (g *game) LeaveUser(user User) {
	g.leaveUser <- user
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
			user.Send(EventMessage{
				Action: ActionJoinUserToRobby,
				Status: StatusOK,
				Param:  user.ToView(),
			})
		case user := <-g.leaveUser:
			delete(g.users, user)
		case req := <-g.createParty:
			id := uuid.NewV4().String()
			party := &party{
				id:        id,
				owner:     req.User,
				maxUsers:  req.MaxUsers,
				isPrivate: req.IsPrivate,
				join:      make(chan JoinPartyRequest),
				leave:     make(chan User),
				users: map[User]struct{}{
					req.User: struct{}{},
				},
			}
			g.parties[party.GetID()] = party
			go party.Start()
			req.User.Send(EventMessage{
				ID:     req.ID,
				Action: ActionCreateParty,
				Status: StatusOK,
				Param:  party.ToView(),
			})
		case req := <-g.joinParty:
			party, ok := g.parties[req.PartyID]
			if !ok {
				req.User.Send(EventMessage{
					ID:     req.ID,
					Action: ActionJoinParty,
					Status: StatusNG,
					Param:  struct{}{},
				})
				continue
			}
			party.Join(req)
		case party := <-g.removeParty:
			delete(g.parties, party.GetID())
		case gameMode := <-g.addGameMode:
			g.gameModes[gameMode.GetID()] = gameMode
		case gameMode := <-g.removeGameMode:
			delete(g.gameModes, gameMode.GetID())
		}
	}
}
