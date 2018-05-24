package server

// Party パーティ
type Party interface {
	GetID() string
	Join(req JoinPartyRequest)
	Leave(req LeavePartyRequest)
	ToView() PartyView
}

type party struct {
	id        string
	owner     User
	join      chan JoinPartyRequest
	leave     chan LeavePartyRequest
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
				act := ActionModifyParty
				if user == req.User {
					id = req.ID
					act = ActionJoinParty
				}
				user.Send(EventMessage{
					ID:     id,
					Action: act,
					Status: StatusOK,
					Param:  view,
				})
			}
		case req := <-p.leave:
			if _, ok := p.users[req.User]; ok {
				delete(p.users, req.User)
			}
			if p.owner == req.User {
				p.owner = nil
				for user := range p.users {
					p.owner = user
					break
				}
			}
		}
	}
}

func (p *party) Join(req JoinPartyRequest) {
	p.join <- req
}

func (p *party) Leave(req LeavePartyRequest) {
	p.leave <- req
}
