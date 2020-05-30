package controllers

import (
	"fmt"
	"github.com/MoonSHRD/logger"
	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
	"github.com/MoonSHRD/sonis/services"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

type TicketTypeNameController struct {
	app *app.App
	ticketTypeNameService *services.TicketTypeNameService
}

func NewTicketTypeNameController(a *app.App, ttns *services.TicketTypeNameService) *TicketTypeNameController {
	return &TicketTypeNameController{
		app: a,
		ticketTypeNameService: ttns,
	}
}

func (ttnc *TicketTypeNameController) PutTicketTypeName(eCtx echo.Context) error {
	var ttn *models.TicketTypeName
	err := eCtx.Bind(&ttn)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	err = ttnc.ticketTypeNameService.Put(ttn)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	_ = eCtx.JSON(http.StatusOK, ttn)
	return nil
}

func (ttnc *TicketTypeNameController) GetTicketTypeName(eCtx echo.Context) error {
	eventID := eCtx.QueryParam("eventID")
	typeIDStr := eCtx.QueryParam("typeID")
	if eventID == "" || typeIDStr == "" {
		err := fmt.Errorf("query params are empty")
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	typeID, err := strconv.Atoi(typeIDStr)
	if err != nil {
		err := fmt.Errorf("typeID must be an int")
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	ttn, err := ttnc.ticketTypeNameService.Get(eventID, typeID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logger.Error(err.Error())
			ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
			return err
		} else {
			err = fmt.Errorf("ticket type name is not found")
			ReturnHTTPError(eCtx, err, http.StatusBadRequest)
			return err
		}
	}
	_ = eCtx.JSON(http.StatusOK, ttn)
	return nil
}

func (ttnc *TicketTypeNameController) UpdateTicketTypeName(eCtx echo.Context) error {
	var ttn *models.TicketTypeName
	err := eCtx.Bind(&ttn)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	err = ttnc.ticketTypeNameService.Update(ttn)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	_ = eCtx.JSON(http.StatusOK, ttn)
	return nil
}

func (ttnc *TicketTypeNameController) DeleteTicketTypeName(eCtx echo.Context) error {
	var ttn *models.TicketTypeName
	err := eCtx.Bind(&ttn)
	if err != nil {
		ReturnHTTPError(eCtx, err, http.StatusBadRequest)
		return err
	}
	err = ttnc.ticketTypeNameService.Delete(ttn)
	if err != nil {
		logger.Error(err.Error())
		ReturnHTTPError(eCtx, err, http.StatusInternalServerError)
		return err
	}
	res := map[string]interface{} {
		"result": "ticket type name was deleted",
	}
	_ = eCtx.JSON(http.StatusOK, res)
	return nil
}