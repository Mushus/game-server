package server

// CreatePartyRequest パーティ作成リクエスト
type CreatePartyRequest struct {
	ID        string
	User      User
	IsPrivate bool
	MaxUsers  int
}

// JoinPartyRequest パーティ参加リクエスト
type JoinPartyRequest struct {
	ID      string
	User    User
	PartyID string
}

// LeavePartyRequest パーティから退室リクエスト
type LeavePartyRequest struct {
	ID   string
	User User
}
