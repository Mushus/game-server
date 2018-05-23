package server

type Party interface {
	GetID() string
}

type party struct {
	id    string
	join  chan User
	leave chan User
	users map[User]struct{}
}

func (p *party) GetID() string {
	return p.id
}

func (p *party) Start() {
	for {
		select {
		case user := <-p.join:
			p.users[user] = struct{}{}
		case user := <-p.leave:
			delete(p.users, user)
		}
	}
}
