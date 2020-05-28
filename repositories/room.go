package repositories

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MoonSHRD/sonis/app"

	"github.com/MoonSHRD/sonis/models"
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

	RoomCollectionName = "rooms"
)

type RoomRepository struct {
	app          *app.App
	dbCollection *mongo.Collection
	ctx          context.Context
}

func NewRoomRepository(a *app.App) (*RoomRepository, error) {
	ctx := context.Background()
	col := a.MainDatabase.Collection(RoomCollectionName)
	roomRepo := &RoomRepository{
		app:          a,
		dbCollection: col,
		ctx:          ctx,
	}

	err := roomRepo.initMongo()
	if err != nil {
		return nil, err
	}

	return roomRepo, nil
}

func (rr *RoomRepository) initMongo() error {
	expiresAtIndex := mongo.IndexModel{
		Keys: bson.M{
			"expiresAt": 1,
		}, Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := rr.dbCollection.Indexes().CreateMany(rr.ctx, []mongo.IndexModel{expiresAtIndex})
	return err
}

func (rr *RoomRepository) PutRoom(room *models.Room) (*models.Room, error) {
	room.ID = primitive.NewObjectID()
	err := rr.prepareRoomForInserting(room)
	if err != nil {
		return nil, err
	}
	ttlDuration := time.Duration(room.TTL)
	now := time.Now().UTC()
	expiresAtTime := now.Add(ttlDuration * time.Second)
	room.ExpiresAt = expiresAtTime
	room.CreatedAt = now

	res, err := rr.dbCollection.InsertOne(rr.ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	return room, nil
}

func (rr *RoomRepository) prepareRoomForInserting(room *models.Room) error {
	if room.TTL <= 0 {
		return fmt.Errorf("TTL is invalid")
	}
	if room.TTL > Month {
		return fmt.Errorf("creating chats for more than one month is prohibited")
	}
	for i, v := range room.Categories {
		for k, j := range room.Categories { // check every element if it's already present in the array
			if j.ID == v.ID && i != k {
				return fmt.Errorf("two or more identical categories detected")
			}
		}
	}

	return nil
}

func (rr *RoomRepository) GetRoomsByCoords(userLat float64, userLon float64, radius int) (*[]models.Room, error) {
	var rooms []models.Room
	pipeline := mongo.Pipeline{
		{{
			"$addFields", bson.D{{
				"radius", bson.D{{
					"$sqrt", bson.D{{
						"$add", bson.A{
							bson.D{{
								"$pow", bson.A{
									bson.D{{
										"$subtract", bson.A{
											"$latitude", userLat,
										},
									}},
									2,
								},
							}},
							bson.D{{
								"$pow", bson.A{
									bson.D{{
										"$subtract", bson.A{
											"$longitude", userLon,
										},
									}},
									2,
								},
							}},
						},
					}},
				}},
			}},
		}},
		{{
			"$match", bson.D{{
				"radius", bson.D{{
					"$lt", radius,
				}},
			}},
		}},
	}
	cursor, err := rr.dbCollection.Aggregate(rr.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(rr.ctx, &rooms); err != nil {
		return nil, err
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}

	return &rooms, nil
}

func (rr *RoomRepository) GetRoomByRoomID(roomID string) (*models.Room, error) {
	return rr.findOne(bson.D{{"roomID", roomID}})
}

func (rr *RoomRepository) GetRoomByID(id primitive.ObjectID) (*models.Room, error) {
	return rr.findOne(bson.D{{"_id", id}})
}

func (rr *RoomRepository) GetAllRooms() ([]models.Room, error) {
	return rr.findMany(bson.D{{}})
}

func (rr *RoomRepository) GetRoomsByCategoryID(id int) ([]models.Room, error) {
	return rr.findMany(bson.D{{"categories.id", id}})
}

func (rr *RoomRepository) findOne(filter interface{}) (*models.Room, error) {
	var room models.Room
	err := rr.dbCollection.FindOne(rr.ctx, filter).Decode(&room)
	return &room, err
}

func (rr *RoomRepository) findMany(filter interface{}) ([]models.Room, error) {
	var rooms []models.Room
	cursor, err := rr.dbCollection.Find(rr.ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(rr.ctx, &rooms); err != nil {
		return nil, err
	}
	if rooms == nil {
		rooms = make([]models.Room, 0)
	}
	return rooms, nil
}

func (rr *RoomRepository) GetRoomsByParentGroupID(parentGroupID string) ([]models.Room, error) {
	return rr.findMany(bson.D{{"parentGroupID", parentGroupID}})
}

func (rr *RoomRepository) UpdateRoom(room *models.Room) (*models.Room, error) {
	err := rr.prepareRoomForInserting(room)
	if err != nil {
		return nil, err
	}
	ttlDuration := time.Duration(room.TTL)
	expiresAtTime := room.CreatedAt.Add(ttlDuration * time.Second)
	room.ExpiresAt = expiresAtTime

	// update all fields except ones which have "readonly" tag
	var setElements bson.D
	v := reflect.ValueOf(*room)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		value := v.FieldByName(name).Interface()
		fieldNameBson := t.Field(i).Tag.Get("bson")
		if fieldNameBson == "" {
			fieldNameBson = name
		}
		if t.Field(i).Tag.Get("readonly") != "true" {
			setElements = append(setElements, bson.E{fieldNameBson, value})
		}
	}

	_, err = rr.dbCollection.UpdateOne(rr.ctx, bson.D{{"_id", room.ID}}, bson.D{{"$set", setElements}})
	return room, err
}

func (rr *RoomRepository) DeleteRoom(id primitive.ObjectID) error {
	_, err := rr.dbCollection.DeleteOne(rr.ctx, bson.D{{"_id", id}})
	return err
}
