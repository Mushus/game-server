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
	// ReceiveActionModifyUser ユーザーを変更する
	ReceiveActionModifyUser = "modify_user"
	// ReceiveActionCreateParty パーティを作成する
	ReceiveActionCreateParty = "create_party"
	// ReceiveActionGetParty パーティを取得する
	ReceiveActionGetParty = "get_party"
	// ReceiveActionJoinPerty パーティに参加依頼する
	ReceiveActionJoinPerty = "join_party"
	// ReceiveActionLeaveParty パーティに参加依頼する
	ReceiveActionLeaveParty = "leave_party"
)

// ParamSocketReceive ソケットの取得する形式
type ParamSocketReceive struct {
	Action string           `json:"action"`
	ID     string           `json:"id"`
	Param  *json.RawMessage `json:"param"`
}

// ParamModifyUser ユーザーが更新された時のパラメータ
type ParamModifyUser struct{}

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

type ParamJoinPerty struct {
	PartyID string `json:"partyId"`
}

// ErrorResponse エラーレスポン
// エラーの場合はこのjsonがパラメータとして返ります
type ErrorResponse struct {
	Message string `json:"message"`
}

// WebSocketResponse レスポンス
type WebSocketResponse struct {
	Action string                  `json:"action"`
	Status WebScoketResponseStatus `json:"status"`
	ID     string                  `json:"id"`
	Param  interface{}             `json:"param"`
}

// WebScoketResponseStatus レスポンスのステータス
type WebScoketResponseStatus string

const (
	// ResponseStatusOK 成功
	ResponseStatusOK WebScoketResponseStatus = "ok"
	// ResponseStatusNG 失敗
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

// PartyNotFoundErrorResponse パーティは存在しませんエラー
var PartyNotFoundErrorResponse = WebSocketResponse{
	Action: ReceiveActionCreateParty,
	Status: ResponseStatusNG,
	Param: ErrorResponse{
		Message: "party not found",
	},
}
