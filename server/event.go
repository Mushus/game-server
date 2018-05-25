package server

// ModifyPartyEvent パーティの更新
type ModifyPartyEvent struct {
	Party PartyView
}

// LeavePartyEvent パーティの退席
type LeavePartyEvent struct {
	Party PartyView
}

// RequestP2PEvent P2P接続要求
type RequestP2PEvent struct {
	Offer  string
	UserID string
}

// ResponseP2PEvent P2P接続応答
type ResponseP2PEvent struct {
	Answer string
	UserID string
}
