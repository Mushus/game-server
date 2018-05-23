package server

import "sync"

type Server interface {
}

type server struct {
	games map[string]Game
}

func NewServer() Server {
	return &server{}
}

func (s *server) Start() {
	wg := &sync.WaitGroup{}
	for _, g := range s.games {
		wg.Add(1)
		gamePtr := g.(*game)
		go func() {
			gamePtr.start()
			wg.Done()
		}()
	}
	wg.Wait()
}

func (s *server) AddGame(id string, game Game) {
	s.games[id] = game
}
