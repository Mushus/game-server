package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mushus/game-server/server"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

var (
	robbyIDs = []string{"hoge"}
	repo     Repository
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
		log.Printf("%#v", userName)
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

			event := make(chan server.Event)
			user := server.NewUser(userName, event)
			game.JoinUser(user)
			defer game.LeaveUser(user)

			// イベントを受け取って、レスポンスを返す
			go func() {
				for resp := range event {
					websocket.JSON.Send(ws, resp)
				}
			}()
			// wsのリクエストを処理する
			for {
				req := ParamSocketReceive{}
				err := websocket.JSON.Receive(ws, &req)
				if err != nil {
					return
				}
				log.Printf("rq: %#v", req)

				switch req.Action {
				case ReceiveActionCreateParty:
					param := &ParamCreateParty{}
					json.Unmarshal(*req.Param, &param)
					game.CreateParty(server.CreatePartyRequest{
						ID:       req.ID,
						User:     user,
						MaxUsers: param.MaxUsers,
					})
					/*				case ReceiveActionModifyUser:
									param := &ParamModifyUser{}
									json.Unmarshal(*req.Param, &param)
									log.Printf("%#v", param)*/
				case ReceiveActionJoinPerty:
					param := &ParamJoinPerty{}
					json.Unmarshal(*req.Param, &param)
					game.JoinParty(server.JoinPartyRequest{
						ID:      req.ID,
						User:    user,
						PartyID: param.PartyID,
					})
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
