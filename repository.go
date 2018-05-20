package main

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// ========================================

// Repository サーバーの状態を管理するところ
type Repository interface {
	GetRobby(robbyID string) Robby
}

// Robby ゲームロビー
type Robby interface {
	ToView() RobbyView
	GetRoom(roomID string) Room
	GetRoomViews() []RoomView
	CreateRoom(name string, password string, maxUser int, isAutoMatching bool) Room
	CreateParty(name string, password string, isPrivate bool, maxUsers int) Party
	Listen(lrf ListenRobbyFunc) (close func())
}

// Room マッチングするための部屋
type Room interface {
	ToView() RoomView
	CanJoin() bool
	Leave()
	Join()
}

// Party 人の集まり
type Party interface {
	ToView() PartyView
	Join() func()
}

// ========================================

// ListenRobbyFunc イベント
type ListenRobbyFunc func(robby Robby)

// ========================================

type repository struct {
	// 部屋情報
	robby map[string]robby
	mu    *sync.RWMutex
}

type robby struct {
	rooms    map[string]room
	party    map[string]party
	listener map[*ListenRobbyFunc]struct{}
	mu       *sync.RWMutex
}

type party struct {
	id        string
	name      string
	password  string
	isPrivate bool
	maxUsers  int
	userCount int
	mu        *sync.RWMutex
}

type room struct {
	id             string
	name           string
	password       string
	maxUsers       int
	usersCount     int
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
	// Name 部屋名
	Name string `json:"name"`
	// Password パスワード
	HasPassword bool `json:"hasPasswrod"`
	// MaxUsers ユーザー数
	MaxUsers int `json:"maxUsers"`
	// IsAutoMatching オートマッチング対応
	IsAutoMatching bool `json:"isAutoMatching"`
}

// PartyView パーティ
type PartyView struct {
	ID string `json:"id"`
	// Name パーティ名
	Name string `json:"name"`
	// Password パスワード
	HasPassword bool `json:"hasPasswrod"`
	// プライベートパーティかどうか
	// ロビーから不可視になります
	isPrivate bool `json:"isPrivate"`
	// MaxUsers　ユーザー数
	MaxUsers  int `json:"maxUsers"`
	UserCount int `json:"maxUsers"`
}

// ========================================

// NewRepository リポジトリを初期化する
func NewRepository(robbyIDs []string) Repository {
	rby := map[string]robby{}
	for _, v := range robbyIDs {
		rby[v] = robby{
			rooms:    map[string]room{},
			listener: map[*ListenRobbyFunc]struct{}{},
			mu:       &sync.RWMutex{},
		}
	}

	return &repository{
		robby: rby,
		mu:    &sync.RWMutex{},
	}
}

// GetRobby 部屋情報を取得する
func (r *repository) GetRobby(robbyID string) Robby {
	r.mu.RLock()
	defer r.mu.RUnlock()

	robby, ok := r.robby[robbyID]
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

func (r *robby) CreateParty(name string, password string, isPrivate bool, maxUsers int) Party {
	partyID := uuid.NewV4().String()

	party := party{
		id:        partyID,
		name:      name,
		password:  password,
		isPrivate: isPrivate,
		maxUsers:  maxUsers,
		userCount: 0,
		mu:        &sync.RWMutex{},
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.party[partyID] = party
	return &party
}

// ロビーの状況を購読する
func (r *robby) Listen(lrf ListenRobbyFunc) (close func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listener[&lrf] = struct{}{}
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.listener, &lrf)
	}
}

// ========================================

func (r *room) ToView() RoomView {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return RoomView{
		ID:             r.id,
		Name:           r.name,
		HasPassword:    r.password != "",
		MaxUsers:       r.maxUsers,
		IsAutoMatching: r.isAutoMatching,
	}
}

func (r *room) CanJoin() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.maxUsers >= r.usersCount || r.usersCount <= 0
}

func (r *room) Join() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.usersCount++
}

func (r *room) Leave() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.usersCount--
}

// ========================================

func (p *party) ToView() PartyView {
	p.mu.RLock()
	defer p.mu.Lock()
	return PartyView{
		ID:          p.id,
		Name:        p.name,
		HasPassword: p.password != "",
		isPrivate:   p.isPrivate,
		MaxUsers:    p.maxUsers,
		UserCount:   0,
	}
}

func (p *party) Join() func() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.userCount++
	return func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.userCount--
	}
}
