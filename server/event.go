package server

type event string

const ()

// EventMessage イベントのメッセージ
type EventMessage struct {
	event  event
	Status bool
	Param  interface{}
}
