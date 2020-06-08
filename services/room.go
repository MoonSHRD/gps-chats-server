package services

import (
	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
	"github.com/MoonSHRD/sonis/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService struct {
	app            *app.App
	roomRepository *repositories.RoomRepository
}

func NewRoomService(a *app.App, rr *repositories.RoomRepository) *RoomService {
	return &RoomService{
		app:            a,
		roomRepository: rr,
	}
}

func (rs *RoomService) PutRoom(room *models.Room) (*models.Room, error) {
	return rs.roomRepository.PutRoom(room)
}

func (rs *RoomService) GetRoomsByCoords(lat float64, lon float64, radius int) (*[]models.Room, error) {
	return rs.roomRepository.GetRoomsByCoords(lat, lon, radius)
}

func (rs *RoomService) GetRoomByID(id string) (*models.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return rs.roomRepository.GetRoomByID(objectID)
}

func (rs *RoomService) GetAllRooms() ([]models.Room, error) {
	return rs.roomRepository.GetAllRooms()
}

func (rs *RoomService) GetRoomsByCategoryID(categoryID int) ([]models.Room, error) {
	return rs.roomRepository.GetRoomsByCategoryID(categoryID)
}

func (rs *RoomService) GetRoomsByParentGroupID(parentGroupID string) ([]models.Room, error) {
	return rs.roomRepository.GetRoomsByParentGroupID(parentGroupID)
}

func (rs *RoomService) UpdateRoom(room *models.Room) (*models.Room, error) {
	return rs.roomRepository.UpdateRoom(room)
}

func (rs *RoomService) DeleteRoom(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return rs.roomRepository.DeleteRoom(objectID)
}
