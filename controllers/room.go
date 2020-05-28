package controllers

import (
	"net/http"
	"strconv"

	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
	"github.com/MoonSHRD/sonis/services"
	"github.com/google/logger"
	"github.com/labstack/echo/v4"
)

type RoomController struct {
	app         *app.App
	roomService *services.RoomService
}

func NewRoomController(a *app.App, rs *services.RoomService) *RoomController {
	return &RoomController{
		app:         a,
		roomService: rs,
	}
}

func (rc *RoomController) PutRoom(eCtx echo.Context) error {
	var room *models.Room
	err := eCtx.Bind(&room)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	room, err = rc.roomService.PutRoom(room)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (rc *RoomController) GetRoomsByCoords(eCtx echo.Context) error {
	latStr := eCtx.QueryParam("gps_lat")
	lonStr := eCtx.QueryParam("gps_lon")
	radiusStr := eCtx.QueryParam("radius")

	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	lon, err := strconv.ParseFloat(lonStr, 32)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}

	rooms, err := rc.roomService.GetRoomsByCoords(lat, lon, radius)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (rc *RoomController) GetRoomByID(eCtx echo.Context) error {
	id := eCtx.Param("id")
	room, err := rc.roomService.GetRoomByID(id)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (rc *RoomController) GetAllRooms(eCtx echo.Context) error {
	rooms, err := rc.roomService.GetAllRooms()
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	_ = eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (rc *RoomController) GetRoomsByCategoryID(eCtx echo.Context) error {
	categoryIDStr := eCtx.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return err
	}
	rooms, err := rc.roomService.GetRoomsByCategoryID(categoryID)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	_ = eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (rc *RoomController) GetRoomsByParentGroupID(eCtx echo.Context) error {
	categoryID := eCtx.Param("parent_group_id")
	rooms, err := rc.roomService.GetRoomsByParentGroupID(categoryID)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	_ = eCtx.JSON(http.StatusOK, rooms)
	return nil
}

func (rc *RoomController) UpdateRoom(eCtx echo.Context) error {
	var room *models.Room
	err := eCtx.Bind(&room)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	room, err = rc.roomService.UpdateRoom(room)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	eCtx.JSON(http.StatusOK, room)
	return nil
}

func (rc *RoomController) DeleteRoom(eCtx echo.Context) error {
	id := eCtx.Param("id")
	err := rc.roomService.DeleteRoom(id)
	if err != nil {
		logger.Errorf(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	return eCtx.String(http.StatusOK, "")
}
