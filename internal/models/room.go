package models

type Room struct {
	ID        int64   `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longtitude"`
	TTL       int     `json:"ttl"`
	RoomID    string  `json:"roomId"`
}
