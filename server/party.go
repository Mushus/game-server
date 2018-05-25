package server

type party struct {
	id        string
	owner     *user
	users     map[string]*user
	maxUsers  int
	isPrivate bool
}

// PartyView パーティJSON
type PartyView struct {
	ID        string     `json:"id"`
	Owner     UserView   `json:"owner"`
	IsPrivate bool       `json:"isPrivate"`
	Users     []UserView `json:"users"`
	MaxUsers  int        `json:"maxUsers"`
}

func (p *party) ToView() PartyView {
	users := []UserView{}
	for _, user := range p.users {
		users = append(users, user.ToView())
	}
	return PartyView{
		ID:        p.id,
		Owner:     p.owner.ToView(),
		IsPrivate: p.isPrivate,
		Users:     users,
		MaxUsers:  p.maxUsers,
	}
}
