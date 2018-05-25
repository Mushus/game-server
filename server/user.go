package server

type user struct {
	id    string
	name  string
	party *party
	event chan interface{}
}

// UserView jsonに変換するためのユーザー情報
type UserView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *user) ToView() UserView {
	return UserView{
		ID:   u.id,
		Name: u.name,
	}
}

func (u *user) Send(event interface{}) {
	u.event <- event
}
