package repositories

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/models"
	"github.com/prprprus/scheduler"
)

type RoomRepository struct {
	db                    *database.Database
	deletingRoomScheduler *scheduler.Scheduler
	logger                *logrus.Logger
}

func NewRoomRepository(db *database.Database) (*RoomRepository, error) {
	if db != nil {
		deletingRoomScheduler, err := scheduler.NewScheduler(10000)
		if err != nil {
			return nil, err
		}
		return &RoomRepository{
			db:                    db,
			deletingRoomScheduler: deletingRoomScheduler,
			logger:                logrus.New(),
		}, nil
	}
	return nil, fmt.Errorf("database connection is null")
}

func (rr *RoomRepository) PutRoom(room *models.Room) (*models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex("INSERT INTO rooms (latitude, longitude, ttl, room_id, category) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;")
	if err != nil {
		return nil, err
	}
	stmt.QueryRow(room.Latitude, room.Longitude, room.TTL, room.RoomID, room.Category).Scan(&room.ID, &room.CreatedAt)
	rr.deletingRoomScheduler.Delay().Second(room.TTL).Do(func() {
		stmt, err := rr.db.GetDatabaseConnection().Preparex("DELETE FROM rooms WHERE id = $1;")
		_, err = stmt.Exec(room.ID)
		if err != nil {
			rr.logger.Errorf("Cannot delete room %d. Reason: %s", room.ID, err.Error())
		}
	})
	return room, nil
}

func (rr *RoomRepository) GetRoomsByCoords(userLat float64, userLon float64, radius int) (*[]models.Room, error) {
	var rooms []models.Room
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms WHERE SQRT(POWER(latitude-$1, 2) + POWER(longitude-$2, 2)) < $3")
	if err != nil {
		return nil, err
	}
	err = stmt.Select(&rooms, userLat, userLon, radius)
	if err != nil {
		return nil, err
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return &rooms, nil
}

func (rr *RoomRepository) GetRoomByRoomID(roomID string) (*models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms WHERE room_id = ?")
	if err != nil {
		return nil, err
	}
	var room models.Room
	stmt.Get(&room, roomID)
	return &room, nil
}
