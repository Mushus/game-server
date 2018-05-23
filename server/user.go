package server

import "golang.org/x/net/websocket"

type User interface {
}

type user struct {
	ws *websocket.Conn
}
