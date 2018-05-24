package server

type Event interface {
}

// Status 操作の成功可否
type Status string

// Action 操作名
type Action string

const (
	// StatusOK 操作が成功
	StatusOK Status = "ok"
	// StatusNG 操作が失敗
	StatusNG Status = "ng"
	// ActionJoinUserToRobby ロビーに参加
	ActionJoinUserToRobby Action = "join_user_to_robby"
	// ActionCreateParty 部屋を建てる
	ActionCreateParty Action = "create_party"
	// ActionJoinParty 部屋に参加する
	ActionJoinParty Action = "join_perty"
	// ActionJoinRoom 部屋に参加
	ActionJoinRoom Action = "join_party"
	// ActionLeaveRoom 部屋から退席
	ActionLeaveRoom Action = "leave_party"
)

// EventMessage イベントのメッセージ
type EventMessage struct {
	ID     string      `json:"id"`
	Action Action      `json:"action"`
	Status Status      `json:"status"`
	Param  interface{} `json:"param"`
}
