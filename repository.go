package main

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// ========================================

type Repository interface {
	GetRobby(gameID string) Robby
}

type Robby interface {
	ToView() RobbyView
	GetRoom(roomID string) Room
	GetRoomViews() []RoomView
	CreateRoom(name string, password string, maxUser int, isAutoMatching bool) Room
}

type Room interface {
	ToView() RoomView
}

// ========================================

// Repository 部屋情報
type repository struct {
	// 部屋情報
	robby map[string]robby
	mu    *sync.RWMutex
}

// GameRobby 部屋一覧
type robby struct {
	rooms map[string]room
	mu    *sync.RWMutex
}

type room struct {
	id             string
	name           string
	password       string
	maxUsers       int
	isAutoMatching bool
	mu             *sync.RWMutex
}

// ========================================

// RobbyView 部屋一覧
type RobbyView struct {
	// 部屋一覧
	Rooms []RoomView `json:"rooms"`
}

// RoomView 部屋
type RoomView struct {
	ID string `json:"id"`
	// Name　部屋名
	Name string `json:"name"`
	// Password　パスワード
	HasPassword bool `json:"has_passwrod"`
	// MaxUsers　ユーザー数
	MaxUser int `json:"max_user"`
	// IsAutoMatching　オートマッチング対応
	IsAutoMatching bool `json:"is_auto_matching"`
}

// ========================================

// NewRepository リポジトリを初期化する
func NewRepository(gameIDs []string) Repository {
	rby := map[string]robby{}
	for _, v := range gameIDs {
		rby[v] = robby{
			rooms: map[string]room{},
			mu:    &sync.RWMutex{},
		}
	}

	return &repository{
		robby: rby,
		mu:    &sync.RWMutex{},
	}
}

// GetRobby 部屋情報を取得する
func (r *repository) GetRobby(gameID string) Robby {
	r.mu.RLock()
	defer r.mu.RUnlock()

	robby, ok := r.robby[gameID]
	if !ok {
		return nil
	}
	return &robby
}

// ========================================

func (r *robby) ToView() RobbyView {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return RobbyView{
		Rooms: r.getRoomViews(),
	}
}

func (r *robby) getRoom(roomID string) Room {
	room, ok := r.rooms[roomID]
	if !ok {
		return nil
	}
	return &room
}

func (r *robby) GetRoom(roomID string) Room {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getRoom(roomID)
}

func (r *robby) getRoomViews() []RoomView {
	rv := []RoomView{}
	for _, v := range r.rooms {
		rv = append(rv, v.ToView())
	}
	return rv
}

// GetRoomViews 部屋一覧を取得する
func (r *robby) GetRoomViews() []RoomView {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getRoomViews()
}

// CreateRoom 部屋を立てる
func (r *robby) CreateRoom(name string, password string, maxUser int, isAutoMatching bool) Room {
	id := uuid.NewV4().String()

	room := room{
		id:             id,
		name:           name,
		password:       password,
		maxUsers:       maxUser,
		isAutoMatching: isAutoMatching,
		mu:             &sync.RWMutex{},
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[id] = room

	return &room
}

// ========================================

func (r *room) ToView() RoomView {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return RoomView{
		ID:             r.id,
		Name:           r.name,
		HasPassword:    r.password != "",
		MaxUser:        r.maxUsers,
		IsAutoMatching: r.isAutoMatching,
	}
}
