package server

type Room interface {
	GetUserCount() int
}

type room struct {
	userCount int
}

func (r *room) GetUserCount() int {
	return r.userCount
}

func (r *room) start() {
	/*for {
		select{
		case
		}
	}*/
}
