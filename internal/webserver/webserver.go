package webserver

import (
	"context"
	"fmt"
	"github.com/MoonSHRD/sonis/internal/utils"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/httpHandler"
	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Webserver struct {
	echo        *echo.Echo
	httpHandler *httpHandler.HttpHandler
	logger      *logrus.Logger
}

func New(cfg *utils.Config, db *database.Database) (*Webserver, error) {
	httpHandler, err := httpHandler.New(db)
	if err != nil {
		return nil, err
	}
	webserver := &Webserver{
		echo:        echo.New(),
		httpHandler: httpHandler,
		logger:      logrus.New(),
	}

	webserver.echo.POST("/rooms/put", httpHandler.HandlePutRoomRequest)
	webserver.echo.GET("/rooms/getByCoords", httpHandler.HandleGetRoomsRequest)
	webserver.echo.GET("/rooms/:room_id", httpHandler.HandleGetRoomByRoomID)
	webserver.echo.GET("/rooms", httpHandler.HandleGetAllRooms)
	webserver.echo.GET("/rooms/byCategory/:category_id", httpHandler.HandleGetRoomsByCategoryID)
	webserver.echo.GET("/categories", httpHandler.HandleGetAllCategories)
	webserver.echo.GET("/rooms/byParentGroupId/:parent_group_id", httpHandler.HandleGetRoomsByParentGroupID)

	webserver.echo.HideBanner = true
	if cfg.WebserverPort <= 0 || cfg.WebserverPort > 65535 {
		return nil, fmt.Errorf("incorrect port %d", cfg.WebserverPort)
	}
	err = webserver.echo.Start(fmt.Sprintf(":%d", cfg.WebserverPort))
	if err != nil {
		return nil, err
	}
	webserver.logger.Info("Started listening on http://localhost:37642/")
	return webserver, nil
}

func (ws *Webserver) Destroy() {
	ws.echo.Shutdown(context.Background())
}
