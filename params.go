package main

import "encoding/json"

// RoomParam TODO
type RoomParam struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	MaxUsers       int    `json:"maxUsers"`
	IsAutoMatching bool   `json:"isAutoMatching"`
}

const (
	// ReceiveActionCreateParty パーティを作成する
	ReceiveActionCreateParty = "create_party"
)

// ParamSocketReceive ソケットの取得する形式
type ParamSocketReceive struct {
	Action string           `json:"action"`
	Param  *json.RawMessage `json:"param"`
}

// ParamCreateParty TODO
type ParamCreateParty struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	IsPrivate bool   `json:"isPrivate"`
	maxUsers  int    `json:"maxUsers"`
}
