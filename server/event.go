package server

// ModifyPartyEvent パーティの更新
type ModifyPartyEvent struct {
	Party PartyView
}

// LeavePartyEvent パーティの退席
type LeavePartyEvent struct {
	Party PartyView
}
