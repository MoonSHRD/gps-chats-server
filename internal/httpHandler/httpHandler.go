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
	PutRoomError    = 1
	GetRoomsError   = 2
	GetRoomByRoomID = 3
)

type HttpHandler struct {
	roomRepository         *repositories.RoomRepository
	chatCategoryRepository *repositories.ChatCategoryRepository
	logger                 *logrus.Logger
}

func New(db *database.Database) (*HttpHandler, error) {
	chatCategoryRepository, err := repositories.NewChatCategoryRepository(db)
	if err != nil {
		return nil, err
	}
	roomRepository, err := repositories.NewRoomRepository(db, chatCategoryRepository)
	if err != nil {
		return nil, err
	}
	return &HttpHandler{
		roomRepository:         roomRepository,
		chatCategoryRepository: chatCategoryRepository,
		logger:                 logrus.New(),
	}, nil
}

func (h *HttpHandler) HandlePutRoomRequest(eCtx echo.Context) error {
	var room *models.Room
	err := eCtx.Bind(&room)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/put request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusInternalServerError, makeHTTPError(PutRoomError, err.Error()))
		return err
	}
	room, err = h.roomRepository.PutRoom(room)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/put request failed! Reason: %s", err.Error())
		eCtx.JSON(http.StatusInternalServerError, makeHTTPError(PutRoomError, err.Error()))
		return err
	}
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (h *HttpHandler) HandleGetRoomsRequest(eCtx echo.Context) error {
	userLatStr := eCtx.QueryParam("gps_lat")
	userLonStr := eCtx.QueryParam("gps_lon")
	radiusStr := eCtx.QueryParam("radius")

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

func (h *HttpHandler) HandleGetRoomByRoomID(eCtx echo.Context) error {
	roomID := eCtx.Param("room_id")
	room, err := h.roomRepository.GetRoomByRoomID(roomID)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/%s request failed! Reason: %s", roomID, err.Error())
		eCtx.JSON(http.StatusInternalServerError, makeHTTPError(GetRoomByRoomID, err.Error()))
		return err
	}
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (h *HttpHandler) HandleGetAllRooms(eCtx echo.Context) error {
	rooms, err := h.roomRepository.GetAllRooms()
	if err != nil {
		h.logger.Errorf("Processing of /rooms request failed! Reason: %s", err.Error())
		_ = eCtx.JSON(http.StatusInternalServerError, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}
	_ = eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (h *HttpHandler) HandleGetRoomsByCategoryID(eCtx echo.Context) error {
	categoryIDStr := eCtx.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return err
	}
	rooms, err := h.roomRepository.GetRoomsByCategoryID(categoryID)
	if err != nil {
		h.logger.Errorf("Processing of /rooms/byCategory/%d request failed! Reason: %s", categoryID, err.Error())
		_ = eCtx.JSON(http.StatusInternalServerError, makeHTTPError(GetRoomsError, err.Error()))
		return err
	}
	_ = eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (h *HttpHandler) HandleGetAllCategories(eCtx echo.Context) error {
	categories, err := h.chatCategoryRepository.GetAllCategories()
	if err != nil {
		h.logger.Errorf("Processing of /categories request failed! Reason: %s", err.Error())
		_ = eCtx.JSON(http.StatusInternalServerError, makeHTTPError(GetRoomByRoomID, err.Error()))
		return err
	}
	_ = eCtx.JSON(http.StatusOK, categories)
	return nil
}

func makeHTTPError(errCode int, errText string) *models.HTTPError {
	return &models.HTTPError{
		Ok:      false,
		ErrCode: errCode,
		ErrText: errText,
	}
}
