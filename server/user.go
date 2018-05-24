package server

// User ユーザーです
type User interface {
	Send(event Event)
}

type user struct {
	name  string
	event chan Event
}

// NewUser name のユーザーを作成します
// event のチャンネルにユーザー宛のメッセージが送信されます
func NewUser(name string, event chan Event) User {
	return &user{
		name:  name,
		event: event,
	}
}

func (u *user) Send(event Event) {
	u.event <- event
}
