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

func (rr *RoomRepository) PutRoom(room *models.Room) error {
	err := rr.db.GetDatabaseConnection().Insert(room)
	if err != nil {
		return err
	}
	rr.deletingRoomScheduler.Delay().Second(room.TTL).Do(func() {
		err := rr.db.GetDatabaseConnection().Delete(room)
		if err != nil {
			rr.logger.Error("Cannot delete room " + room.RoomID + ". Reason: " + err.Error())
		}
	})
	return nil
}
