package router

import (
	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/controllers"
	"github.com/MoonSHRD/sonis/repositories"
	"github.com/MoonSHRD/sonis/services"
	"github.com/labstack/echo/v4"
)

func NewRouter(a *app.App) (*echo.Echo, error) {
	router := echo.New()
	router.HideBanner = true

	// NOTE Create repositories here
	ccr := repositories.NewChatCategoryRepository(a)
	rr := repositories.NewRoomRepository(a, ccr)

	// NOTE Create services here
	rs := services.NewRoomService(a, rr, ccr)

	// NOTE Create controllers here
	rc := controllers.NewRoomController(a, rs)

	// NOTE Add routes here
	router.POST("/rooms/put", rc.PutRoom)
	router.GET("/rooms/getByCoords", rc.GetRoomsByCoords)
	router.GET("/rooms/:id", rc.GetRoomByID)
	router.GET("/rooms", rc.GetAllRooms)
	router.GET("/rooms/byCategory/:category_id", rc.GetRoomsByCategoryID)
	router.GET("/categories", rc.GetAllCategories)
	router.GET("/rooms/byParentGroupId/:parent_group_id", rc.GetRoomsByParentGroupID)
	router.PUT("/rooms/update", rc.UpdateRoom) // TODO change to "PUT /rooms"
	router.DELETE("/rooms/:id", rc.DeleteRoom)

	return router, nil
}
