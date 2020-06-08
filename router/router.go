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
	rr, err := repositories.NewRoomRepository(a)
	if err != nil {
		return nil, err
	}
	ttnr, err := repositories.NewTicketTypeNameRepository(a)
	if err != nil {
		return nil, err
	}

	// NOTE Create services here
	rs := services.NewRoomService(a, rr)
	ttns := services.NewTicketTypeNameService(a, ttnr)

	// NOTE Create controllers here
	rc := controllers.NewRoomController(a, rs)
	ttnc := controllers.NewTicketTypeNameController(a, ttns)

	// NOTE Add routes here
	router.POST("/rooms/put", rc.PutRoom)
	router.GET("/rooms/getByCoords", rc.GetRoomsByCoords)
	router.GET("/rooms/:id", rc.GetRoomByID)
	router.GET("/rooms", rc.GetAllRooms)
	router.GET("/rooms/byCategory/:category_id", rc.GetRoomsByCategoryID)
	router.GET("/rooms/byParentGroupId/:parent_group_id", rc.GetRoomsByParentGroupID)
	router.PUT("/rooms/update", rc.UpdateRoom) // TODO change to "PUT /rooms"
	router.DELETE("/rooms/:id", rc.DeleteRoom)

	router.POST("/ticketTypeNames", ttnc.PutTicketTypeName)
	router.GET("/ticketTypeNames", ttnc.GetTicketTypeName)
	router.PUT("/ticketTypeNames", ttnc.UpdateTicketTypeName)
	router.DELETE("/ticketTypeNames", ttnc.DeleteTicketTypeName)


	return router, nil
}
