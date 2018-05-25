package server

type createUserRequest struct {
	resp     chan createUserResponse
	userName string
	event    chan interface{}
}

type createUserResponse struct {
	user UserView
}

type leaveUserFromGameRequest struct {
	resp   chan leaveUserFromGameResponse
	userID string
}

type leaveUserFromGameResponse struct {
}

type createPartyRequest struct {
	resp      chan createPartyResponse
	userID    string
	isPrivate bool
	maxUsers  int
}

type createPartyResponse struct {
	party  *PartyView
	status bool
}

type joinPartyRequest struct {
	resp    chan joinPartyResponse
	userID  string
	partyID string
}

type joinPartyResponse struct {
	party  *PartyView
	status bool
}

type leaveUserFromPartyRequest struct {
	resp   chan leaveUserFromPartyResponse
	userID string
}

type leaveUserFromPartyResponse struct {
	status bool
}
