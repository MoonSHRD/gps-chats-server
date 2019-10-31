package httpHandler

import (
	"net/http"
	"strconv"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/models"
	"github.com/MoonSHRD/sonis/internal/repositories"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	PutRoomError  = 1
	GetRoomsError = 2
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
		h.logger.Errorf("Processing of /rooms/put request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusInternalServerError, makeHTTPError(PutRoomError, err.Error()))
		return err
	}
	h.roomRepository.PutRoom(&room)
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (h *HttpHandler) HandleGetRoomsRequest(eCtx echo.Context) error {
	userLatStr := eCtx.Param("gps_lat")
	userLonStr := eCtx.Param("gps_lon")
	radiusStr := eCtx.Param("radius")

	userLat, err := strconv.ParseFloat(userLatStr, 32)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/getByCoords request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusBadRequest, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}
	userLon, err := strconv.ParseFloat(userLonStr, 32)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/getByCoords request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusBadRequest, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}
	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/getByCoords request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusBadRequest, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}

	rooms, err := h.roomRepository.GetRoomsByCoords(userLat, userLon, radius)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/getByCoords request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusInternalServerError, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}
	eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func makeHTTPError(errCode int, errText string) *models.HTTPError {
	return &models.HTTPError{
		Ok:      false,
		ErrCode: errCode,
		ErrText: errText,
	}
}
