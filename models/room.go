package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Latitude       float32            `json:"latitude,omitempty" bson:"latitude"`
	Longitude      float32            `json:"longitude,omitempty" bson:"longitude"`
	CreatedAt      time.Time          `json:"createdAt,omitempty" bson:"createdAt" readonly:"true"`
	TTL            int                `json:"ttl,omitempty" bson:"ttl"`
	Categories     []ChatCategory     `json:"categories,omitempty" bson:"categories"`
	RoomID         string             `json:"roomID,omitempty" bson:"roomID"`
	ParentGroupID  string             `json:"parentGroupId,omitempty" bson:"parentGroupID"`
	EventStartDate time.Time          `json:"eventStartDate,omitempty" bson:"eventStartDate"`
	Name           string             `json:"name" bson:"name"`
	Address        string             `json:"address" bson:"address"`
	ExpiresAt      time.Time          `json:"expiresAt" bson:"expiresAt"`
}
