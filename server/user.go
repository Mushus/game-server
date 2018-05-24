package server

// User ユーザーです
type User interface {
	Send(event Event)
	ToView() UserView
}

type user struct {
	id    string
	name  string
	event chan Event
}

// UserView jsonに変換するためのユーザー情報
type UserView struct {
	Name string `json:"name"`
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

func (u *user) ToView() UserView {
	return UserView{
		Name: u.name,
	}
}
