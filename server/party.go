package server

// Party パーティ
type Party interface {
	GetID() string
	Join(req JoinPartyRequest)
	ToView() PartyView
}

type party struct {
	id        string
	owner     User
	join      chan JoinPartyRequest
	leave     chan User
	users     map[User]struct{}
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

func (p *party) GetID() string {
	return p.id
}

func (p *party) ToView() PartyView {
	users := []UserView{}
	for user := range p.users {
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

func (p *party) Start() {
	for {
		select {
		case req := <-p.join:
			_, ok := p.users[req.User]
			if ok || p.isPrivate || (p.maxUsers != 0 && len(p.users) >= p.maxUsers) {
				req.User.Send(EventMessage{
					Action: ActionJoinParty,
					Status: StatusNG,
					Param:  struct{}{},
				})
				continue
			}
			p.users[req.User] = struct{}{}
			view := p.ToView()
			for user := range p.users {
				id := ""
				if user == req.User {
					id = req.ID
				}
				user.Send(EventMessage{
					ID:     id,
					Action: ActionJoinParty,
					Status: StatusOK,
					Param:  view,
				})
			}
		case user := <-p.leave:
			if _, ok := p.users[user]; !ok {
				user.Send(EventMessage{
					Action: ActionLeaveRoom,
					Status: StatusNG,
					Param:  struct{}{},
				})
				continue
			}
			delete(p.users, user)
			user.Send(EventMessage{
				Action: ActionLeaveRoom,
				Status: StatusOK,
				Param:  struct{}{},
			})
		}
	}
}

func (p *party) Join(req JoinPartyRequest) {
	p.join <- req
}
