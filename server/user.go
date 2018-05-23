package server

import "golang.org/x/net/websocket"

type User interface {
	Send(event Event)
}

type user struct {
	ws    *websocket.Conn
	event chan Event
}

func (u *user) Send(event Event) {
	u.event <- event
}
