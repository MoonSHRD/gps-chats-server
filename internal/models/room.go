package models

type Room struct {
	ID        int64
	Latitude  float32
	Longitude float32
	TTL       int
	RoomID    string
}
