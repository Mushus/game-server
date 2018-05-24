package server

import (
	"github.com/satori/go.uuid"
)

// User ユーザーです
type User interface {
	Send(event Event)
	LeaveParty(req LeavePartyRequest)
	GetID() string
	ToView() UserView
}

type user struct {
	id         string
	name       string
	party      Party
	leaveParty chan LeavePartyRequest
	event      chan Event
}

// UserView jsonに変換するためのユーザー情報
type UserView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewUser name のユーザーを作成します
// event のチャンネルにユーザー宛のメッセージが送信されます
func NewUser(name string, event chan Event) User {
	id := uuid.NewV4().String()
	user := &user{
		id:         id,
		name:       name,
		leaveParty: make(chan LeavePartyRequest),
		event:      event,
	}
	go user.Start()
	return user
}

func (u *user) Send(event Event) {
	u.event <- event
}

func (u *user) LeaveParty(req LeavePartyRequest) {
	u.leaveParty <- req
}

func (u *user) GetID() string {
	return u.id
}

func (u *user) ToView() UserView {
	return UserView{
		Name: u.name,
	}
}

func (u *user) Start() {
	for {
		select {
		case req := <-u.leaveParty:
			if u.party == nil {
				u.Send(EventMessage{
					ID:     req.ID,
					Action: ActionLeaveParty,
					Status: StatusNG,
					Param:  struct{}{},
				})
			}
			u.party.Leave(req)
			u.party = nil
			u.Send(EventMessage{
				Action: ActionLeaveParty,
				Status: StatusOK,
				Param:  struct{}{},
			})
		}
	}
}
