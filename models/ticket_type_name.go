package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TicketTypeName struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	EventID  string             `json:"eventID" bson:"eventID"`
	TypeName string             `json:"typeName" bson:"typeName"`
	TypeID   int                `json:"typeID" bson:"typeID"`
}
