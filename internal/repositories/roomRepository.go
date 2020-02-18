package repositories

import (
	"fmt"

	"github.com/MoonSHRD/logger"

	"github.com/MoonSHRD/sonis/internal/utils"

	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/models"
)

const (
	// SecondMillisecs represents one second in milliseconds
	SecondMillisecs = 1000

	// Minute represents one minute in seconds
	Minute = 60

	// Hour represents one hour in seconds
	Hour = 60 * Minute

	// Day represents one day in seconds
	Day = 24 * Hour

	// Month represents one month in seconds
	Month = 31 * Day
)

type RoomRepository struct {
	db                     *database.Database
	chatCategoryRepository *ChatCategoryRepository
}

func NewRoomRepository(db *database.Database, chatCategoryRepository *ChatCategoryRepository) (*RoomRepository, error) {
	if db != nil {
		roomRepo := &RoomRepository{
			db:                     db,
			chatCategoryRepository: chatCategoryRepository,
		}
		roomRepo.clearExpiredRecords()
		utils.SetInterval(func(args ...interface{}) {
			roomRepo.clearExpiredRecords()
		}, SecondMillisecs*30, true)
		return roomRepo, nil
	}
	return nil, fmt.Errorf("database connection is null")
}

func (rr *RoomRepository) PutRoom(room *models.Room) (*models.Room, error) {
	if room.TTL <= 0 {
		return nil, fmt.Errorf("TTL is invalid")
	}
	stmt, err := rr.db.GetDatabaseConnection().Preparex(`
		INSERT INTO rooms (latitude, longitude, ttl, room_id, parent_group_id, event_start_date) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at;
	`)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(room.Latitude, room.Longitude, room.TTL, room.RoomID, room.ParentGroupID, room.EventStartDate).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, x := range room.Categories {
		stmt, err := rr.db.GetDatabaseConnection().Preparex("INSERT INTO roomsChatCategoriesLink (categoryId, roomId) VALUES ($1, $2);")
		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec(x.Id, room.ID)
		if err != nil {
			return nil, err
		}
	}

	seconds := room.TTL
	if seconds > Month {
		return nil, fmt.Errorf("creating chats for more than one month is prohibited")
	}
	return room, nil
}

func (rr *RoomRepository) GetRoomsByCoords(userLat float64, userLon float64, radius int) (*[]models.Room, error) {
	var rooms []models.Room
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms WHERE SQRT(POWER(latitude-$1, 2) + POWER(longitude-$2, 2)) < $3;")
	if err != nil {
		return nil, err
	}
	err = stmt.Select(&rooms, userLat, userLon, radius)
	if err != nil {
		return nil, err
	}

	for i, room := range rooms {
		categories, err := rr.getCategoriesByRoomID(room.ID)
		if err != nil {
			return nil, err
		}
		rooms[i].Categories = categories
	}

	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return &rooms, nil
}

func (rr *RoomRepository) GetRoomByRoomID(roomID string) (*models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms WHERE room_id = $1")
	if err != nil {
		return nil, err
	}
	var room models.Room
	err = stmt.Get(&room, roomID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rr *RoomRepository) GetRoomByID(id int) (*models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms WHERE id = $1")
	if err != nil {
		return nil, err
	}
	var room models.Room
	err = stmt.Get(&room, id)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rr *RoomRepository) GetAllRooms() ([]models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex("SELECT * FROM rooms;")
	if err != nil {
		return nil, err
	}
	var rooms []models.Room
	err = stmt.Select(&rooms)
	if err != nil {
		return nil, err
	}

	for i, room := range rooms {
		categories, err := rr.getCategoriesByRoomID(room.ID)
		if err != nil {
			return nil, err
		}
		rooms[i].Categories = categories
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return rooms, nil
}

func (rr *RoomRepository) GetRoomsByCategoryID(id int) ([]models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex(`
		SELECT r.id, r.latitude, r.longitude, r.ttl, r.room_id, r.created_at, r.parent_group_id, r.event_start_date
		FROM rooms as r
		INNER JOIN roomsChatCategoriesLink AS rccl
		ON rccl.categoryId = $1 AND rccl.roomId = r.id;
	`)
	if err != nil {
		return nil, err
	}
	var rooms []models.Room
	err = stmt.Select(&rooms, id)
	if err != nil {
		return nil, err
	}

	for i, room := range rooms {
		categories, err := rr.getCategoriesByRoomID(room.ID)
		if err != nil {
			return nil, err
		}
		rooms[i].Categories = categories
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return rooms, nil
}

func (rr *RoomRepository) GetRoomsByParentGroupID(parentGroupID string) ([]models.Room, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex(`
		SELECT r.id, r.latitude, r.longitude, r.ttl, r.room_id, r.created_at, r.parent_group_id, r.event_start_date
		FROM rooms as r
		INNER JOIN roomsChatCategoriesLink AS rccl
		ON rccl.roomId = r.id
		WHERE r.parent_group_id = $1;
	`)
	if err != nil {
		return nil, err
	}
	var rooms []models.Room
	err = stmt.Select(&rooms, parentGroupID)
	if err != nil {
		return nil, err
	}

	for i, room := range rooms {
		categories, err := rr.getCategoriesByRoomID(room.ID)
		if err != nil {
			return nil, err
		}
		rooms[i].Categories = categories
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return rooms, nil
}

func (rr *RoomRepository) UpdateRoom(room *models.Room) (*models.Room, error) {
	if room.TTL <= 0 {
		return nil, fmt.Errorf("TTL is invalid")
	}
	if room.TTL > Month {
		return nil, fmt.Errorf("creating chats for more than one month is prohibited")
	}
	args := map[string]interface{}{
		"id":               room.ID,
		"latitude":         room.Latitude,
		"longitude":        room.Longitude,
		"ttl":              room.TTL,
		"room_id":          room.RoomID,
		"parent_group_id":  room.ParentGroupID,
		"event_start_date": room.EventStartDate,
	}
	_, err := rr.db.GetDatabaseConnection().NamedExec(
		`update rooms set 
			latitude = :latitude,
			longitude = :longitude,
			ttl = :ttl,
			room_id = :room_id, 
			parent_group_id = :parent_group_id,
			event_start_date = :event_start_date
		where id = :id
	`, args)
	if err != nil {
		return nil, err
	}
	stmt, err := rr.db.GetDatabaseConnection().Preparex("delete from roomschatcategorieslink where roomid = $1")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(room.ID)
	if err != nil {
		return nil, err
	}

	for _, v := range room.Categories {
		stmt, err := rr.db.GetDatabaseConnection().Preparex("insert into roomsChatCategoriesLink (categoryId, roomId) values ($1, $2);")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(v.Id, room.ID)
		if err != nil {
			return nil, err
		}
	}

	return room, nil
}

func (rr *RoomRepository) getCategoriesByRoomID(id int) ([]models.ChatCategory, error) {
	stmt, err := rr.db.GetDatabaseConnection().Preparex(`
			SELECT cc.id, cc.categoryname
			FROM chatCategories AS cc
         	INNER JOIN roomsChatCategoriesLink AS rccl
        	ON rccl.roomId = $1 AND cc.id = rccl.categoryId;
	`)
	if err != nil {
		return nil, err
	}
	var categories []models.ChatCategory
	err = stmt.Select(&categories, id)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = make([]models.ChatCategory, 0)
	}
	return categories, nil
}

func (rr *RoomRepository) clearExpiredRecords() {
	res, err := rr.db.GetDatabaseConnection().Exec(`
		WITH x AS (
			DELETE
			FROM rooms
			WHERE created_at <= now() AT TIME ZONE 'UTC' - interval '1 second' * rooms.ttl
			RETURNING id
		)
		DELETE FROM roomsChatCategoriesLink AS rccl
		USING x
		WHERE rccl.roomID = x.id;
	`)
	if err != nil {
		logger.Errorf("Failed to clear expired records in database! %s", err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Errorf("Failed to clear expired records in database! %s", err.Error())
	}
	if rowsAffected > 0 {
		logger.Infof("Cleaned up database from %d expired records", rowsAffected)
	}
}
