package main

import (
	"encoding/json"
	"fmt"
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
	//e.use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/:gameID", ShowRobby)
	e.POST("/:gameID/room", CreateRoom)
	e.GET("/:gameID/room/:roomID/join", JoinRoom)
	e.GET("/:gameID/room/:roomID", ShowRoom)

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
				rsvPrm := ParamCreateParty{}
				// NOTE: 既にパースされてるのでエラーの確認は不要のはず
				json.Unmarshal(*srp.Param, &rsvPrm)
				Close(closeRobby)
				myParty = robby.CreateParty(rsvPrm.Name, rsvPrm.Password, rsvPrm.IsPrivate, rsvPrm.maxUsers)
				closeRobby = myParty.Join()
				rsp = myParty.ToView()
			default:
				continue
			}
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

// ShowRobby ルーム一覧を表示する
func ShowRobby(c echo.Context) error {
	gameID := c.Param("gameID")

	robby := repo.GetRobby(gameID)
	if robby == nil {
		return c.String(http.StatusNotFound, "robby not found")
	}

	return c.JSON(http.StatusOK, robby.ToView())
}

// CreateRoom 部屋を作る
func CreateRoom(c echo.Context) error {
	gameID := c.Param("gameID")

	rp := RoomParam{}
	if err := c.Bind(&rp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid room info",
		})
	}

	robby := repo.GetRobby(gameID)
	if robby == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "game not found",
		})
	}

	room := robby.CreateRoom(rp.Name, rp.Password, rp.MaxUsers, rp.IsAutoMatching)
	return c.JSON(http.StatusOK, room.ToView())
}

// ShowRoom 部屋の詳細を表示する
func ShowRoom(c echo.Context) error {
	gameID := c.Param("gameID")
	roomID := c.Param("roomID")

	robby := repo.GetRobby(gameID)
	if robby == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "game not found",
		})
	}

	room := robby.GetRoom(roomID)
	if room == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "room not found",
		})
	}

	return c.JSON(http.StatusOK, room.ToView())
}

// JoinRoom 部屋に入る
func JoinRoom(c echo.Context) error {
	gameID := c.Param("gameID")
	roomID := c.Param("roomID")

	robby := repo.GetRobby(gameID)
	if robby == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "game not found",
		})
	}

	room := robby.GetRoom(roomID)
	if room == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "room not found",
		})
	}

	if !room.CanJoin() {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "cannot join this room",
		})
	}

	room.Join()
	websocket.Handler(func(ws *websocket.Conn) {
		defer room.Leave()
		defer ws.Close()
		fmt.Printf("%v\n", ws.IsClientConn())
		fmt.Printf("%v\n", ws.IsServerConn())
		for {
			err := websocket.JSON.Send(ws, map[string]interface{}{"message": "hello"})
			if err != nil {
				fmt.Print("ended")
				return
			}
			time.Sleep(time.Second)
		}
		fmt.Print("hoge")
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
