package server

type Event interface {
}

type Status string
type Action string

const (
	StatusOK        Status = "ok"
	StatusNG        Status = "ng"
	ActionJoinRoom  Action = "join_party"
	ActionLeaveRoom Action = "leave_party"
)

type EventMessage struct {
	Action Action      `json:"action"`
	Status Status      `json:"status"`
	Param  interface{} `json:"param"`
}
