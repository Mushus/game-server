package server

type user struct {
	id    string
	name  string
	party *party
	event chan EventMessage
}

// UserView jsonに変換するためのユーザー情報
type UserView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *user) ToView() UserView {
	return UserView{
		Name: u.name,
	}
}

func (u *user) Send(event EventMessage) {
	u.event <- event
}
