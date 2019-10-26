package httpHandler

import (
	"net/http"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/models"
	"github.com/MoonSHRD/sonis/internal/repositories"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type HttpHandler struct {
	roomRepository *repositories.RoomRepository
	logger         *logrus.Logger
}

func New(db *database.Database) (*HttpHandler, error) {
	roomRepository, err := repositories.NewRoomRepository(db)
	if err != nil {
		return nil, err
	}
	return &HttpHandler{
		roomRepository: roomRepository,
		logger:         logrus.New(),
	}, nil
}

func (h *HttpHandler) HandlePutRoomRequest(eCtx echo.Context) error {
	var room models.Room
	err := eCtx.Bind(&room)
	if err != nil {
		h.logger.Error("Processing of /putRoom request failed! Reason: " + err.Error())
		httpError := models.HttpError{
			Ok:      false,
			ErrCode: -1,
			ErrText: err.Error(),
		}
		eCtx.JSON(http.StatusInternalServerError, httpError)
		return err
	}
	h.roomRepository.PutRoom(&room)
	eCtx.JSON(http.StatusOK, room)
	return nil
}
