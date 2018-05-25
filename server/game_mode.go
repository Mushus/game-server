package server

type GameMode interface {
	GetID() string
}

type gameMode struct {
	name         string
	id           string
	maxUser      int
	maxPartyUser int
	rooms        map[Room]struct{}
}

func NewGameMode(name string, maxUser int, maxPartyUser int) GameMode {
	return &gameMode{
		name:         name,
		maxUser:      maxUser,
		maxPartyUser: maxPartyUser,
		rooms:        map[Room]struct{}{},
	}
}

func (g *gameMode) GetID() string {
	return g.id
}
