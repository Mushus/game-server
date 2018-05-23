package server

type Party interface {
	GetID() string
}

type party struct {
	id       string
	join     chan User
	leave    chan User
	users    map[User]struct{}
	maxUsers int
}

func (p *party) GetID() string {
	return p.id
}

func (p *party) Start() {
	for {
		select {
		case user := <-p.join:
			if len(p.users) < p.maxUsers {
				p.users[user] = struct{}{}
				user.Send(EventMessage{
					Action: ActionJoinRoom,
					Status: StatusOK,
					Param:  struct{}{},
				})
				continue
			}
			user.Send(EventMessage{
				Action: ActionJoinRoom,
				Status: StatusNG,
				Param:  struct{}{},
			})
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
