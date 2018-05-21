package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

var (
	robbyIDs = []string{"hoge"}
	repo     Repository
)

func main() {
	repo = NewRepository(robbyIDs)

	go func() {
		for {
			roomManagement()
			time.Sleep(1000)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/:robbyID/matching", Matching)

	e.Logger.Fatal(e.Start(":8090"))
}

func roomManagement() {

}

// Matching マッチングのwebsocket
func Matching(c echo.Context) error {
	robbyID := c.Param("robbyID")
	robby := repo.GetRobby(robbyID)
	if robby == nil {
		return c.String(http.StatusNotFound, "robby not found")
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var myParty Party
		var partyClose func()

		defer Close(partyClose)

		// ロビーを購読して変化があったらクライアントに伝える
		closeRobby := robby.Listen(func(r Robby) {
			// TODO: エラーの検知は不要？
			websocket.JSON.Send(ws, r.ToView())
		})
		defer closeRobby()

		for {
			srp := ParamSocketReceive{}
			err := websocket.JSON.Receive(ws, &srp)
			if err != nil {
				return
			}

			var rsp interface{}
			switch srp.Action {
			case ReceiveActionCreateParty:
				// パーティを作る
				rsvPrm := ParamCreateParty{}
				// パラメータが存在しない
				if srp.Param == nil {
					websocket.JSON.Send(ws, WebSocketResponse{
						Action: ReceiveActionCreateParty,
						Status: ResponseStatusNG,
						Param: ErrorResponse{
							Message: "invalid parameter",
						},
					})
					continue
				}
				// NOTE: 既にパースされてるのでエラーの確認は不要のはず
				json.Unmarshal(*srp.Param, &rsvPrm)
				Close(partyClose)
				myParty = robby.CreateParty(rsvPrm.IsPrivate, rsvPrm.maxUsers)
				partyClose = myParty.Join()
				rsp = WebSocketResponse{
					Action: ReceiveActionCreateParty,
					Status: ResponseStatusOK,
					Param:  myParty.ToView(),
				}
			default:
				continue
			}
			log.Printf("%#v", rsp)
			err = websocket.JSON.Send(ws, rsp)
			if err != nil {
				fmt.Print("ended")
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// Close 閉じる
func Close(fn func()) {
	if fn != nil {
		fn()
	}
}
