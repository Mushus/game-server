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

type event string

const (
	// EventModifyParty パーティの変更を検知
	EventModifyParty event = "modify_party"
)

// Request websocket のリクエスト
type Request struct {
	ID     string           `json:"id"`
	Action action           `json:"action"`
	Param  *json.RawMessage `json:param`
}

// Response websocket のレスポンス
type Response struct {
	Event  event       `json:"event,omitempty"`
	ID     string      `json:"id,omitempty"`
	Status status      `json:"status,omitempty"`
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
