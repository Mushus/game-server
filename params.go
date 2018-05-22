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
	// ReceiveActionGetParty パーティを取得する
	ReceiveActionGetParty = "get_party"
)

// ParamSocketReceive ソケットの取得する形式
type ParamSocketReceive struct {
	Action string           `json:"action"`
	ID     string           `json:"id"`
	Param  *json.RawMessage `json:"param"`
}

// ParamCreateParty パーティを作成する
type ParamCreateParty struct {
	// IsPrivate パーティに入れるかどうか
	IsPrivate bool `json:"isPrivate"`
	// maxUsers パーティの人数制限
	MaxUsers int `json:"maxUsers"`
}

// ParamGetParty パーティを取得する
type ParamGetParty struct {
	// PartyID パーティID
	PartyID string `json:"partyId"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
type WebSocketResponse struct {
	Action string                  `json:"action"`
	Status WebScoketResponseStatus `json:"status"`
	ID     string                  `json:"id"`
	Param  interface{}             `json:"param"`
}

type WebScoketResponseStatus string

const (
	ResponseStatusOK WebScoketResponseStatus = "ok"
	ResponseStatusNG WebScoketResponseStatus = "ng"
)

// InvalidParameterErrorResponse パラメータが間違ってるエラー
var InvalidParameterErrorResponse = WebSocketResponse{
	Action: ReceiveActionCreateParty,
	Status: ResponseStatusNG,
	Param: ErrorResponse{
		Message: "invalid parameter",
	},
}
