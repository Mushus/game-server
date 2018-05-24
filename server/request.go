package server

// CreatePartyRequest パーティ作成リクエスト
type CreatePartyRequest struct {
	ID        string
	User      User
	IsPrivate bool
	MaxUsers  int
}

// JoinPertyRequest パーティ参加リクエスト
type JoinPartyRequest struct {
	ID      string
	User    User
	PartyID string
}
