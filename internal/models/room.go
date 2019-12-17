package models

import (
	"time"
)

type Room struct {
	ID         int            `json:"id,omitempty" db:"id"`
	Latitude   float32        `json:"latitude,omitempty" db:"latitude"`
	Longitude  float32        `json:"longitude,omitempty" db:"longitude"`
	CreatedAt  time.Time      `json:"created_at,omitempty" db:"created_at"`
	TTL        int            `json:"ttl,omitempty" db:"ttl"`
	Categories []ChatCategory `json:"categories,omitempty" db:"-"`
	RoomID     string         `json:"roomId,omitempty" db:"room_id"`
	EventID    string         `json:"eventId,omitempty" db:"event_id"`
}
