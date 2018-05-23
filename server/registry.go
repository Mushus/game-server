package server

import (
	"github.com/satori/go.uuid"
)

type Registry interface {
	CreateParty() Party
}

type registry struct {
	joinUser    chan User
	leaveUser   chan User
	addParty    chan Party
	removeParty chan Party
	users       map[User]struct{}
	parties     map[string]Party
}

func NewRegistry() Registry {
	return &registry{
		joinUser:    make(chan User),
		leaveUser:   make(chan User),
		addParty:    make(chan Party),
		removeParty: make(chan Party),
		users:       map[User]struct{}{},
		parties:     map[string]Party{},
	}
}

func (r *registry) CreateParty() Party {
	id := uuid.NewV4().String()
	party := &party{
		id:    id,
		join:  make(chan User),
		leave: make(chan User),
	}
	r.addParty <- party

	return party
}

func (r *registry) Start() {
	for {
		select {
		case user := <-r.joinUser:
			r.users[user] = struct{}{}
		case user := <-r.leaveUser:
			delete(r.users, user)
		case party := <-r.addParty:
			r.parties[party.GetID()] = party
		case party := <-r.removeParty:
			delete(r.parties, party.GetID())
		}
	}
}
