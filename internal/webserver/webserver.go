package webserver

import (
	"context"

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

func New(db *database.Database) (*Webserver, error) {
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

	webserver.echo.HideBanner = true
	err = webserver.echo.Start(":37642")
	if err != nil {
		return nil, err
	}
	webserver.logger.Info("Started listening on http://localhost:37642/")
	return webserver, nil
}

func (ws *Webserver) Destroy() {
	ws.echo.Shutdown(context.Background())
}
