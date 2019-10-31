package models

type Room struct {
	ID        int64   `json:"id,omitempty" db:"id"`
	Latitude  float32 `json:"latitude,omitempty" db:"latitude"`
	Longitude float32 `json:"longtitude,omitempty" db:"longitude"`
	CreatedAt float64 `json:"created_at,omitempty" db:"created_at"`
	TTL       int     `json:"ttl,omitempty" db:"ttl"`
	Category  string  `json:"category,omitempty" db:"category"`
	RoomID    string  `json:"roomId,omitempty" db:"room_id"`
}
