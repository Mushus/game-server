package main

import (
	"encoding/json"
)

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
	// ActionRequestP2P p2p接続をリクエストする
	ActionRequestP2P action = "request_p2p"
	// ActionResponseP2P p2p接続に応答する
	ActionResponseP2P action = "response_p2p"
)

type event string

const (
	// EventCreateUser ユーザーを作成する
	EventCreateUser event = "create_user"
	// EventModifyParty パーティの変更を検知
	EventModifyParty event = "modify_party"
	// EventRequestP2P P2P接続要求検知
	EventRequestP2P event = "request_p2p"
	// EventResponseP2P P2P接続応答検知
	EventResponseP2P event = "response_p2p"
)

// ===========================================================================

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

// ===========================================================================
// RequestParam

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

// ParamRequestP2P p2p接続をリクエストする
type ParamRequestP2P struct {
	UserID string `json:"userId"`
	Offer  string `json:"offer"`
}

// ParamResponseP2P p2p接続に応答する
type ParamResponseP2P struct {
	UserID string `json:"userId"`
	Answer string `json:"Answer"`
}

// ===========================================================================
// EventParam

// EventParamRequestP2P p2p接続をリクエストする
type EventParamRequestP2P struct {
	UserID string `json:"userId"`
	Offer  string `json:"offer"`
}

// EventParamResponseP2P p2p接続に応答する
type EventParamResponseP2P struct {
	UserID string `json:"userId"`
	Answer string `json:"answer"`
}
