package server

import (
	uuid "github.com/satori/go.uuid"
)

// Game ゲームです
type Game interface {
	CreateUserRequest(userName string, event chan interface{}) UserView
	LeaveUserFromGameRequest(userID string)
	CreatePartyRequest(userID string, isPrivate bool, maxUsers int) (*PartyView, bool)
	JoinPartyRequest(userID string, partyID string) (*PartyView, bool)
	LeaveUserFromPartyRequest(userID string) bool
	RequestP2PRequest(userID string, targetID string, offer string) bool
	ResponseP2PRequest(userID string, targetID string, answer string) bool
}

type game struct {
	users                map[string]*user
	parties              map[string]*party
	gameModes            map[string]*gameMode
	createUserCh         chan createUserRequest
	leaveUserFromGameCh  chan leaveUserFromGameRequest
	createPartyCh        chan createPartyRequest
	joinPartyCh          chan joinPartyRequest
	leaveUserFromPartyCh chan leaveUserFromPartyRequest
	requestP2PCh         chan requestP2PRequest
	responseP2PCh        chan responseP2PRequest
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
		users:                map[string]*user{},
		parties:              map[string]*party{},
		gameModes:            gameModes,
		createUserCh:         make(chan createUserRequest),
		leaveUserFromGameCh:  make(chan leaveUserFromGameRequest),
		createPartyCh:        make(chan createPartyRequest),
		joinPartyCh:          make(chan joinPartyRequest),
		leaveUserFromPartyCh: make(chan leaveUserFromPartyRequest),
		requestP2PCh:         make(chan requestP2PRequest),
		responseP2PCh:        make(chan responseP2PRequest),
	}
}

func (g *game) start() {
	for {
		select {
		case req := <-g.createUserCh:
			user := g.createUserAction(req.userName, req.event)
			req.resp <- createUserResponse{
				user: user,
			}
		case req := <-g.leaveUserFromGameCh:
			g.leaveUserFromGameAction(req.userID)
			req.resp <- leaveUserFromGameResponse{}
		case req := <-g.createPartyCh:
			party, status := g.createPartyAction(req.userID, req.isPrivate, req.maxUsers)
			req.resp <- createPartyResponse{
				party:  party,
				status: status,
			}
		case req := <-g.joinPartyCh:
			party, status := g.joinPartyAction(req.userID, req.partyID)
			req.resp <- joinPartyResponse{
				party:  party,
				status: status,
			}
		case req := <-g.leaveUserFromPartyCh:
			status := g.leaveUserFromPartyAction(req.userID)
			req.resp <- leaveUserFromPartyResponse{
				status: status,
			}
		case req := <-g.requestP2PCh:
			status := g.requestP2PAction(req.userID, req.targetID, req.offer)
			req.resp <- requestP2PResponse{
				status: status,
			}
		case req := <-g.responseP2PCh:
			status := g.responseP2PAction(req.userID, req.targetID, req.answer)
			req.resp <- responseP2PResponse{
				status: status,
			}
		}
	}
}

// ===========================================================================
// サーバーの動作

func (g *game) createUserAction(userName string, event chan interface{}) UserView {
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
	user, ok := g.users[userID]
	if ok {
		g.leaveUserFromParty(user)
		delete(g.users, userID)
	}
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
		users:     map[string]*user{},
		maxUsers:  maxUsers,
		isPrivate: isPrivate,
	}
	g.parties[id] = party
	g.joinParty(owner, party)
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

	g.joinParty(rookie, targetParty)
	// TODO: ModifyParty
	pv := targetParty.ToView()
	return &pv, true
}

func (g *game) leaveUserFromPartyAction(userID string) bool {
	leaver, ok := g.users[userID]
	if !ok {
		return false
	}
	g.leaveUserFromParty(leaver)
	return true
}

func (g *game) requestP2PAction(userID string, targetID string, offer string) bool {
	_, ok := g.users[userID]
	if !ok {
		return false
	}

	target, ok := g.users[targetID]
	if !ok {
		return false
	}

	target.Send(RequestP2PEvent{
		Offer:  offer,
		UserID: userID,
	})
	return true
}

func (g *game) responseP2PAction(userID string, targetID string, answer string) bool {
	_, ok := g.users[userID]
	if !ok {
		return false
	}

	target, ok := g.users[targetID]
	if !ok {
		return false
	}

	target.Send(ResponseP2PEvent{
		Answer: answer,
		UserID: userID,
	})
	return true
}

// ===========================================================================
// 共通操作

func (g *game) joinParty(rookie *user, party *party) {
	// すでにパーティにいる場合は退席
	if rookie.party != nil {
		g.leaveUserFromParty(rookie)
	}

	users := map[string]*user{}
	for k, v := range party.users {
		users[k] = v
	}
	// パーティに参加
	party.users[rookie.id] = rookie
	rookie.party = party

	// 変更を通知
	pv := party.ToView()
	for _, member := range users {
		m := member
		m.Send(ModifyPartyEvent{
			Party: pv,
		})
	}
}

func (g *game) leaveUserFromParty(user *user) {
	party := user.party

	// paryから退出
	delete(party.users, user.id)
	// 人がいないパーティは破棄
	if len(party.users) == 0 {
		delete(g.parties, party.id)
	}
	if party.owner.id == user.id {
		party.owner = nil
		for _, member := range party.users {
			party.owner = member
			break
		}
	}
	// パーティ変更を通知
	pv := party.ToView()
	for _, member := range party.users {
		m := member
		go m.Send(ModifyPartyEvent{
			Party: pv,
		})
	}
}

// ===========================================================================
// サーバーに対するリクエスト

func (g *game) CreateUserRequest(userName string, event chan interface{}) UserView {
	respCh := make(chan createUserResponse)
	g.createUserCh <- createUserRequest{
		resp:     respCh,
		userName: userName,
		event:    event,
	}
	resp := <-respCh
	return resp.user
}

func (g *game) LeaveUserFromGameRequest(userID string) {
	respCh := make(chan leaveUserFromGameResponse)
	g.leaveUserFromGameCh <- leaveUserFromGameRequest{
		resp:   respCh,
		userID: userID,
	}
	<-respCh
}

func (g *game) CreatePartyRequest(userID string, isPrivate bool, maxUsers int) (*PartyView, bool) {
	respCh := make(chan createPartyResponse)
	g.createPartyCh <- createPartyRequest{
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
	g.joinPartyCh <- joinPartyRequest{
		resp:    respCh,
		userID:  userID,
		partyID: partyID,
	}
	resp := <-respCh
	return resp.party, resp.status
}

func (g *game) LeaveUserFromPartyRequest(userID string) bool {
	respCh := make(chan leaveUserFromPartyResponse)
	g.leaveUserFromPartyCh <- leaveUserFromPartyRequest{
		resp:   respCh,
		userID: userID,
	}
	resp := <-respCh
	return resp.status
}

func (g *game) RequestP2PRequest(userID string, targetID string, offer string) bool {
	respCh := make(chan requestP2PResponse)
	g.requestP2PCh <- requestP2PRequest{
		resp:     respCh,
		userID:   userID,
		targetID: targetID,
		offer:    offer,
	}
	resp := <-respCh
	return resp.status
}

func (g *game) ResponseP2PRequest(userID string, targetID string, answer string) bool {
	respCh := make(chan responseP2PResponse)
	g.responseP2PCh <- responseP2PRequest{
		resp:     respCh,
		userID:   userID,
		targetID: targetID,
		answer:   answer,
	}
	resp := <-respCh
	return resp.status
}
