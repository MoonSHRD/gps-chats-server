package services

import (
	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
	"github.com/MoonSHRD/sonis/repositories"
)

type TicketTypeNameService struct {
	app                      *app.App
	ticketTypeNameRepositoru *repositories.TicketTypeNameRepository
}

func NewTicketTypeNameService(a *app.App, ttnr *repositories.TicketTypeNameRepository) *TicketTypeNameService {
	return &TicketTypeNameService{
		app:                      a,
		ticketTypeNameRepositoru: ttnr,
	}
}

func (ttns *TicketTypeNameService) Put(ttn *models.TicketTypeName) error {
	return ttns.ticketTypeNameRepositoru.Put(ttn)
}

func (ttns *TicketTypeNameService) Get(eventID string, typeID int) (*models.TicketTypeName, error) {
	return ttns.ticketTypeNameRepositoru.Get(eventID, typeID)
}

func (ttns *TicketTypeNameService) Update(ttn *models.TicketTypeName) error {
	return ttns.ticketTypeNameRepositoru.Update(ttn)
}

func (ttns *TicketTypeNameService) Delete(ttn *models.TicketTypeName) error {
	return ttns.ticketTypeNameRepositoru.Delete(ttn)
}