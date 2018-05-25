package main

import "encoding/json"

type status string

const (
	// StatusOK ok
	StatusOK status = "ok"
	// StatusNG ng
	StatusNG status = "ng"
)

type action string

const (
	// ActionCreateParty パーティ作成
	ActionCreateParty action = "create_party"
	// ActionJoinParty パーティに参加する
	ActionJoinParty action = "join_party"
	// ActionLeaveUserFromParty パーティから退出する
	ActionLeaveUserFromParty action = "leave_user_from_party"
)

// Request websocket のリクエスト
type Request struct {
	ID     string           `json:"id"`
	Action action           `json:"action"`
	Param  *json.RawMessage `json:param`
}

// Response websocket のレスポンス
type Response struct {
	ID     string      `json:"id"`
	Status status      `json:"status"`
	Param  interface{} `json:"param"`
}

// ParamCreateParty パーティを作成する
type ParamCreateParty struct {
	// IsPrivate パーティに入れるかどうか
	IsPrivate bool `json:"isPrivate"`
	// maxUsers パーティの人数制限
	MaxUsers int `json:"maxUsers"`
}

// ParamJoinPerty パーティに参加する
type ParamJoinPerty struct {
	PartyID string `json:"partyId"`
}
