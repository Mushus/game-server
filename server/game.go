package server

import (
	uuid "github.com/satori/go.uuid"
)

// Game ゲームです
type Game interface {
	CreateUserRequest(userName string, event chan EventMessage) UserView
	LeaveUserFromGameRequest(userID string)
	CreatePartyRequest(userID string, isPrivate bool, maxUsers int) (*PartyView, bool)
	JoinPartyRequest(userID string, partyID string) (*PartyView, bool)
	LeaveUserFromPartyRequest(userID string) bool
}

type game struct {
	users              map[string]*user
	parties            map[string]*party
	gameModes          map[string]*gameMode
	createUser         chan createUserRequest
	leaveUserFromGame  chan leaveUserFromGameRequest
	createParty        chan createPartyRequest
	joinParty          chan joinPartyRequest
	leaveUserFromParty chan leaveUserFromPartyRequest
}

// NewGame ゲームを作成します
//
func NewGame(gameModeList []GameMode) Game {
	gameModes := map[string]*gameMode{}
	for _, igm := range gameModeList {
		gm, ok := igm.(*gameMode)
		if !ok {
			continue
		}
		gameModes[gm.id] = gm
	}

	return &game{
		users:             map[string]*user{},
		parties:           map[string]*party{},
		gameModes:         gameModes,
		createUser:        make(chan createUserRequest),
		leaveUserFromGame: make(chan leaveUserFromGameRequest),
		createParty:       make(chan createPartyRequest),
	}
}

func (g *game) start() {
	for {
		select {
		case req := <-g.createUser:
			user := g.createUserAction(req.userName, req.event)
			req.resp <- createUserResponse{
				user: user,
			}
		case req := <-g.leaveUserFromGame:
			g.leaveUserFromGameAction(req.userID)
			req.resp <- leaveUserFromGameResponse{}
		case req := <-g.createParty:
			party, status := g.createPartyAction(req.userID, req.isPrivate, req.maxUsers)
			req.resp <- createPartyResponse{
				party:  party,
				status: status,
			}
		case req := <-g.joinParty:
			party, status := g.joinPartyAction(req.userID, req.partyID)
			req.resp <- joinPartyResponse{
				party:  party,
				status: status,
			}
		case req := <-g.leaveUserFromParty:
			status := g.leaveUserFromPartyAction(req.userID)
			req.resp <- leaveUserFromPartyResponse{
				status: status,
			}
		}
	}
}

// ===========================================================================
// サーバーの動作

func (g *game) createUserAction(userName string, event chan EventMessage) UserView {
	id := uuid.NewV4().String()
	user := &user{
		id:    id,
		name:  userName,
		event: event,
	}
	g.users[id] = user
	return user.ToView()
}

func (g *game) leaveUserFromGameAction(userID string) {
	delete(g.users, userID)
}

func (g *game) createPartyAction(userID string, isPrivate bool, maxUsers int) (*PartyView, bool) {
	owner, ok := g.users[userID]
	if !ok {
		return nil, false
	}
	id := uuid.NewV4().String()
	party := &party{
		id:        id,
		owner:     owner,
		users:     map[string]*user{userID: owner},
		maxUsers:  maxUsers,
		isPrivate: isPrivate,
	}
	owner.party = party
	pv := party.ToView()
	return &pv, true
}

func (g *game) joinPartyAction(userID string, partyID string) (*PartyView, bool) {
	rookie, ok := g.users[userID]
	if !ok {
		return nil, false
	}

	targetParty, ok := g.parties[partyID]
	if !ok {
		return nil, false
	}

	targetParty.users[userID] = rookie
	rookie.party = targetParty
	// TODO: ModifyParty
	pv := targetParty.ToView()
	return &pv, true
}

func (g *game) leaveUserFromPartyAction(userID string) bool {
	leaver, ok := g.users[userID]
	if !ok {
		return false
	}

	party := leaver.party
	delete(party.users, userID)
	if len(party.users) == 0 {
		delete(g.parties, party.id)
	}
	if party.owner.id == userID {
		party.owner = nil
		for _, member := range party.users {
			party.owner = member
			break
		}
	}
	return true
}

// ===========================================================================
// サーバーに対するリクエスト

func (g *game) CreateUserRequest(userName string, event chan EventMessage) UserView {
	respCh := make(chan createUserResponse)
	g.createUser <- createUserRequest{
		resp:     respCh,
		userName: userName,
		event:    event,
	}
	resp := <-respCh
	return resp.user
}

func (g *game) LeaveUserFromGameRequest(userID string) {
	respCh := make(chan leaveUserFromGameResponse)
	g.leaveUserFromGame <- leaveUserFromGameRequest{
		resp:   respCh,
		userID: userID,
	}
	<-respCh
}

func (g *game) CreatePartyRequest(userID string, isPrivate bool, maxUsers int) (*PartyView, bool) {
	respCh := make(chan createPartyResponse)
	g.createParty <- createPartyRequest{
		resp:      respCh,
		userID:    userID,
		isPrivate: isPrivate,
		maxUsers:  maxUsers,
	}
	resp := <-respCh
	return resp.party, resp.status
}

func (g *game) JoinPartyRequest(userID string, partyID string) (*PartyView, bool) {
	respCh := make(chan joinPartyResponse)
	g.joinParty <- joinPartyRequest{
		resp:    respCh,
		userID:  userID,
		partyID: partyID,
	}
	resp := <-respCh
	return resp.party, resp.status
}

func (g *game) LeaveUserFromPartyRequest(userID string) bool {
	respCh := make(chan leaveUserFromPartyResponse)
	g.leaveUserFromParty <- leaveUserFromPartyRequest{
		resp:   respCh,
		userID: userID,
	}
	resp := <-respCh
	return resp.status
}
