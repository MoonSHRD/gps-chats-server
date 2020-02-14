package webserver

import (
	"context"
	"fmt"

	"github.com/MoonSHRD/logger"

	"github.com/MoonSHRD/sonis/internal/utils"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/httpHandler"
	echo "github.com/labstack/echo/v4"
)

type Webserver struct {
	echo        *echo.Echo
	httpHandler *httpHandler.HttpHandler
}

func New(cfg *utils.Config, db *database.Database) (*Webserver, error) {
	httpHandler, err := httpHandler.New(db)
	if err != nil {
		return nil, err
	}
	webserver := &Webserver{
		echo:        echo.New(),
		httpHandler: httpHandler,
	}

	webserver.echo.POST("/rooms/put", httpHandler.HandlePutRoomRequest)
	webserver.echo.GET("/rooms/getByCoords", httpHandler.HandleGetRoomsRequest)
	webserver.echo.GET("/rooms/:room_id", httpHandler.HandleGetRoomByRoomID)
	webserver.echo.GET("/rooms", httpHandler.HandleGetAllRooms)
	webserver.echo.GET("/rooms/byCategory/:category_id", httpHandler.HandleGetRoomsByCategoryID)
	webserver.echo.GET("/categories", httpHandler.HandleGetAllCategories)
	webserver.echo.GET("/rooms/byParentGroupId/:parent_group_id", httpHandler.HandleGetRoomsByParentGroupID)

	webserver.echo.HideBanner = true
	go func() {
		logger.Fatal(webserver.echo.Start(fmt.Sprintf(":%d", cfg.WebserverPort)))
	}()
	logger.Infof("Started listening on http://localhost:%d/", cfg.WebserverPort)
	return webserver, nil
}

func (ws *Webserver) Destroy() {
	ws.echo.Shutdown(context.Background())
}
