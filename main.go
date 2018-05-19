package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	gameIDs = []string{"hoge"}
	repo    Repository
)

func main() {
	repo = NewRepository(gameIDs)

	go func() {
		for {
			RoomManagement()
			time.Sleep(1000)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS())
	//e.use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/:gameID", ShowRobby)
	e.POST("/:gameID/room", CreateRoom)
	e.OPTIONS("/:gameID/room", CreateRoom)
	e.GET("/:gameID/room/:roomID", ShowRoom)

	e.Logger.Fatal(e.Start(":8090"))
}

func RoomManagement() {

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
