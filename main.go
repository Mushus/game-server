package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mushus/game-server/server"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

var (
	robbyIDs = []string{"hoge"}
)

func main() {

	srv := server.NewServer()
	srv.AddGame(
		"hoge",
		server.NewGame([]server.GameMode{
			server.NewGameMode("simple", 1, 1),
		}),
	)
	go srv.Start()

	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/:gameID", Connect(srv))

	e.Logger.Fatal(e.Start("localhost:8090"))
}

// Connect クライアントとの通信のエンドポイント
func Connect(srv server.Server) func(echo.Context) error {
	return func(c echo.Context) error {
		gameID := c.Param("gameID")
		game := srv.GetGame(gameID)
		if game == nil {
			return c.JSON(
				http.StatusNotFound,
				map[string]interface{}{
					"message": "game not found",
				},
			)
		}

		// userName をリクエストから取り出す
		userName := c.FormValue("userName")
		if userName == "" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"message": "invalid or empty user name",
				},
			)
		}

		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			event := make(chan interface{})
			defer close(event)
			user := game.CreateUserRequest(userName, event)
			userID := user.ID
			defer game.LeaveUserFromGameRequest(userID)
			// イベントを受け取って、レスポンスを返す
			go func() {
				for resp := range event {
					switch r := resp.(type) {
					case server.ModifyPartyEvent:
						websocket.JSON.Send(ws, Response{
							Event: EventModifyParty,
							Param: r.Party,
						})
					case server.RequestP2PEvent:
						websocket.JSON.Send(ws, Response{
							Event: EventRequestP2P,
							Param: EventParamRequestP2P{
								Offer: r.Offer,
							},
						})
					case server.ResponseP2PEvent:
						websocket.JSON.Send(ws, Response{
							Event: EventResponseP2P,
							Param: EventParamResponseP2P{
								Answer: r.Answer,
							},
						})
					}
				}
			}()
			// wsのリクエストを処理する
			for {
				req := &Request{}
				err := websocket.JSON.Receive(ws, &req)
				if err != nil {
					break
				}

				var resp interface{}
				status := false
				switch req.Action {
				case ActionCreateParty:
					param := &ParamCreateParty{}
					json.Unmarshal(*req.Param, &param)
					resp, status = game.CreatePartyRequest(userID, param.IsPrivate, param.MaxUsers)
				case ActionJoinParty:
					param := &ParamJoinPerty{}
					json.Unmarshal(*req.Param, &param)
					resp, status = game.JoinPartyRequest(userID, param.PartyID)
				case ActionLeaveUserFromParty:
					status = game.LeaveUserFromPartyRequest(userID)
				case ActionRequestP2P:
					param := &ParamRequestP2P{}
					json.Unmarshal(*req.Param, &param)
					status = game.RequestP2PRequest(userID, param.TargetID, param.Offer)
				case ActionResponseP2P:
					param := &ParamResponseP2P{}
					json.Unmarshal(*req.Param, &param)
					status = game.ResponseP2PRequest(userID, param.TargetID, param.Answer)
				}
				statusText := StatusNG
				if status {
					statusText = StatusOK
				}
				websocket.JSON.Send(ws, Response{
					ID:     req.ID,
					Status: statusText,
					Param:  resp,
				})
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
