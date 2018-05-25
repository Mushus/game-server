package server

// ModifyPartyEvent パーティの更新
type ModifyPartyEvent struct {
	party PartyView
}

// LeavePartyEvent パーティの退席
type LeavePartyEvent struct {
	party PartyView
}
