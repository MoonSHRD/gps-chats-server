package repositories

import (
	"context"

	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TicketTypeNamesCollectionName = "ticket_type_names"
)

type TicketTypeNameRepository struct {
	app          *app.App
	dbCollection *mongo.Collection
	ctx          context.Context
}

func NewTicketTypeNameRepository(a *app.App) (*TicketTypeNameRepository, error) {
	ctx := context.Background()
	col := a.MainDatabase.Collection(TicketTypeNamesCollectionName)
	ttnr := &TicketTypeNameRepository{
		app:          a,
		dbCollection: col,
		ctx:          ctx,
	}

	err := ttnr.initMongo()
	if err != nil {
		return nil, err
	}

	return ttnr, nil
}

func (ttnr *TicketTypeNameRepository) initMongo() error {
	// add indexes

	//_, err := ttnr.dbCollection.Indexes().CreateMany(ttnr.ctx, []mongo.IndexModel{})
	return nil
}

func (ttnr *TicketTypeNameRepository) findOne(filter interface{}) (*models.TicketTypeName, error) {
	var ttn models.TicketTypeName
	err := ttnr.dbCollection.FindOne(ttnr.ctx, filter).Decode(&ttn)
	return &ttn, err
}

func (ttnr *TicketTypeNameRepository) findMany(filter interface{}) ([]models.TicketTypeName, error) {
	var ttns []models.TicketTypeName
	cursor, err := ttnr.dbCollection.Find(ttnr.ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ttnr.ctx, &ttns); err != nil {
		return nil, err
	}
	if ttns == nil {
		ttns = make([]models.TicketTypeName, 0)
	}
	return ttns, nil
}

func (ttnr *TicketTypeNameRepository) Put(ttn *models.TicketTypeName) error {
	ttn.ID = primitive.NewObjectID()
	_, err := ttnr.dbCollection.InsertOne(ttnr.ctx, ttn)
	return err
}

func (ttnr *TicketTypeNameRepository) Get(eventID string, typeID int) (*models.TicketTypeName, error) {
	return ttnr.findOne(bson.D{{"eventID", eventID}, {"typeID", typeID}})
}

func (ttnr *TicketTypeNameRepository) Update(ticketTypeName *models.TicketTypeName) error {
	_, err := ttnr.dbCollection.UpdateOne(ttnr.ctx, bson.D{{"_id", ticketTypeName.ID}}, bson.D{{"$set", ticketTypeName}})
	return err
}

func (ttnr *TicketTypeNameRepository) Delete(ticketTypeName *models.TicketTypeName) error {
	_, err := ttnr.dbCollection.DeleteOne(ttnr.ctx, bson.D{{"_id", ticketTypeName.ID}})
	return err
}

