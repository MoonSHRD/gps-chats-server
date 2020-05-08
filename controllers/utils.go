package controllers

import (
	"github.com/MoonSHRD/sonis/models"
	"github.com/labstack/echo/v4"
)

func ReturnHTTPError(ctx echo.Context, err error, code int) {
	res := models.HTTPError{
		Error: models.HTTPErrorMessage{
			Message: err.Error(),
		},
	}
	ctx.JSON(code, res)
}
