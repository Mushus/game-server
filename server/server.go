package server

import "sync"

// Server ゲームサーバーです
type Server interface {
	Start()
	// AddGame ゲームを追加します。
	AddGame(id string, game Game)
	// GetGame ゲームを取得します
	// 対応した id のゲームが事前に追加されていない場合、nilを返します
	GetGame(id string) Game
}

type server struct {
	games map[string]Game
}

// NewServer ゲームサーバーを作成します
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

func (s *server) GetGame(id string) Game {
	game, ok := s.games[id]
	if !ok {
		return nil
	}
	return game
}
